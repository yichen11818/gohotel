package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/pkg/errors"

	"gorm.io/gorm"
)

// RoomCategoryService 房型分类服务
type RoomCategoryService struct {
	db               *gorm.DB
	roomCategoryRepo repository.RoomCategoryRepository
}

// CreateRoomCategoryRequest 创建房型分类请求
type CreateRoomCategoryRequest struct {
	Name        string `json:"name" binding:"required"` // 房型名称
	Description string `json:"description"`             // 房型描述
	Facilities  string `json:"facilities"`              // 设施 JSON 数组字符串
	Images      string `json:"images"`                  // 预览图 JSON 数组字符串
}

// UpdateRoomCategoryRequest 更新房型分类请求
type UpdateRoomCategoryRequest struct {
	Name        *string `json:"name"`        // 房型名称
	Description *string `json:"description"` // 房型描述
	Facilities  *string `json:"facilities"`  // 设施 JSON 数组字符串
	Images      *string `json:"images"`      // 预览图 JSON 数组字符串
}

func NewRoomCategoryService(db *gorm.DB, roomCategoryRepo repository.RoomCategoryRepository) *RoomCategoryService {
	return &RoomCategoryService{
		db:               db,
		roomCategoryRepo: roomCategoryRepo,
	}
}

func (s *RoomCategoryService) CreateRoomCategory(req *CreateRoomCategoryRequest) (*models.RoomCategory, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.NewBadRequestError("房型名称不能为空")
	}

	facilities, err := normalizeJSONArrayString(req.Facilities, "设施")
	if err != nil {
		return nil, err
	}
	images, err := normalizeJSONArrayString(req.Images, "预览图")
	if err != nil {
		return nil, err
	}

	existing, err := s.roomCategoryRepo.FindByName(name)
	if err == nil && existing != nil {
		return nil, errors.NewConflictError("房型分类已存在")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.NewDatabaseError("find room category by name", err)
	}

	category := &models.RoomCategory{
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		Facilities:  facilities,
		Images:      images,
	}

	if err := s.roomCategoryRepo.Create(category); err != nil {
		return nil, errors.NewDatabaseError("create room category", err)
	}

	if err := s.syncRoomCategoryFields(category.Name, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *RoomCategoryService) UpdateRoomCategory(id uint, req *UpdateRoomCategoryRequest) (*models.RoomCategory, error) {
	category, err := s.roomCategoryRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("房型分类不存在")
		}
		return nil, errors.NewDatabaseError("find room category", err)
	}

	oldName := category.Name
	newName := category.Name
	if req.Name != nil {
		newName = strings.TrimSpace(*req.Name)
		if newName == "" {
			return nil, errors.NewBadRequestError("房型名称不能为空")
		}
	}

	facilities := category.Facilities
	if req.Facilities != nil {
		facilities, err = normalizeJSONArrayString(*req.Facilities, "设施")
		if err != nil {
			return nil, err
		}
	}

	images := category.Images
	if req.Images != nil {
		images, err = normalizeJSONArrayString(*req.Images, "预览图")
		if err != nil {
			return nil, err
		}
	}

	if newName != oldName {
		existing, findErr := s.roomCategoryRepo.FindByName(newName)
		if findErr == nil && existing != nil && existing.ID != category.ID {
			return nil, errors.NewConflictError("房型分类名称已存在")
		}
		if findErr != nil && findErr != gorm.ErrRecordNotFound {
			return nil, errors.NewDatabaseError("find room category by name", findErr)
		}
	}

	if req.Description != nil {
		category.Description = strings.TrimSpace(*req.Description)
	}
	category.Name = newName
	category.Facilities = facilities
	category.Images = images

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(category).Error; err != nil {
			return err
		}

		if oldName != newName {
			if err := tx.Model(&models.Room{}).Where("room_type = ?", oldName).Update("room_type", newName).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.RoomInventory{}).Where("room_type = ?", oldName).Update("room_type", newName).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.PricingRule{}).Where("room_type = ?", oldName).Update("room_type", newName).Error; err != nil {
				return err
			}
		}

		if err := s.syncRoomCategoryFieldsTx(tx, newName, category); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.NewDatabaseError("update room category", err)
	}

	return category, nil
}

func (s *RoomCategoryService) DeleteRoomCategory(id uint) error {
	category, err := s.roomCategoryRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("房型分类不存在")
		}
		return errors.NewDatabaseError("find room category", err)
	}

	var roomCount int64
	if err := s.db.Model(&models.Room{}).Where("room_type = ?", category.Name).Count(&roomCount).Error; err != nil {
		return errors.NewDatabaseError("count rooms by category", err)
	}
	if roomCount > 0 {
		return errors.NewConflictError(fmt.Sprintf("该房型已关联 %d 个房间，请先调整房间分类后再删除", roomCount))
	}

	var inventoryCount int64
	if err := s.db.Model(&models.RoomInventory{}).Where("room_type = ?", category.Name).Count(&inventoryCount).Error; err != nil {
		return errors.NewDatabaseError("count inventories by category", err)
	}
	if inventoryCount > 0 {
		return errors.NewConflictError("该房型已存在库存配置，请先处理库存后再删除")
	}

	var pricingCount int64
	if err := s.db.Model(&models.PricingRule{}).Where("room_type = ?", category.Name).Count(&pricingCount).Error; err != nil {
		return errors.NewDatabaseError("count pricing rules by category", err)
	}
	if pricingCount > 0 {
		return errors.NewConflictError("该房型已存在定价规则，请先处理规则后再删除")
	}

	if err := s.roomCategoryRepo.Delete(id); err != nil {
		return errors.NewDatabaseError("delete room category", err)
	}
	return nil
}

func (s *RoomCategoryService) GetRoomCategoryByID(id uint) (*models.RoomCategory, error) {
	category, err := s.roomCategoryRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("房型分类不存在")
		}
		return nil, errors.NewDatabaseError("find room category", err)
	}

	if err := s.attachUsageCounts([]*models.RoomCategory{category}); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *RoomCategoryService) ListRoomCategories(page, pageSize int) ([]models.RoomCategory, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 20
	}

	categories, total, err := s.roomCategoryRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("list room categories", err)
	}

	categoryRefs := make([]*models.RoomCategory, 0, len(categories))
	for i := range categories {
		categoryRefs = append(categoryRefs, &categories[i])
	}
	if err := s.attachUsageCounts(categoryRefs); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (s *RoomCategoryService) syncRoomCategoryFields(roomType string, category *models.RoomCategory) error {
	return s.syncRoomCategoryFieldsTx(s.db, roomType, category)
}

func (s *RoomCategoryService) syncRoomCategoryFieldsTx(tx *gorm.DB, roomType string, category *models.RoomCategory) error {
	updates := map[string]interface{}{
		"description": category.Description,
		"facilities":  category.Facilities,
		"images":      category.Images,
	}
	return tx.Model(&models.Room{}).Where("room_type = ?", roomType).Updates(updates).Error
}

func (s *RoomCategoryService) attachUsageCounts(categories []*models.RoomCategory) error {
	if len(categories) == 0 {
		return nil
	}

	names := make([]string, 0, len(categories))
	categoryMap := make(map[string]*models.RoomCategory, len(categories))
	for _, category := range categories {
		names = append(names, category.Name)
		categoryMap[category.Name] = category
	}

	type usageRow struct {
		RoomType string `gorm:"column:room_type"`
		Count    int64  `gorm:"column:count"`
	}

	var usage []usageRow
	if err := s.db.Model(&models.Room{}).
		Select("room_type, COUNT(*) as count").
		Where("room_type IN ?", names).
		Group("room_type").
		Scan(&usage).Error; err != nil {
		return errors.NewDatabaseError("count room category usage", err)
	}

	for _, row := range usage {
		if category := categoryMap[row.RoomType]; category != nil {
			category.UsageCount = row.Count
		}
	}

	return nil
}

func normalizeJSONArrayString(raw string, fieldName string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "[]", nil
	}

	var parsed []string
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return "", errors.NewBadRequestError(fmt.Sprintf("%s必须是 JSON 数组字符串", fieldName))
	}

	normalized, err := json.Marshal(parsed)
	if err != nil {
		return "", errors.NewBadRequestError(fmt.Sprintf("%s格式错误", fieldName))
	}

	return string(normalized), nil
}
