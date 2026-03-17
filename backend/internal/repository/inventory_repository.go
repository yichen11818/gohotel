package repository

import (
	"context"
	"gohotel/internal/models"
	"time"

	"gorm.io/gorm"
)

type InventoryRepository interface {
	GetByRoomTypeAndDate(ctx context.Context, roomType string, date time.Time) (*models.RoomInventory, error)
	GetByDateRange(ctx context.Context, roomType string, startDate, endDate time.Time) ([]models.RoomInventory, error)
	UpdateBookedCount(ctx context.Context, id uint, delta int) error
	Create(ctx context.Context, inventory *models.RoomInventory) error
	BatchCreate(ctx context.Context, inventories []models.RoomInventory) error
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) GetByRoomTypeAndDate(ctx context.Context, roomType string, date time.Time) (*models.RoomInventory, error) {
	var inventory models.RoomInventory
	err := r.db.WithContext(ctx).Where("room_type = ? AND date = ?", roomType, date.Format("2006-01-02")).First(&inventory).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) GetByDateRange(ctx context.Context, roomType string, startDate, endDate time.Time) ([]models.RoomInventory, error) {
	var inventories []models.RoomInventory
	query := r.db.WithContext(ctx).Where("date >= ? AND date <= ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if roomType != "" {
		query = query.Where("room_type = ?", roomType)
	}
	err := query.Order("date ASC").Find(&inventories).Error
	return inventories, err
}

func (r *inventoryRepository) UpdateBookedCount(ctx context.Context, id uint, delta int) error {
	return r.db.WithContext(ctx).Model(&models.RoomInventory{}).Where("id = ?", id).
		UpdateColumn("booked_count", gorm.Expr("booked_count + ?", delta)).Error
}

func (r *inventoryRepository) Create(ctx context.Context, inventory *models.RoomInventory) error {
	return r.db.WithContext(ctx).Create(inventory).Error
}

func (r *inventoryRepository) BatchCreate(ctx context.Context, inventories []models.RoomInventory) error {
	return r.db.WithContext(ctx).Create(&inventories).Error
}
