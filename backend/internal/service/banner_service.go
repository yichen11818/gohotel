package service

import (
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/pkg/utils"
	"sync"
	"time"
)

// BannerService 活动横幅业务逻辑层
type BannerService struct {
	bannerRepo *repository.BannerRepository
	cosService *CosService
	timeWheel  *utils.MultiTimeWheel       // 时间轮实例（重构后使用多层时间轮）
	taskMap    map[int64]map[string]string // 任务ID映射，key: bannerID, value: map[taskType]taskID
	taskMutex  sync.RWMutex                // 保护taskMap的互斥锁
}

// BannerTaskExecutor banner任务执行器，用于处理banner相关的定时任务
type BannerTaskExecutor struct {
	bannerService *BannerService
}

// NewBannerTaskExecutor 创建banner任务执行器实例
func NewBannerTaskExecutor(bannerService *BannerService) *BannerTaskExecutor {
	return &BannerTaskExecutor{
		bannerService: bannerService,
	}
}

// Execute 执行banner任务
func (e *BannerTaskExecutor) Execute(task *utils.Task) {
	// 从任务元数据中获取bannerID和任务类型
	bannerID, ok := task.Meta["banner_id"].(int64)
	if !ok {
		return
	}

	taskType, ok := task.Meta["task_type"].(string)
	if !ok {
		return
	}

	// 根据任务类型执行相应操作
	switch taskType {
	case "start":
		// 开始时间到，将banner状态改为active
		e.bannerService.UpdateBannerStatus(bannerID, "active")
	case "end":
		// 结束时间到，将banner状态改为inactive
		e.bannerService.UpdateBannerStatus(bannerID, "inactive")
	}

	// 从任务映射中移除已执行的任务
	e.bannerService.taskMutex.Lock()
	defer e.bannerService.taskMutex.Unlock()

	if taskMap, exists := e.bannerService.taskMap[bannerID]; exists {
		delete(taskMap, taskType)
		// 如果该banner没有任务了，删除映射
		if len(taskMap) == 0 {
			delete(e.bannerService.taskMap, bannerID)
		}
	}
}

// GetTaskType 获取任务类型
func (e *BannerTaskExecutor) GetTaskType() string {
	return "banner"
}

// NewBannerService 创建活动横幅服务实例
func NewBannerService(bannerRepo *repository.BannerRepository, cosService *CosService, timeWheel *utils.MultiTimeWheel) *BannerService {
	service := &BannerService{
		bannerRepo: bannerRepo,
		cosService: cosService,
		timeWheel:  timeWheel,
		taskMap:    make(map[int64]map[string]string),
	}

	// 创建并注册banner任务执行器
	executor := NewBannerTaskExecutor(service)
	timeWheel.RegisterExecutor(executor)

	return service
}

// CreateBannerRequest 创建活动横幅请求结构
type CreateBannerRequest struct {
	Title     string  `json:"title" binding:"required,min=1,max=100"`
	Subtitle  *string `json:"subtitle" binding:"omitempty,max=255"`
	ImageURL  string  `json:"image_url" binding:"required,max=500"`
	LinkURL   *string `json:"link_url" binding:"omitempty,max=500"`
	Sort      int     `json:"sort" binding:"omitempty,min=0"`
	StartTime *string `json:"start_time" binding:"omitempty"`
	EndTime   *string `json:"end_time" binding:"omitempty"`
}

// UpdateBannerRequest 更新活动横幅请求结构
type UpdateBannerRequest struct {
	Title     string  `json:"title" binding:"omitempty,min=1,max=100"`
	Subtitle  *string `json:"subtitle" binding:"omitempty,max=255"`
	ImageURL  string  `json:"image_url" binding:"omitempty,max=500"`
	LinkURL   *string `json:"link_url" binding:"omitempty,max=500"`
	Sort      *int    `json:"sort" binding:"omitempty,min=0"`
	StartTime *string `json:"start_time" binding:"omitempty"`
	EndTime   *string `json:"end_time" binding:"omitempty"`
}

// parseTimeString 将字符串转换为*time.Time
// 支持的格式：2006-01-02 15:04:05
// 如果字符串为空或解析失败，返回nil
func parseTimeString(timeStr *string) *time.Time {
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

// CreateBanner 创建活动横幅
func (s *BannerService) CreateBanner(req *CreateBannerRequest) (*models.Banner, error) {
	// 解析时间
	startTime := parseTimeString(req.StartTime)
	endTime := parseTimeString(req.EndTime)
	now := time.Now()

	// 根据时间自动计算状态
	status := "inactive"
	// 如果开始时间已过且结束时间未过，设为active
	if (startTime == nil || now.After(*startTime)) &&
		(endTime == nil || now.Before(*endTime)) {
		status = "active"
	}

	// 创建Banner对象
	banner := &models.Banner{
		Title:     req.Title,
		Subtitle:  req.Subtitle,
		ImageURL:  req.ImageURL,
		LinkURL:   req.LinkURL,
		Status:    status,
		Sort:      req.Sort,
		StartTime: startTime,
		EndTime:   endTime,
	}

	// 保存到数据库
	if err := s.bannerRepo.Create(banner); err != nil {
		return nil, err
	}

	// 添加时间轮任务
	s.addBannerTasks(banner)

	return banner, nil
}

// GetBannerByID 根据ID获取活动横幅
func (s *BannerService) GetBannerByID(id int64) (*models.Banner, error) {
	return s.bannerRepo.FindByID(id)
}

// GetAllBanners 获取所有活动横幅（带分页）
func (s *BannerService) GetAllBanners(page, pageSize int) ([]models.Banner, int64, error) {
	return s.bannerRepo.FindAll(page, pageSize)
}

// GetActiveBanners 获取激活的活动横幅（前端展示用）
func (s *BannerService) GetActiveBanners() ([]models.Banner, error) {
	return s.bannerRepo.FindActive()
}

// UpdateBanner 更新活动横幅信息
func (s *BannerService) UpdateBanner(id int64, req *UpdateBannerRequest) (*models.Banner, error) {
	// 获取现有Banner
	banner, err := s.bannerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 保存旧图片URL
	oldImageURL := banner.ImageURL

	// 更新字段
	if req.Title != "" {
		banner.Title = req.Title
	}
	if req.Subtitle != nil {
		banner.Subtitle = req.Subtitle
	}
	if req.ImageURL != "" {
		// 更新图片URL
		banner.ImageURL = req.ImageURL
		// 删除旧图片（如果存在且与新图片不同）
		if oldImageURL != "" && oldImageURL != req.ImageURL {
			if err := s.cosService.DeleteFile(oldImageURL); err != nil {
				return nil, err
			}
		}
	}
	if req.LinkURL != nil {
		banner.LinkURL = req.LinkURL
	}
	if req.Sort != nil {
		banner.Sort = *req.Sort
	}

	// 更新时间字段
	if req.StartTime != nil {
		banner.StartTime = parseTimeString(req.StartTime)
	}
	if req.EndTime != nil {
		banner.EndTime = parseTimeString(req.EndTime)
	}

	// 根据时间自动计算状态
	now := time.Now()
	status := "inactive"
	// 如果开始时间已过且结束时间未过，设为active
	if (banner.StartTime == nil || now.After(*banner.StartTime)) &&
		(banner.EndTime == nil || now.Before(*banner.EndTime)) {
		status = "active"
	}
	banner.Status = status

	// 保存更新
	if err := s.bannerRepo.Update(banner); err != nil {
		return nil, err
	}

	// 更新时间轮任务：先删除旧任务，再添加新任务
	s.removeBannerTasks(id)
	s.addBannerTasks(banner)

	return banner, nil
}

// UpdateBannerStatus 更新活动横幅状态
// 内部使用，由时间轮任务调用
func (s *BannerService) UpdateBannerStatus(id int64, status string) error {
	return s.bannerRepo.UpdateStatus(id, status)
}

// removeBannerTasks 删除banner的时间轮任务
func (s *BannerService) removeBannerTasks(bannerID int64) {
	s.taskMutex.Lock()
	defer s.taskMutex.Unlock()

	// 获取该banner的任务映射
	if taskMap, exists := s.taskMap[bannerID]; exists {
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

		// 如果该banner没有任务了，删除映射
		if len(taskMap) == 0 {
			delete(s.taskMap, bannerID)
		}
	}
}

// addBannerTasks 为banner添加时间轮任务
func (s *BannerService) addBannerTasks(banner *models.Banner) {
	bannerID := banner.ID

	// 初始化该banner的任务映射
	s.taskMutex.Lock()
	if _, exists := s.taskMap[bannerID]; !exists {
		s.taskMap[bannerID] = make(map[string]string)
	}
	taskMap := s.taskMap[bannerID]
	s.taskMutex.Unlock()

	// 添加开始时间任务
	if banner.StartTime != nil {
		startTime := banner.StartTime
		// 准备任务元数据
		meta := map[string]interface{}{
			"executor_type": "banner", // 执行器类型
			"banner_id":     bannerID,
			"task_type":     "start",
		}
		// 到时间将状态改为active
		taskID := s.timeWheel.AddTask(*startTime, func() {
			s.UpdateBannerStatus(bannerID, "active")
			// 任务执行后从映射中移除
			s.taskMutex.Lock()
			delete(taskMap, "start")
			if len(taskMap) == 0 {
				delete(s.taskMap, bannerID)
			}
			s.taskMutex.Unlock()
		}, meta)

		// 保存任务ID
		s.taskMutex.Lock()
		taskMap["start"] = taskID
		s.taskMutex.Unlock()
	}

	// 添加结束时间任务
	if banner.EndTime != nil {
		endTime := banner.EndTime
		// 准备任务元数据
		meta := map[string]interface{}{
			"executor_type": "banner", // 执行器类型
			"banner_id":     bannerID,
			"task_type":     "end",
		}
		// 到时间将状态改为inactive
		taskID := s.timeWheel.AddTask(*endTime, func() {
			s.UpdateBannerStatus(bannerID, "inactive")
			// 任务执行后从映射中移除
			s.taskMutex.Lock()
			delete(taskMap, "end")
			if len(taskMap) == 0 {
				delete(s.taskMap, bannerID)
			}
			s.taskMutex.Unlock()
		}, meta)

		// 保存任务ID
		s.taskMutex.Lock()
		taskMap["end"] = taskID
		s.taskMutex.Unlock()
	}
}

// DeleteBanner 删除活动横幅
func (s *BannerService) DeleteBanner(id int64) error {
	banner, err := s.bannerRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := s.cosService.DeleteFile(banner.ImageURL); err != nil {
		return err
	}

	// 删除时间轮任务
	s.removeBannerTasks(banner.ID)

	//删除数据库中的banner记录
	return s.bannerRepo.Delete(id)
}
