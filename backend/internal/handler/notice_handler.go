package handler

import (
	"gohotel/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// NoticeHandler 公告API处理层
type NoticeHandler struct {
	noticeService *service.NoticeService
}

// NewNoticeHandler 创建公告处理器实例
func NewNoticeHandler(noticeService *service.NoticeService) *NoticeHandler {
	return &NoticeHandler{
		noticeService: noticeService,
	}
}

// CreateNotice 创建公告
// @Summary 创建公告
// @Description 创建新的公告
// @Tags 公告管理
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param title formData string true "公告标题"
// @Param link_url formData string false "点击跳转链接"
// @Param sort formData int false "展示顺序"
// @Param start_time formData string false "公告开始时间"
// @Param end_time formData string false "公告结束时间"
// @Success 200 {object} models.Notice
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/notices [post]
func (h *NoticeHandler) CreateNotice(c *gin.Context) {
	// 从表单中获取公告信息
	title := c.PostForm("title")
	linkURL, hasLinkURL := c.GetPostForm("link_url")
	sortStr, hasSort := c.GetPostForm("sort")
	startTime, hasStartTime := c.GetPostForm("start_time")
	endTime, hasEndTime := c.GetPostForm("end_time")

	// 验证必要字段
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题不能为空"})
		return
	}

	// 解析sort字段
	var sort *int
	if hasSort && sortStr != "" {
		if s, parseErr := strconv.Atoi(sortStr); parseErr == nil {
			sort = &s
		}
	}

	var linkURLPtr *string
	if hasLinkURL {
		linkURLPtr = &linkURL
	}
	var startTimePtr *string
	if hasStartTime {
		startTimePtr = &startTime
	}
	var endTimePtr *string
	if hasEndTime {
		endTimePtr = &endTime
	}

	// 构建请求对象
	req := &service.CreateNoticeRequest{
		Title:     title,
		LinkURL:   linkURLPtr,
		Sort:      sort,
		StartTime: startTimePtr,
		EndTime:   endTimePtr,
	}

	// 调用服务层创建公告
	notice, err := h.noticeService.CreateNotice(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, notice)
}

// GetNoticeByID 根据ID获取公告
// @Summary 根据ID获取公告
// @Description 根据ID获取公告详情
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "公告ID"
// @Success 200 {object} models.Notice
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/notices/{id} [get]
func (h *NoticeHandler) GetNoticeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64) // ParseInt returns int64
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	notice, err := h.noticeService.GetNoticeByID(int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "公告不存在"})
		return
	}

	c.JSON(http.StatusOK, notice)
}

// GetAllNotices 获取所有公告（带分页）
// @Summary 获取所有公告
// @Description 获取所有公告，支持分页
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码，默认1"
// @Param pageSize query int false "每页条数，默认10"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/admin/notices [get]
func (h *NoticeHandler) GetAllNotices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	notices, total, err := h.noticeService.GetAllNotices(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notices":  notices,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetActiveNotices 获取激活的公告（前端展示用）
// @Summary 获取激活的公告
// @Description 获取激活状态的公告，用于前端展示
// @Tags 公告管理
// @Accept json
// @Produce json
// @Success 200 {array} models.Notice
// @Failure 500 {object} map[string]string
// @Router /api/notices/active [get]
func (h *NoticeHandler) GetActiveNotices(c *gin.Context) {
	notices, err := h.noticeService.GetActiveNotices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, notices)
}

// UpdateNotice 更新公告
// @Summary 更新公告
// @Description 更新公告信息
// @Tags 公告管理
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param id path string true "公告ID"
// @Param title formData string false "公告标题"
// @Param link_url formData string false "点击跳转链接"
// @Param sort formData int false "展示顺序"
// @Param start_time formData string false "公告开始时间"
// @Param end_time formData string false "公告结束时间"
// @Success 200 {object} models.Notice
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/notices/{id} [post]
func (h *NoticeHandler) UpdateNotice(c *gin.Context) {
	// 获取Banner ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	// 检查公告是否存在
	_, err = h.noticeService.GetNoticeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "公告不存在"})
		return
	}

	// 从表单中获取公告信息
	title, hasTitle := c.GetPostForm("title")
	linkURL, hasLinkURL := c.GetPostForm("link_url")
	sortStr, hasSort := c.GetPostForm("sort")
	startTime, hasStartTime := c.GetPostForm("start_time")
	endTime, hasEndTime := c.GetPostForm("end_time")

	// 解析sort字段
	var sort *int
	if hasSort && sortStr != "" {
		if s, parseErr := strconv.Atoi(sortStr); parseErr == nil {
			sort = &s
		}
	}

	var linkURLPtr *string
	if hasLinkURL {
		linkURLPtr = &linkURL
	}
	var startTimePtr *string
	if hasStartTime {
		startTimePtr = &startTime
	}
	var endTimePtr *string
	if hasEndTime {
		endTimePtr = &endTime
	}

	// 构建请求对象
	req := &service.UpdateNoticeRequest{
		LinkURL:   linkURLPtr,
		Sort:      sort,
		StartTime: startTimePtr,
		EndTime:   endTimePtr,
	}
	if hasTitle {
		req.Title = title
	}

	// 调用服务层更新公告
	notice, err := h.noticeService.UpdateNotice(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, notice)
}

// DeleteBanner 删除活动横幅
// @Summary 删除活动横幅
// @Description 删除指定ID的活动横幅
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "公告ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/notices/{id}/delete [post]
func (h *NoticeHandler) DeleteNotice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.noticeService.DeleteNotice(int64(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
