package repository

import (
	"context"
	"gohotel/internal/models"

	"gorm.io/gorm"
)

type WorkOrderRepository interface {
	CreateMaintenance(ctx context.Context, m *models.Maintenance) error
	UpdateMaintenance(ctx context.Context, m *models.Maintenance) error
	GetMaintenanceByID(ctx context.Context, id uint) (*models.Maintenance, error)
	ListMaintenance(ctx context.Context, filter map[string]interface{}) ([]models.Maintenance, error)

	CreateHousekeeping(ctx context.Context, h *models.Housekeeping) error
	UpdateHousekeeping(ctx context.Context, h *models.Housekeeping) error
	GetHousekeepingByID(ctx context.Context, id uint) (*models.Housekeeping, error)
	ListHousekeeping(ctx context.Context, filter map[string]interface{}) ([]models.Housekeeping, error)
}

type workOrderRepository struct {
	db *gorm.DB
}

func NewWorkOrderRepository(db *gorm.DB) WorkOrderRepository {
	return &workOrderRepository{db: db}
}

func (r *workOrderRepository) CreateMaintenance(ctx context.Context, m *models.Maintenance) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *workOrderRepository) UpdateMaintenance(ctx context.Context, m *models.Maintenance) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *workOrderRepository) GetMaintenanceByID(ctx context.Context, id uint) (*models.Maintenance, error) {
	var m models.Maintenance
	err := r.db.WithContext(ctx).Preload("Room").First(&m, id).Error
	return &m, err
}

func (r *workOrderRepository) ListMaintenance(ctx context.Context, filter map[string]interface{}) ([]models.Maintenance, error) {
	var list []models.Maintenance
	err := r.db.WithContext(ctx).Where(filter).Preload("Room").Find(&list).Error
	return list, err
}

func (r *workOrderRepository) CreateHousekeeping(ctx context.Context, h *models.Housekeeping) error {
	return r.db.WithContext(ctx).Create(h).Error
}

func (r *workOrderRepository) UpdateHousekeeping(ctx context.Context, h *models.Housekeeping) error {
	return r.db.WithContext(ctx).Save(h).Error
}

func (r *workOrderRepository) GetHousekeepingByID(ctx context.Context, id uint) (*models.Housekeeping, error) {
	var h models.Housekeeping
	err := r.db.WithContext(ctx).Preload("Room").First(&h, id).Error
	return &h, err
}

func (r *workOrderRepository) ListHousekeeping(ctx context.Context, filter map[string]interface{}) ([]models.Housekeeping, error) {
	var list []models.Housekeeping
	err := r.db.WithContext(ctx).Where(filter).Preload("Room").Find(&list).Error
	return list, err
}
