package handler

import (
	"gohotel/internal/service"
	"gohotel/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BannerHandler 活动横幅API处理层
type BannerHandler struct {
	bannerService *service.BannerService
	cosService    *service.CosService
}

// NewBannerHandler 创建活动横幅处理器实例
func NewBannerHandler(bannerService *service.BannerService, cosService *service.CosService) *BannerHandler {
	return &BannerHandler{
		bannerService: bannerService,
		cosService:    cosService,
	}
}

// CreateBanner 创建活动横幅
// @Summary 创建活动横幅
// @Description 创建新的活动横幅，使用临时图片URL
// @Tags 活动管理
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param title formData string true "活动横幅标题"
// @Param subtitle formData string false "活动横幅副标题"
// @Param temp_url formData string true "临时图片URL（通过通用上传接口获取）"
// @Param link_url formData string false "点击跳转链接"
// @Param sort formData int false "展示顺序"
// @Param start_time formData string false "活动开始时间"
// @Param end_time formData string false "活动结束时间"
// @Success 200 {object} models.Banner
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/banners [post]
func (h *BannerHandler) CreateBanner(c *gin.Context) {
	// 从表单中获取Banner信息
	title := c.PostForm("title")
	subtitle := c.PostForm("subtitle")
	tempURL := c.PostForm("temp_url")
	linkURL := c.PostForm("link_url")
	sortStr := c.PostForm("sort")
	startTime := c.PostForm("start_time")
	endTime := c.PostForm("end_time")

	// 验证必要字段
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题不能为空"})
		return
	}
	if tempURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "临时图片URL不能为空"})
		return
	}

	// 确认上传，获取正式URL
	imageURL, err := h.cosService.ConfirmUpload(tempURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "确认图片上传失败: " + err.Error()})
		return
	}

	// 解析sort字段
	sort := 0
	if sortStr != "" {
		if s, parseErr := strconv.Atoi(sortStr); parseErr == nil {
			sort = s
		}
	}

	// 构建请求对象
	req := &service.CreateBannerRequest{
		Title:     title,
		Subtitle:  &subtitle,
		ImageURL:  imageURL,
		LinkURL:   &linkURL,
		Sort:      sort,
		StartTime: &startTime,
		EndTime:   &endTime,
	}

	// 调用服务层创建Banner
	banner, err := h.bannerService.CreateBanner(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, banner)
}

// GetBannerByID 根据ID获取活动横幅
// @Summary 根据ID获取活动横幅
// @Description 根据ID获取活动横幅详情
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "活动横幅ID"
// @Success 200 {object} models.Banner
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/banners/{id} [get]
func (h *BannerHandler) GetBannerByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	banner, err := h.bannerService.GetBannerByID(int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "活动横幅不存在"})
		return
	}

	c.JSON(http.StatusOK, banner)
}

// GetAllBanners 获取所有活动横幅（带分页）
// @Summary 获取所有活动横幅
// @Description 获取所有活动横幅，支持分页
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码，默认1"
// @Param pageSize query int false "每页条数，默认10"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/admin/banners [get]
func (h *BannerHandler) GetAllBanners(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	banners, total, err := h.bannerService.GetAllBanners(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"banners":  banners,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetActiveBanners 获取激活的活动横幅（前端展示用）
// @Summary 获取激活的活动横幅
// @Description 获取激活状态的活动横幅，用于前端展示
// @Tags 活动管理
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Banner}
// @Failure 500 {object} utils.Response
// @Router /api/banners/active [get]
func (h *BannerHandler) GetActiveBanners(c *gin.Context) {
	banners, err := h.bannerService.GetActiveBanners()
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, banners)
}

// UpdateBanner 更新活动横幅
// @Summary 更新活动横幅
// @Description 更新活动横幅信息，使用临时图片URL
// @Tags 活动管理
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param id path string true "活动横幅ID"
// @Param title formData string false "活动横幅标题"
// @Param subtitle formData string false "活动横幅副标题"
// @Param temp_url formData string false "临时图片URL（通过通用上传接口获取）"
// @Param link_url formData string false "点击跳转链接"
// @Param sort formData int false "展示顺序"
// @Param start_time formData string false "活动开始时间"
// @Param end_time formData string false "活动结束时间"
// @Success 200 {object} models.Banner
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/banners/{id} [post]
func (h *BannerHandler) UpdateBanner(c *gin.Context) {
	// 获取Banner ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	// 检查Banner是否存在
	_, err = h.bannerService.GetBannerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Banner不存在"})
		return
	}

	// 从表单中获取Banner信息
	title, hasTitle := c.GetPostForm("title")
	subtitle, hasSubtitle := c.GetPostForm("subtitle")
	tempURL, hasTempURL := c.GetPostForm("temp_url")
	linkURL, hasLinkURL := c.GetPostForm("link_url")
	sortStr, hasSort := c.GetPostForm("sort")
	startTime, hasStartTime := c.GetPostForm("start_time")
	endTime, hasEndTime := c.GetPostForm("end_time")

	imageURL := ""
	if hasTempURL {
		if tempURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "临时图片URL不能为空"})
			return
		}
		// 确认上传，获取正式URL
		imageURL, err = h.cosService.ConfirmUpload(tempURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "确认图片上传失败: " + err.Error()})
			return
		}
	}

	// 解析sort字段
	var sort *int
	if hasSort && sortStr != "" {
		if s, parseErr := strconv.Atoi(sortStr); parseErr == nil {
			sort = &s
		}
	}

	var subtitlePtr *string
	if hasSubtitle {
		subtitlePtr = &subtitle
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
	req := &service.UpdateBannerRequest{
		Subtitle:  subtitlePtr,
		ImageURL:  imageURL,
		LinkURL:   linkURLPtr,
		Sort:      sort,
		StartTime: startTimePtr,
		EndTime:   endTimePtr,
	}
	if hasTitle {
		req.Title = title
	}

	// 调用服务层更新Banner
	banner, err := h.bannerService.UpdateBanner(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, banner)
}

// DeleteBanner 删除活动横幅
// @Summary 删除活动横幅
// @Description 删除指定ID的活动横幅
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "活动横幅ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/banners/{id}/delete [post]
func (h *BannerHandler) DeleteBanner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.bannerService.DeleteBanner(int64(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
