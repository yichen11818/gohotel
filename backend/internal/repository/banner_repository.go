package repository

import (
	"gohotel/internal/models"

	"gorm.io/gorm"
)

// BannerRepository 活动横幅数据访问层
type BannerRepository struct {
	db *gorm.DB
}

// NewBannerRepository 创建活动横幅仓库实例
func NewBannerRepository(db *gorm.DB) *BannerRepository {
	return &BannerRepository{db: db}
}

// Create 创建活动横幅
func (r *BannerRepository) Create(banner *models.Banner) error {
	return r.db.Create(banner).Error
}

// FindByID 根据ID查找活动横幅
func (r *BannerRepository) FindByID(id int64) (*models.Banner, error) {
	var banner models.Banner
	err := r.db.First(&banner, id).Error
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

// FindAll 查找所有活动横幅（带分页）
func (r *BannerRepository) FindAll(page, pageSize int) ([]models.Banner, int64, error) {
	var banners []models.Banner
	var total int64

	query := r.db.Model(&models.Banner{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("sort ASC, created_at DESC").Find(&banners).Error
	return banners, total, err
}

// FindActive 查找激活的活动横幅（前端展示用）
func (r *BannerRepository) FindActive() ([]models.Banner, error) {
	var banners []models.Banner
	err := r.db.Where("status = ?", "active").
		Order("sort ASC, created_at DESC").Find(&banners).Error
	return banners, err
}

// Update 更新活动横幅信息
func (r *BannerRepository) Update(banner *models.Banner) error {
	return r.db.Save(banner).Error
}

// UpdateStatus 更新活动横幅状态
func (r *BannerRepository) UpdateStatus(id int64, status string) error {
	return r.db.Model(&models.Banner{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateSort 更新活动横幅排序
func (r *BannerRepository) UpdateSort(id int64, sort int) error {
	return r.db.Model(&models.Banner{}).Where("id = ?", id).Update("sort", sort).Error
}

// Delete 删除活动横幅
func (r *BannerRepository) Delete(id int64) error {
	return r.db.Delete(&models.Banner{}, id).Error
}
