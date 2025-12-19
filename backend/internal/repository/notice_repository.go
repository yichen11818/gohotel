package repository

import (
	"gohotel/internal/models"

	"gorm.io/gorm"
)

// NoticeRepository 公告数据访问层
type NoticeRepository struct {
	db *gorm.DB
}

// NewNoticeRepository 创建公告仓库实例
func NewNoticeRepository(db *gorm.DB) *NoticeRepository { return &NoticeRepository{db: db} }

// Create 创建公告横幅
func (r *NoticeRepository) Create(notice *models.Notice) error { return r.db.Create(notice).Error }

// FindByID 根据ID查找公告
func (r *NoticeRepository) FindByID(id int64) (*models.Notice, error) {
	var notice models.Notice
	err := r.db.First(&notice, id).Error
	if err != nil {
		return nil, err
	}
	return &notice, nil
}

// FindAll 查找所有活动横幅 (带分页)
func (r *NoticeRepository) FindAll(page, pageSize int) ([]models.Notice, int64, error) {
	var notices []models.Notice
	var total int64

	query := r.db.Model(&models.Notice{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("sort ASC,created_at DESC").Find(&notices).Error
	return notices, total, err
}

// FindActive 查找激活的活动横幅（前端展示用）
func (r *NoticeRepository) FindActive() ([]models.Notice, error) {
	var notices []models.Notice
	err := r.db.Where("status = ?", "active").
		Order("sort ASC, created_at DESC").Find(&notices).Error
	return notices, err
}

// Update 更新活动横幅信息
func (r *NoticeRepository) Update(Notice *models.Notice) error {
	return r.db.Save(Notice).Error
}

// UpdateStatus 更新活动横幅状态
func (r *NoticeRepository) UpdateStatus(id int64, status string) error {
	return r.db.Model(&models.Notice{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateSort 更新活动横幅排序
func (r *NoticeRepository) UpdateSort(id int64, sort int) error {
	return r.db.Model(&models.Notice{}).Where("id = ?", id).Update("sort", sort).Error
}

// Delete 删除活动横幅
func (r *NoticeRepository) Delete(id int64) error {
	return r.db.Delete(&models.Notice{}, id).Error
}
