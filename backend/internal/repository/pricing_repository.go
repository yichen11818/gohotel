package repository

import (
	"context"
	"gohotel/internal/models"
	"time"

	"gorm.io/gorm"
)

type PricingRepository interface {
	GetActiveRules(ctx context.Context, roomType string, date time.Time) ([]models.PricingRule, error)
	CreateRule(ctx context.Context, rule *models.PricingRule) error
	ListRules(ctx context.Context) ([]models.PricingRule, error)
}

type pricingRepository struct {
	db *gorm.DB
}

func NewPricingRepository(db *gorm.DB) PricingRepository {
	return &pricingRepository{db: db}
}

func (r *pricingRepository) GetActiveRules(ctx context.Context, roomType string, date time.Time) ([]models.PricingRule, error) {
	var rules []models.PricingRule
	err := r.db.WithContext(ctx).
		Where("start_date <= ? AND end_date >= ?", date, date).
		Where("(room_type = ? OR room_type = '')", roomType).
		Order("priority DESC").
		Find(&rules).Error
	return rules, err
}

func (r *pricingRepository) CreateRule(ctx context.Context, rule *models.PricingRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

func (r *pricingRepository) ListRules(ctx context.Context) ([]models.PricingRule, error) {
	var rules []models.PricingRule
	err := r.db.WithContext(ctx).Find(&rules).Error
	return rules, err
}
