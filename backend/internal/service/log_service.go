package service

import (
	"gohotel/internal/models"
	"gohotel/internal/repository"
)

// LogReportRequest 日志上报请求（支持批量）
type LogReportRequest struct {
	Logs []LogEntry `json:"logs" binding:"required,min=1"`
}

// LogEntry 单条日志
type LogEntry struct {
	Level   string `json:"level" binding:"required,oneof=debug info warn error"`
	Message string `json:"message" binding:"required"`
}

// LogListResponse Swagger 用日志列表响应结构
type LogListResponse []models.Log

// LogService 日志业务逻辑层
type LogService struct {
	logRepo *repository.LogRepository
}

// NewLogService 创建日志服务实例
func NewLogService(logRepo *repository.LogRepository) *LogService {
	return &LogService{logRepo: logRepo}
}

// ReportLogs 批量上报日志
func (s *LogService) ReportLogs(req *LogReportRequest) error {
	logs := make([]*models.Log, 0, len(req.Logs))
	for _, entry := range req.Logs {
		logs = append(logs, &models.Log{
			Level:   entry.Level,
			Message: entry.Message,
		})
	}
	return s.logRepo.BatchCreate(logs)
}

// CreateLog 创建单条日志
func (s *LogService) CreateLog(level, message string) error {
	log := &models.Log{
		Level:   level,
		Message: message,
	}
	return s.logRepo.Create(log)
}

// GetLogs 分页获取日志
func (s *LogService) GetLogs(page, pageSize int) ([]models.Log, int64, error) {
	return s.logRepo.GetLogs(page, pageSize)
}

// DeleteLog 删除日志
func (s *LogService) DeleteLog(id uint) error {
	return s.logRepo.Delete(id)
}

// GetLogByID 根据 ID 获取日志
func (s *LogService) GetLogByID(id uint) (*models.Log, error) {
	return s.logRepo.GetByID(id)
}
