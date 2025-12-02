package service

import (
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/pkg/utils"
)

// BannerService 活动横幅业务逻辑层
type BannerService struct {
	bannerRepo *repository.BannerRepository
	cosService *CosService
}

// NewBannerService 创建活动横幅服务实例
func NewBannerService(bannerRepo *repository.BannerRepository, cosService *CosService) *BannerService {
	return &BannerService{
		bannerRepo: bannerRepo,
		cosService: cosService,
	}
}

// CreateBannerRequest 创建活动横幅请求结构
type CreateBannerRequest struct {
	Title     string  `json:"title" binding:"required,min=1,max=100"`
	Subtitle  *string `json:"subtitle" binding:"omitempty,max=255"`
	ImageURL  string  `json:"image_url" binding:"required,max=500"`
	LinkURL   *string `json:"link_url" binding:"omitempty,max=500"`
	Status    string  `json:"status" binding:"omitempty,oneof=active inactive"`
	Sort      int     `json:"sort" binding:"omitempty,min=0"`
	StartTime *string `json:"start_time" binding:"omitempty"`
	EndTime   *string `json:"end_time" binding:"omitempty"`
}

// UpdateBannerRequest 更新活动横幅请求结构
type UpdateBannerRequest struct {
	Title     string  `json:"title" binding:"omitempty,min=1,max=100"`
	Subtitle  *string `json:"subtitle" binding:"omitempty,max=255"`
	ImageURL  string  `json:"image_url" binding:"omitempty,max=500"`
	LinkURL   *string `json:"link_url" binding:"omitempty,max=500"`
	Status    string  `json:"status" binding:"omitempty,oneof=active inactive"`
	Sort      int     `json:"sort" binding:"omitempty,min=0"`
	StartTime *string `json:"start_time" binding:"omitempty"`
	EndTime   *string `json:"end_time" binding:"omitempty"`
}

// CreateBanner 创建活动横幅
func (s *BannerService) CreateBanner(req *CreateBannerRequest) (*models.Banner, error) {
	// 设置默认值
	status := "active"
	if req.Status != "" {
		status = req.Status
	}

	// 创建Banner对象
	banner := &models.Banner{
		ID:       utils.JSONInt64(utils.GenID()),
		Title:    req.Title,
		Subtitle: req.Subtitle,
		ImageURL: req.ImageURL,
		LinkURL:  req.LinkURL,
		Status:   status,
		Sort:     req.Sort,
	}

	// 保存到数据库
	if err := s.bannerRepo.Create(banner); err != nil {
		return nil, err
	}

	return banner, nil
}

// GetBannerByID 根据ID获取活动横幅
func (s *BannerService) GetBannerByID(id int64) (*models.Banner, error) {
	return s.bannerRepo.FindByID(id)
}

// GetAllBanners 获取所有活动横幅（带分页）
func (s *BannerService) GetAllBanners(page, pageSize int) ([]models.Banner, int64, error) {
	return s.bannerRepo.FindAll(page, pageSize)
}

// GetActiveBanners 获取激活的活动横幅（前端展示用）
func (s *BannerService) GetActiveBanners() ([]models.Banner, error) {
	return s.bannerRepo.FindActive()
}

// UpdateBanner 更新活动横幅信息
func (s *BannerService) UpdateBanner(id int64, req *UpdateBannerRequest) (*models.Banner, error) {
	// 获取现有Banner
	banner, err := s.bannerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Title != "" {
		banner.Title = req.Title
	}
	if req.Subtitle != nil {
		banner.Subtitle = req.Subtitle
	}
	if req.ImageURL != "" {
		banner.ImageURL = req.ImageURL
	}
	if req.LinkURL != nil {
		banner.LinkURL = req.LinkURL
	}
	if req.Status != "" {
		banner.Status = req.Status
	}
	if req.Sort != 0 {
		banner.Sort = req.Sort
	}

	// 保存更新
	if err := s.bannerRepo.Update(banner); err != nil {
		return nil, err
	}

	return banner, nil
}

// UpdateBannerStatus 更新活动横幅状态
func (s *BannerService) UpdateBannerStatus(id int64, status string) error {
	return s.bannerRepo.UpdateStatus(id, status)
}

// DeleteBanner 删除活动横幅
func (s *BannerService) DeleteBanner(id int64) error {
	return s.bannerRepo.Delete(id)
}
