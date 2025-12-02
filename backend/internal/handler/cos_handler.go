package handler

import (
	"gohotel/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CosHandler 腾讯云对象存储API处理层
type CosHandler struct {
	cosService *service.CosService
}

// NewCosHandler 创建对象存储处理器实例
func NewCosHandler(cosService *service.CosService) *CosHandler {
	return &CosHandler{cosService: cosService}
}

// UploadImage 通用图片上传接口
// @Summary 通用图片上传
// @Description 上传图片到腾讯云对象存储，返回临时URL
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param image formData file true "图片文件"
// @Param type formData string true "资源类型：banner/room/user等"
// @Success 200 {object} map[string]string "{\"temp_url\": \"临时图片URL\"}"
// @Failure 400 {object} map[string]string "{\"error\": \"错误信息\"}"
// @Failure 500 {object} map[string]string "{\"error\": \"错误信息\"}"
// @Router /api/upload/image [post]
func (h *CosHandler) UploadImage(c *gin.Context) {
	// 获取上传的文件和资源类型
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取上传文件失败: " + err.Error()})
		return
	}
	resourceType := c.PostForm("type")
	if resourceType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "资源类型不能为空"})
		return
	}

	// 上传到临时目录
	tempURL, err := h.cosService.UploadTempFile(file, resourceType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "图片上传失败: " + err.Error()})
		return
	}

	// 返回临时URL
	c.JSON(http.StatusOK, gin.H{"temp_url": tempURL})
}
