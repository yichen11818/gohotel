package service

import (
	"context"
	"encoding/json"
	"math"
	"sort"
	"strings"

	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"

	"gorm.io/gorm"
)

// RecommendationService 混合推荐服务
// 通过协同过滤分数 + 轻量神经网络打分，生成房型推荐结果。
type RecommendationService struct {
	db       *gorm.DB
	roomRepo repository.RoomRepository
	userRepo *repository.UserRepository
}

// TrackBehaviorRequest 行为埋点请求
type TrackBehaviorRequest struct {
	RoomID       int64                  `json:"room_id" binding:"required"`
	BehaviorType string                 `json:"behavior_type" binding:"required"`
	Source       string                 `json:"source"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// RecommendedRoom 推荐结果
type RecommendedRoom struct {
	Room               models.Room `json:"room"`
	Score              float64     `json:"score"`
	CollaborativeScore float64     `json:"collaborative_score"`
	NeuralScore        float64     `json:"neural_score"`
	PopularityScore    float64     `json:"popularity_score"`
	Reason             string      `json:"reason"`
	Tags               []string    `json:"tags"`
}

type userRoomSignal struct {
	UserID int64   `gorm:"column:user_id"`
	RoomID int64   `gorm:"column:room_id"`
	Weight float64 `gorm:"column:weight"`
}

type roomFeatureStats struct {
	MinPrice      float64
	MaxPrice      float64
	MinArea       float64
	MaxArea       float64
	MinCapacity   float64
	MaxCapacity   float64
	MaxPopularity float64
}

type userPreferenceSummary struct {
	HasHistory     bool
	PreferredTypes map[string]float64
	AvgPrice       float64
	AvgArea        float64
	AvgCapacity    float64
	SpendLevel     float64
	LevelBoost     float64
}

// NewRecommendationService 创建推荐服务实例
func NewRecommendationService(db *gorm.DB, roomRepo repository.RoomRepository, userRepo *repository.UserRepository) *RecommendationService {
	return &RecommendationService{
		db:       db,
		roomRepo: roomRepo,
		userRepo: userRepo,
	}
}

// TrackBehavior 记录用户行为
func (s *RecommendationService) TrackBehavior(userID int64, req *TrackBehaviorRequest) error {
	if userID <= 0 {
		return errors.NewUnauthorizedError("未登录")
	}

	room, err := s.roomRepo.FindByID(uint(req.RoomID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("房间不存在")
		}
		return errors.NewDatabaseError("find room", err)
	}

	weight, ok := recommendationBehaviorWeight(req.BehaviorType)
	if !ok {
		return errors.NewBadRequestError("不支持的行为类型")
	}

	metadataPayload := ""
	if len(req.Metadata) > 0 {
		payload, err := json.Marshal(req.Metadata)
		if err != nil {
			return errors.NewBadRequestError("行为元数据格式错误")
		}
		metadataPayload = string(payload)
	}

	behavior := &models.UserBehavior{
		ID:           utils.JSONInt64(utils.GenID()),
		UserID:       utils.JSONInt64(userID),
		RoomID:       int64(room.ID),
		BehaviorType: normalizeBehaviorType(req.BehaviorType),
		Source:       strings.TrimSpace(req.Source),
		Weight:       weight,
		MetadataJSON: metadataPayload,
	}

	if err := s.db.WithContext(context.Background()).Create(behavior).Error; err != nil {
		return errors.NewDatabaseError("create user behavior", err)
	}

	return nil
}

// GetRoomRecommendations 获取房型推荐
func (s *RecommendationService) GetRoomRecommendations(userID *int64, limit int, excludeRoomID int64) ([]RecommendedRoom, error) {
	if limit <= 0 || limit > 12 {
		limit = 6
	}

	rooms, _, err := s.roomRepo.FindAvailable(1, 100)
	if err != nil {
		return nil, errors.NewDatabaseError("list available rooms", err)
	}
	if len(rooms) == 0 {
		return []RecommendedRoom{}, nil
	}

	roomMap := make(map[int64]models.Room, len(rooms))
	for _, room := range rooms {
		roomMap[int64(room.ID)] = room
	}

	interactionMap, popularityMap, err := s.buildInteractionMatrix()
	if err != nil {
		return nil, err
	}

	var profile *models.User
	if userID != nil && *userID > 0 {
		profile, _ = s.userRepo.FindByID(*userID)
	}

	stats := buildRoomFeatureStats(rooms, popularityMap)
	currentSignals := map[int64]float64{}
	hasHistory := false
	if userID != nil && *userID > 0 {
		currentSignals = interactionMap[*userID]
		hasHistory = len(currentSignals) > 0
	}

	preferenceSummary := buildUserPreferenceSummary(profile, currentSignals, roomMap)
	collaborativeScores := computeCollaborativeScores(userID, interactionMap)

	recommendations := make([]RecommendedRoom, 0, len(rooms))
	for _, room := range rooms {
		roomID := int64(room.ID)
		if excludeRoomID > 0 && roomID == excludeRoomID {
			continue
		}

		popularityScore := normalizeValue(popularityMap[roomID], 0, stats.MaxPopularity)
		collaborativeScore := collaborativeScores[roomID]
		neuralScore := computeNeuralScore(preferenceSummary, room, stats, popularityScore)
		finalScore := blendRecommendationScores(hasHistory, collaborativeScore, neuralScore, popularityScore)
		reason, tags := buildRecommendationReason(room, preferenceSummary, collaborativeScore, neuralScore, popularityScore)

		recommendations = append(recommendations, RecommendedRoom{
			Room:               room,
			Score:              roundFloat(finalScore, 4),
			CollaborativeScore: roundFloat(collaborativeScore, 4),
			NeuralScore:        roundFloat(neuralScore, 4),
			PopularityScore:    roundFloat(popularityScore, 4),
			Reason:             reason,
			Tags:               tags,
		})
	}

	sort.SliceStable(recommendations, func(i, j int) bool {
		if recommendations[i].Score == recommendations[j].Score {
			return recommendations[i].Room.Price < recommendations[j].Room.Price
		}
		return recommendations[i].Score > recommendations[j].Score
	})

	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

func (s *RecommendationService) buildInteractionMatrix() (map[int64]map[int64]float64, map[int64]float64, error) {
	interactionMap := make(map[int64]map[int64]float64)
	popularityMap := make(map[int64]float64)

	var behaviorSignals []userRoomSignal
	if err := s.db.Model(&models.UserBehavior{}).
		Select("user_id, room_id, SUM(weight) as weight").
		Group("user_id, room_id").
		Scan(&behaviorSignals).Error; err != nil {
		return nil, nil, errors.NewDatabaseError("list behavior signals", err)
	}

	var bookingSignals []userRoomSignal
	if err := s.db.Model(&models.Booking{}).
		Where("status <> ?", "cancelled").
		Select(`
			user_id,
			room_id,
			SUM(
				CASE
					WHEN payment_status = 'paid' THEN 8.0
					WHEN status = 'checkout' THEN 7.0
					WHEN status = 'checkin' THEN 6.5
					WHEN status = 'confirmed' THEN 6.0
					ELSE 5.0
				END
			) as weight
		`).
		Group("user_id, room_id").
		Scan(&bookingSignals).Error; err != nil {
		return nil, nil, errors.NewDatabaseError("list booking signals", err)
	}

	mergeSignal := func(signal userRoomSignal) {
		if _, exists := interactionMap[signal.UserID]; !exists {
			interactionMap[signal.UserID] = make(map[int64]float64)
		}
		interactionMap[signal.UserID][signal.RoomID] += signal.Weight
		popularityMap[signal.RoomID] += signal.Weight
	}

	for _, signal := range behaviorSignals {
		mergeSignal(signal)
	}
	for _, signal := range bookingSignals {
		mergeSignal(signal)
	}

	return interactionMap, popularityMap, nil
}

func recommendationBehaviorWeight(behaviorType string) (float64, bool) {
	switch normalizeBehaviorType(behaviorType) {
	case "view_room":
		return 1.0, true
	case "view_detail":
		return 1.8, true
	case "click_recommendation":
		return 2.2, true
	case "book_intent":
		return 3.6, true
	default:
		return 0, false
	}
}

func normalizeBehaviorType(behaviorType string) string {
	return strings.ToLower(strings.TrimSpace(behaviorType))
}

func computeCollaborativeScores(currentUserID *int64, interactionMap map[int64]map[int64]float64) map[int64]float64 {
	scores := make(map[int64]float64)
	if currentUserID == nil || *currentUserID <= 0 {
		return scores
	}

	currentVector := interactionMap[*currentUserID]
	if len(currentVector) == 0 {
		return scores
	}

	similaritySums := make(map[int64]float64)
	for otherUserID, otherVector := range interactionMap {
		if otherUserID == *currentUserID {
			continue
		}

		similarity := cosineSimilarity(currentVector, otherVector)
		if similarity <= 0 {
			continue
		}

		for roomID, weight := range otherVector {
			scores[roomID] += similarity * weight
			similaritySums[roomID] += similarity
		}
	}

	maxScore := 0.0
	for roomID, score := range scores {
		if similaritySums[roomID] > 0 {
			score /= similaritySums[roomID]
		}
		scores[roomID] = score
		if score > maxScore {
			maxScore = score
		}
	}

	if maxScore > 0 {
		for roomID, score := range scores {
			scores[roomID] = score / maxScore
		}
	}

	return scores
}

func buildUserPreferenceSummary(profile *models.User, userSignals map[int64]float64, roomMap map[int64]models.Room) userPreferenceSummary {
	summary := userPreferenceSummary{
		PreferredTypes: make(map[string]float64),
	}

	if profile != nil {
		summary.SpendLevel = clamp(profile.TotalSpend/5000, 0, 1)
		switch profile.Level {
		case "platinum":
			summary.LevelBoost = 0.18
		case "gold":
			summary.LevelBoost = 0.12
		case "silver":
			summary.LevelBoost = 0.07
		default:
			summary.LevelBoost = 0.02
		}
	}

	totalWeight := 0.0
	for roomID, weight := range userSignals {
		room, exists := roomMap[roomID]
		if !exists || weight <= 0 {
			continue
		}

		totalWeight += weight
		summary.PreferredTypes[room.RoomType] += weight
		summary.AvgPrice += room.Price * weight
		summary.AvgArea += room.Area * weight
		summary.AvgCapacity += float64(room.Capacity) * weight
	}

	if totalWeight > 0 {
		summary.HasHistory = true
		summary.AvgPrice /= totalWeight
		summary.AvgArea /= totalWeight
		summary.AvgCapacity /= totalWeight
		for roomType, weight := range summary.PreferredTypes {
			summary.PreferredTypes[roomType] = weight / totalWeight
		}
	}

	return summary
}

func buildRoomFeatureStats(rooms []models.Room, popularityMap map[int64]float64) roomFeatureStats {
	stats := roomFeatureStats{}
	if len(rooms) == 0 {
		return stats
	}

	stats.MinPrice = rooms[0].Price
	stats.MaxPrice = rooms[0].Price
	stats.MinArea = rooms[0].Area
	stats.MaxArea = rooms[0].Area
	stats.MinCapacity = float64(rooms[0].Capacity)
	stats.MaxCapacity = float64(rooms[0].Capacity)

	for _, room := range rooms {
		stats.MinPrice = math.Min(stats.MinPrice, room.Price)
		stats.MaxPrice = math.Max(stats.MaxPrice, room.Price)
		stats.MinArea = math.Min(stats.MinArea, room.Area)
		stats.MaxArea = math.Max(stats.MaxArea, room.Area)
		stats.MinCapacity = math.Min(stats.MinCapacity, float64(room.Capacity))
		stats.MaxCapacity = math.Max(stats.MaxCapacity, float64(room.Capacity))
		stats.MaxPopularity = math.Max(stats.MaxPopularity, popularityMap[int64(room.ID)])
	}

	return stats
}

func computeNeuralScore(summary userPreferenceSummary, room models.Room, stats roomFeatureStats, popularityScore float64) float64 {
	priceNorm := normalizeValue(room.Price, stats.MinPrice, stats.MaxPrice)
	areaNorm := normalizeValue(room.Area, stats.MinArea, stats.MaxArea)
	capacityNorm := normalizeValue(float64(room.Capacity), stats.MinCapacity, stats.MaxCapacity)
	discountNorm := clamp(room.GetDiscountRate()/100, 0, 1)
	typePreference := summary.PreferredTypes[room.RoomType]

	priceAffinity := 0.55
	areaAffinity := 0.55
	capacityAffinity := 0.55

	if summary.HasHistory {
		priceAffinity = 1 - math.Min(math.Abs(room.Price-summary.AvgPrice)/(math.Max(stats.MaxPrice-stats.MinPrice, 1)+1), 1)
		areaAffinity = 1 - math.Min(math.Abs(room.Area-summary.AvgArea)/(math.Max(stats.MaxArea-stats.MinArea, 1)+1), 1)
		capacityAffinity = 1 - math.Min(math.Abs(float64(room.Capacity)-summary.AvgCapacity)/(math.Max(stats.MaxCapacity-stats.MinCapacity, 1)+1), 1)
	}

	inputs := []float64{
		priceNorm,
		areaNorm,
		capacityNorm,
		discountNorm,
		popularityScore,
		typePreference,
		priceAffinity,
		areaAffinity,
		capacityAffinity,
		summary.SpendLevel,
		summary.LevelBoost,
	}

	hidden := []float64{
		relu(0.82*inputs[5] + 0.61*inputs[6] + 0.33*inputs[3] + 0.18*inputs[4] - 0.42*inputs[0] + 0.10),
		relu(0.76*inputs[4] + 0.44*inputs[1] + 0.39*inputs[2] + 0.21*inputs[9] - 0.14*inputs[0] + 0.06),
		relu(0.88*inputs[5] + 0.29*inputs[7] + 0.31*inputs[8] + 0.16*inputs[10] + 0.04),
		relu(0.71*inputs[6] + 0.53*inputs[7] + 0.25*inputs[4] + 0.17*inputs[3] + 0.08),
	}

	output := sigmoid(
		0.34*hidden[0] +
			0.27*hidden[1] +
			0.22*hidden[2] +
			0.25*hidden[3] -
			0.12,
	)

	return clamp(output, 0, 1)
}

func blendRecommendationScores(hasHistory bool, collaborativeScore, neuralScore, popularityScore float64) float64 {
	if hasHistory {
		return collaborativeScore*0.48 + neuralScore*0.37 + popularityScore*0.15
	}
	return neuralScore*0.68 + popularityScore*0.32
}

func buildRecommendationReason(room models.Room, summary userPreferenceSummary, collaborativeScore, neuralScore, popularityScore float64) (string, []string) {
	tags := make([]string, 0, 3)

	if collaborativeScore >= 0.45 {
		tags = append(tags, "协同偏好")
	}
	if summary.PreferredTypes[room.RoomType] >= 0.28 {
		tags = append(tags, "房型匹配")
	}
	if popularityScore >= 0.55 {
		tags = append(tags, "热门房型")
	}
	if room.GetDiscountRate() >= 15 {
		tags = append(tags, "高性价比")
	}
	if len(tags) == 0 {
		tags = append(tags, "智能推荐")
	}

	switch {
	case collaborativeScore >= 0.45:
		return "基于相似住客偏好与浏览行为综合推荐", tags
	case summary.PreferredTypes[room.RoomType] >= 0.28:
		return "与你近期关注的房型特征高度匹配", tags
	case popularityScore >= 0.55:
		return "近期预订热度较高，适合作为优先选择", tags
	case neuralScore >= 0.6:
		return "结合价格、面积与入住偏好进行神经打分推荐", tags
	default:
		return "基于房型特征与业务热度的综合推荐", tags
	}
}

func cosineSimilarity(a, b map[int64]float64) float64 {
	dot := 0.0
	normA := 0.0
	normB := 0.0

	for roomID, value := range a {
		normA += value * value
		if otherValue, exists := b[roomID]; exists {
			dot += value * otherValue
		}
	}
	for _, value := range b {
		normB += value * value
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func normalizeValue(value, minValue, maxValue float64) float64 {
	if maxValue <= minValue {
		return 0.5
	}
	return clamp((value-minValue)/(maxValue-minValue), 0, 1)
}

func relu(value float64) float64 {
	if value < 0 {
		return 0
	}
	return value
}

func sigmoid(value float64) float64 {
	return 1 / (1 + math.Exp(-value))
}

func clamp(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func roundFloat(value float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(value*pow) / pow
}
