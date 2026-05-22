package handler

import (
	"strconv"
	"strings"

	"gohotel/internal/service"
	"gohotel/pkg/utils"

	"github.com/gin-gonic/gin"
)

// RecommendationHandler 推荐控制器
type RecommendationHandler struct {
	recommendationService *service.RecommendationService
}

// NewRecommendationHandler 创建推荐控制器实例
func NewRecommendationHandler(recommendationService *service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{recommendationService: recommendationService}
}

// GetRoomRecommendations 获取房型推荐
func (h *RecommendationHandler) GetRoomRecommendations(c *gin.Context) {
	limit := getQueryInt(c, "limit", 6)
	excludeRoomID, _ := strconv.ParseInt(c.DefaultQuery("exclude_room_id", "0"), 10, 64)
	userID := extractOptionalUserID(c)

	recommendations, err := h.recommendationService.GetRoomRecommendations(userID, limit, excludeRoomID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, recommendations)
}

// TrackBehavior 记录推荐行为
func (h *RecommendationHandler) TrackBehavior(c *gin.Context) {
	userIDAny, _ := c.Get("user_id")
	userID, _ := userIDAny.(int64)

	var req service.TrackBehaviorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	if err := h.recommendationService.TrackBehavior(userID, &req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "行为记录成功", nil)
}

func extractOptionalUserID(c *gin.Context) *int64 {
	authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
	if authHeader == "" {
		return nil
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil
	}

	claims, err := utils.ParseToken(parts[1])
	if err != nil {
		return nil
	}

	userID := claims.UserID
	return &userID
}
