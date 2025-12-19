package service

import (
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/pkg/utils"
	"sync"
	"time"
)

// NoticeService 公告业务逻辑层
type NoticeService struct {
	noticeRepo *repository.NoticeRepository
	cosService *CosService
	timeWheel  *utils.MultiTimeWheel       // 时间轮实例（重构后使用多层时间轮）
	taskMap    map[int64]map[string]string // 任务ID映射，key: noticeID, value: map[taskType]taskID
	taskMutex  sync.RWMutex                // 保护taskMap的互斥锁
}

// NoticeTaskExecutor 公告任务执行器，用于处理公告相关的定时任务
type NoticeTaskExecutor struct {
	noticeService *NoticeService
}

// NewNoticeTaskExecutor 创建公告任务执行器实例
func NewNoticeTaskExecutor(noticeService *NoticeService) *NoticeTaskExecutor {
	return &NoticeTaskExecutor{
		noticeService: noticeService,
	}
}

// Execute 执行公告任务
func (e *NoticeTaskExecutor) Execute(task *utils.Task) {
	// 从任务元数据中获取noticeID和任务类型
	noticeIDFloat, ok := task.Meta["notice_id"].(float64)
	if !ok {
		return
	}
	noticeID := int64(noticeIDFloat)

	taskType, ok := task.Meta["task_type"].(string)
	if !ok {
		return
	}

	// 根据任务类型执行相应操作
	switch taskType {
	case "start":
		// 开始时间到，将公告状态改为active
		e.noticeService.UpdateNoticeStatus(noticeID, "active")
	case "end":
		// 结束时间到，将公告状态改为inactive
		e.noticeService.UpdateNoticeStatus(noticeID, "inactive")
	}

	// 从任务映射中移除已执行的任务
	e.noticeService.taskMutex.Lock()
	defer e.noticeService.taskMutex.Unlock()

	if taskMap, exists := e.noticeService.taskMap[noticeID]; exists {
		delete(taskMap, taskType)
		// 如果该公告没有任务了，删除映射
		if len(taskMap) == 0 {
			delete(e.noticeService.taskMap, noticeID)
		}
	}
}

// GetTaskType 获取任务类型
func (e *NoticeTaskExecutor) GetTaskType() string {
	return "notice"
}

// NewNoticeService 创建公告服务实例
func NewNoticeService(noticeRepo *repository.NoticeRepository, cosService *CosService, timeWheel *utils.MultiTimeWheel) *NoticeService {
	service := &NoticeService{
		noticeRepo: noticeRepo,
		cosService: cosService,
		timeWheel:  timeWheel,
		taskMap:    make(map[int64]map[string]string),
	}

	// 创建并注册公告任务执行器
	executor := NewNoticeTaskExecutor(service)
	timeWheel.RegisterExecutor(executor)

	return service
}

// CreateNoticeRequest 创建公告请求结构
type CreateNoticeRequest struct {
	Title     string  `json:"title" binding:"required,min=1,max=100"`
	LinkURL   *string `json:"link_url" binding:"omitempty,max=500"`
	Sort      *int    `json:"sort" binding:"omitempty,min=0"`
	StartTime *string `json:"start_time" binding:"omitempty"`
	EndTime   *string `json:"end_time" binding:"omitempty"`
}

// UpdateNoticeRequest 更新公告请求结构
type UpdateNoticeRequest struct {
	Title     string  `json:"title" binding:"required,min=1,max=100"`
	LinkURL   *string `json:"link_url" binding:"omitempty,max=500"`
	Sort      *int    `json:"sort" binding:"omitempty,min=0"`
	StartTime *string `json:"start_time" binding:"omitempty"`
	EndTime   *string `json:"end_time" binding:"omitempty"`
}

// parseNoticeTimeString 将字符串转换为*time.Time
// 支持的格式：2006-01-02 15:04:05
// 如果字符串为空或解析失败，返回nil
func parseNoticeTimeString(timeStr *string) *time.Time {
	if timeStr == nil || *timeStr == "" {
		return nil
	}

	// 解析时间字符串为本地时间
	t, err := time.ParseInLocation("2006-01-02 15:04:05", *timeStr, time.Local)
	if err != nil {
		return nil
	}

	return &t
}

// CreateNotice 创建公告
func (s *NoticeService) CreateNotice(req *CreateNoticeRequest) (*models.Notice, error) {
	// 解析时间
	startTime := parseNoticeTimeString(req.StartTime)
	endTime := parseNoticeTimeString(req.EndTime)
	now := time.Now()

	sort := 0
	if req.Sort != nil {
		sort = *req.Sort
	}

	// 根据时间自动计算状态
	status := "inactive"
	// 如果开始时间已过且结束时间未过，设为active
	if (startTime == nil || now.After(*startTime)) &&
		(endTime == nil || now.Before(*endTime)) {
		status = "active"
	}

	// 创建公告对象
	notice := &models.Notice{
		ID:        utils.JSONInt64(utils.GenID()),
		Title:     req.Title,
		LinkURL:   req.LinkURL,
		Sort:      sort,
		Status:    status,
		StartTime: startTime,
		EndTime:   endTime,
	}

	// 保存到数据库
	if err := s.noticeRepo.Create(notice); err != nil {
		return nil, err
	}

	// 添加时间轮任务
	s.addNoticeTasks(notice)

	return notice, nil
}

// GetNoticeByID 根据ID获取公告
func (s *NoticeService) GetNoticeByID(id int64) (*models.Notice, error) {
	return s.noticeRepo.FindByID(id)
}

// GetAllNotices 获取所有公告（带分页）
func (s *NoticeService) GetAllNotices(page, pageSize int) ([]models.Notice, int64, error) {
	return s.noticeRepo.FindAll(page, pageSize)
}

// GetActiveNotices 获取激活的公告（前端展示用）
func (s *NoticeService) GetActiveNotices() ([]models.Notice, error) {
	return s.noticeRepo.FindActive()
}

// UpdateNotice 更新公告信息
func (s *NoticeService) UpdateNotice(id int64, req *UpdateNoticeRequest) (*models.Notice, error) {
	// 获取现有公告
	notice, err := s.noticeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Title != "" {
		notice.Title = req.Title
	}

	if req.LinkURL != nil {
		notice.LinkURL = req.LinkURL
	}
	if req.Sort != nil {
		notice.Sort = *req.Sort
	}

	// 更新时间字段
	if req.StartTime != nil {
		notice.StartTime = parseNoticeTimeString(req.StartTime)
	}
	if req.EndTime != nil {
		notice.EndTime = parseNoticeTimeString(req.EndTime)
	}

	// 根据时间自动计算状态
	now := time.Now()
	status := "inactive"
	// 如果开始时间已过且结束时间未过，设为active
	if (notice.StartTime == nil || now.After(*notice.StartTime)) &&
		(notice.EndTime == nil || now.Before(*notice.EndTime)) {
		status = "active"
	}
	notice.Status = status

	// 保存更新
	if err := s.noticeRepo.Update(notice); err != nil {
		return nil, err
	}

	// 更新时间轮任务：先删除旧任务，再添加新任务
	s.removeNoticeTasks(id)
	s.addNoticeTasks(notice)

	return notice, nil
}

// UpdateNoticeStatus 更新公告状态
// 内部使用，由时间轮任务调用
func (s *NoticeService) UpdateNoticeStatus(id int64, status string) error {
	return s.noticeRepo.UpdateStatus(id, status)
}

// removeNoticeTasks 删除公告的时间轮任务
func (s *NoticeService) removeNoticeTasks(noticeID int64) {
	s.taskMutex.Lock()
	defer s.taskMutex.Unlock()

	// 获取该公告的任务映射
	if taskMap, exists := s.taskMap[noticeID]; exists {
		// 删除开始时间任务
		if startTaskID, ok := taskMap["start"]; ok {
			s.timeWheel.RemoveTask(startTaskID)
			delete(taskMap, "start")
		}

		// 删除结束时间任务
		if endTaskID, ok := taskMap["end"]; ok {
			s.timeWheel.RemoveTask(endTaskID)
			delete(taskMap, "end")
		}

		// 如果该公告没有任务了，删除映射
		if len(taskMap) == 0 {
			delete(s.taskMap, noticeID)
		}
	}
}

// addNoticeTasks 为公告添加时间轮任务
func (s *NoticeService) addNoticeTasks(notice *models.Notice) {
	noticeID := int64(notice.ID)

	// 初始化该公告的任务映射
	s.taskMutex.Lock()
	if _, exists := s.taskMap[noticeID]; !exists {
		s.taskMap[noticeID] = make(map[string]string)
	}
	taskMap := s.taskMap[noticeID]
	s.taskMutex.Unlock()

	// 添加开始时间任务
	if notice.StartTime != nil {
		startTime := notice.StartTime
		// 准备任务元数据
		meta := map[string]interface{}{
			"executor_type": "notice", // 执行器类型
			"notice_id":     noticeID,
			"task_type":     "start",
		}
		// 到时间将状态改为active
		taskID := s.timeWheel.AddTask(*startTime, func() {
			s.UpdateNoticeStatus(noticeID, "active")
			// 任务执行后从映射中移除
			s.taskMutex.Lock()
			delete(taskMap, "start")
			if len(taskMap) == 0 {
				delete(s.taskMap, noticeID)
			}
			s.taskMutex.Unlock()
		}, meta)

		// 保存任务ID
		s.taskMutex.Lock()
		taskMap["start"] = taskID
		s.taskMutex.Unlock()
	}

	// 添加结束时间任务
	if notice.EndTime != nil {
		endTime := notice.EndTime
		// 准备任务元数据
		meta := map[string]interface{}{
			"executor_type": "notice", // 执行器类型
			"notice_id":     noticeID,
			"task_type":     "end",
		}
		// 到时间将状态改为inactive
		taskID := s.timeWheel.AddTask(*endTime, func() {
			s.UpdateNoticeStatus(noticeID, "inactive")
			// 任务执行后从映射中移除
			s.taskMutex.Lock()
			delete(taskMap, "end")
			if len(taskMap) == 0 {
				delete(s.taskMap, noticeID)
			}
			s.taskMutex.Unlock()
		}, meta)

		// 保存任务ID
		s.taskMutex.Lock()
		taskMap["end"] = taskID
		s.taskMutex.Unlock()
	}
}

// DeleteNotice 删除公告
func (s *NoticeService) DeleteNotice(id int64) error {
	notice, err := s.noticeRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 删除时间轮任务
	s.removeNoticeTasks(int64(notice.ID))

	//删除数据库中的公告记录
	return s.noticeRepo.Delete(id)
}
