package service

import (
	"context"
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"time"
)

type InventoryService interface {
	CheckAvailability(ctx context.Context, roomType string, startDate, endDate time.Time) (bool, error)
	GetDailyPrice(ctx context.Context, roomType string, date time.Time) (float64, error)
	UpdateInventory(ctx context.Context, roomType string, startDate, endDate time.Time, delta int) error
	InitInventory(ctx context.Context, roomType string, totalCount int, price float64, days int) error
	GetInventoryGrid(ctx context.Context, startDate, endDate time.Time) (map[string][]models.RoomInventory, error)
}

// InventoryGridResponse Swagger 用库存网格响应结构
type InventoryGridResponse map[string][]models.RoomInventory

type inventoryService struct {
	repo        repository.InventoryRepository
	pricingRepo repository.PricingRepository
}

func NewInventoryService(repo repository.InventoryRepository, pricingRepo repository.PricingRepository) InventoryService {
	return &inventoryService{
		repo:        repo,
		pricingRepo: pricingRepo,
	}
}

func (s *inventoryService) CheckAvailability(ctx context.Context, roomType string, startDate, endDate time.Time) (bool, error) {
	inventories, err := s.repo.GetByDateRange(ctx, roomType, startDate, endDate.AddDate(0, 0, -1))
	if err != nil {
		return false, err
	}

	days := int(endDate.Sub(startDate).Hours() / 24)
	if len(inventories) < days {
		return false, nil // Missing inventory records for some days
	}

	for _, inv := range inventories {
		if inv.TotalCount-inv.BookedCount <= 0 {
			return false, nil
		}
	}
	return true, nil
}

func (s *inventoryService) GetDailyPrice(ctx context.Context, roomType string, date time.Time) (float64, error) {
	inv, err := s.repo.GetByRoomTypeAndDate(ctx, roomType, date)
	if err != nil {
		return 0, err
	}

	price := inv.Price
	rules, _ := s.pricingRepo.GetActiveRules(ctx, roomType, date)
	for _, rule := range rules {
		if rule.IsPercent {
			price = price * (1 + rule.Adjustment/100)
		} else {
			price = price + rule.Adjustment
		}
	}
	return price, nil
}

func (s *inventoryService) UpdateInventory(ctx context.Context, roomType string, startDate, endDate time.Time, delta int) error {
	inventories, err := s.repo.GetByDateRange(ctx, roomType, startDate, endDate.AddDate(0, 0, -1))
	if err != nil {
		return err
	}

	for _, inv := range inventories {
		if err := s.repo.UpdateBookedCount(ctx, inv.ID, delta); err != nil {
			return err
		}
	}
	return nil
}

func (s *inventoryService) InitInventory(ctx context.Context, roomType string, totalCount int, price float64, days int) error {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var inventories []models.RoomInventory
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i)
		inventories = append(inventories, models.RoomInventory{
			RoomType:   roomType,
			Date:       date,
			TotalCount: totalCount,
			Price:      price,
		})
	}
	return s.repo.BatchCreate(ctx, inventories)
}

func (s *inventoryService) GetInventoryGrid(ctx context.Context, startDate, endDate time.Time) (map[string][]models.RoomInventory, error) {
	inventories, err := s.repo.GetByDateRange(ctx, "", startDate, endDate)
	if err != nil {
		return nil, err
	}

	grid := make(map[string][]models.RoomInventory)
	for _, inv := range inventories {
		grid[inv.RoomType] = append(grid[inv.RoomType], inv)
	}
	return grid, nil
}
