package repository

import (
	"gohotel/internal/models"

	"gorm.io/gorm"
)

// FacilityRepository 设施数据访问层
type FacilityRepository struct {
	db *gorm.DB
}

// NewFacilityRepository 创建设施仓库实例
func NewFacilityRepository(db *gorm.DB) *FacilityRepository {
	return &FacilityRepository{db: db}
}

// Create 创建设施
func (r *FacilityRepository) Create(facility *models.Facility) error {
	return r.db.Create(facility).Error
}

// FindByID 根据 ID 查找设施
func (r *FacilityRepository) FindByID(id uint) (*models.Facility, error) {
	var facility models.Facility
	err := r.db.First(&facility, id).Error
	if err != nil {
		return nil, err
	}
	return &facility, nil
}

// Update 更新设施
func (r *FacilityRepository) Update(facility *models.Facility) error {
	return r.db.Save(facility).Error
}

// Delete 删除设施
func (r *FacilityRepository) Delete(id uint) error {
	return r.db.Delete(&models.Facility{}, id).Error
}

// FindAll 查询所有设施（分页）
func (r *FacilityRepository) FindAll(page, pageSize int) ([]models.Facility, int64, error) {
	var facilities []models.Facility
	var total int64

	if err := r.db.Model(&models.Facility{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Order("floor, type").Find(&facilities).Error
	return facilities, total, err
}

// FindByFloor 按楼层查询设施
func (r *FacilityRepository) FindByFloor(floor int) ([]models.Facility, error) {
	var facilities []models.Facility
	err := r.db.Where("floor = ?", floor).Order("type").Find(&facilities).Error
	return facilities, err
}

// BatchUpdate 批量更新设施位置
func (r *FacilityRepository) BatchUpdate(facilities []models.Facility) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, f := range facilities {
			if err := tx.Model(&models.Facility{}).Where("id = ?", f.ID).Updates(map[string]interface{}{
				"left":     f.Left,
				"top":      f.Top,
				"width":    f.Width,
				"height":   f.Height,
				"rotation": f.Rotation,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
