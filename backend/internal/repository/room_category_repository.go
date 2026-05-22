package repository

import (
	"gohotel/internal/models"

	"gorm.io/gorm"
)

// RoomCategoryRepository 房型分类仓库接口
type RoomCategoryRepository interface {
	Create(category *models.RoomCategory) error
	FindByID(id uint) (*models.RoomCategory, error)
	FindByName(name string) (*models.RoomCategory, error)
	FindAll(page, pageSize int) ([]models.RoomCategory, int64, error)
	Update(category *models.RoomCategory) error
	Delete(id uint) error
}

type roomCategoryRepository struct {
	db *gorm.DB
}

func NewRoomCategoryRepository(db *gorm.DB) RoomCategoryRepository {
	return &roomCategoryRepository{db: db}
}

func (r *roomCategoryRepository) Create(category *models.RoomCategory) error {
	return r.db.Create(category).Error
}

func (r *roomCategoryRepository) FindByID(id uint) (*models.RoomCategory, error) {
	var category models.RoomCategory
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *roomCategoryRepository) FindByName(name string) (*models.RoomCategory, error) {
	var category models.RoomCategory
	if err := r.db.Where("name = ?", name).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *roomCategoryRepository) FindAll(page, pageSize int) ([]models.RoomCategory, int64, error) {
	var categories []models.RoomCategory
	var total int64

	query := r.db.Model(&models.RoomCategory{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at ASC").Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *roomCategoryRepository) Update(category *models.RoomCategory) error {
	return r.db.Save(category).Error
}

func (r *roomCategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.RoomCategory{}, id).Error
}
