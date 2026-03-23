package service

import (
	"context"
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"time"
)

type WorkOrderService interface {
	// Maintenance
	CreateRepairRequest(ctx context.Context, roomID int64, repairType, desc string, reporterID int64) error
	CompleteRepair(ctx context.Context, id uint, workerID int64, remark string) error
	ListMaintenance(ctx context.Context, status string) ([]models.Maintenance, error)

	// Housekeeping
	CreateCleaningTask(ctx context.Context, roomID int64, cleaningType string) error
	AssignStaff(ctx context.Context, id uint, staffID int64) error
	CompleteCleaning(ctx context.Context, id uint) error
	ListHousekeeping(ctx context.Context, status string) ([]models.Housekeeping, error)
}

// MaintenanceListResponse Swagger 用维修列表响应结构
type MaintenanceListResponse []models.Maintenance

// HousekeepingListResponse Swagger 用清洁列表响应结构
type HousekeepingListResponse []models.Housekeeping

type workOrderService struct {
	repo     repository.WorkOrderRepository
	roomRepo repository.RoomRepository
}

func NewWorkOrderService(repo repository.WorkOrderRepository, roomRepo repository.RoomRepository) WorkOrderService {
	return &workOrderService{
		repo:     repo,
		roomRepo: roomRepo,
	}
}

func (s *workOrderService) CreateRepairRequest(ctx context.Context, roomID int64, repairType, desc string, reporterID int64) error {
	m := &models.Maintenance{
		RoomID:      roomID,
		Type:        repairType,
		Description: desc,
		Status:      "pending",
		ReporterID:  reporterID,
	}
	if err := s.repo.CreateMaintenance(ctx, m); err != nil {
		return err
	}
	// Update room status to maintenance
	return s.roomRepo.UpdateStatus(ctx, roomID, "maintenance")
}

func (s *workOrderService) CompleteRepair(ctx context.Context, id uint, workerID int64, remark string) error {
	m, err := s.repo.GetMaintenanceByID(ctx, id)
	if err != nil {
		return err
	}
	m.Status = "completed"
	m.WorkerID = workerID
	m.Remark = remark
	if err := s.repo.UpdateMaintenance(ctx, m); err != nil {
		return err
	}
	// Restore room status to available (but dirty)
	_ = s.roomRepo.UpdateStatus(ctx, m.RoomID, "available")
	return s.roomRepo.UpdateCleanStatus(ctx, m.RoomID, "dirty")
}

func (s *workOrderService) ListMaintenance(ctx context.Context, status string) ([]models.Maintenance, error) {
	filter := make(map[string]interface{})
	if status != "" {
		filter["status"] = status
	}
	return s.repo.ListMaintenance(ctx, filter)
}

func (s *workOrderService) CreateCleaningTask(ctx context.Context, roomID int64, cleaningType string) error {
	h := &models.Housekeeping{
		RoomID: roomID,
		Type:   cleaningType,
		Status: "pending",
	}
	return s.repo.CreateHousekeeping(ctx, h)
}

func (s *workOrderService) AssignStaff(ctx context.Context, id uint, staffID int64) error {
	h, err := s.repo.GetHousekeepingByID(ctx, id)
	if err != nil {
		return err
	}
	h.StaffID = staffID
	h.Status = "in_progress"
	now := time.Now()
	h.StartTime = &now
	return s.repo.UpdateHousekeeping(ctx, h)
}

func (s *workOrderService) CompleteCleaning(ctx context.Context, id uint) error {
	h, err := s.repo.GetHousekeepingByID(ctx, id)
	if err != nil {
		return err
	}
	h.Status = "completed"
	now := time.Now()
	h.EndTime = &now
	if err := s.repo.UpdateHousekeeping(ctx, h); err != nil {
		return err
	}
	return s.roomRepo.UpdateCleanStatus(ctx, h.RoomID, "clean")
}

func (s *workOrderService) ListHousekeeping(ctx context.Context, status string) ([]models.Housekeeping, error) {
	filter := make(map[string]interface{})
	if status != "" {
		filter["status"] = status
	}
	return s.repo.ListHousekeeping(ctx, filter)
}
