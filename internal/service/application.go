package service

import (
	"errors"
	"openmdm/internal/model"
	"openmdm/internal/repository"
	"time"
)

var (
	ErrAppNotFound       = errors.New("application not found")
	ErrAppAlreadyExists  = errors.New("application already exists")
	ErrAppDeviceNotFound = errors.New("device not found")
)

type ApplicationService struct {
	appRepo     *repository.ApplicationRepository
	deviceRepo  *repository.DeviceRepository
}

func NewApplicationService(appRepo *repository.ApplicationRepository, deviceRepo *repository.DeviceRepository) *ApplicationService {
	return &ApplicationService{
		appRepo:    appRepo,
		deviceRepo: deviceRepo,
	}
}

type CreateAppRequest struct {
	Name        string          `json:"name"`
	BundleID    string          `json:"bundle_id"`
	Platform    model.AppPlatform `json:"platform"`
	Version     string          `json:"version"`
	FileURL     string          `json:"file_url"`
	FileSize    int64           `json:"file_size"`
	IconURL     string          `json:"icon_url"`
	Description string          `json:"description"`
}

func (s *ApplicationService) CreateApp(req *CreateAppRequest) (*model.Application, error) {
	app := &model.Application{
		Name:        req.Name,
		BundleID:    req.BundleID,
		Platform:    req.Platform,
		Version:     req.Version,
		FileURL:     req.FileURL,
		FileSize:    req.FileSize,
		IconURL:     req.IconURL,
		Description: req.Description,
		Status:      model.AppStatusActive,
	}

	if err := s.appRepo.Create(app); err != nil {
		return nil, err
	}

	return app, nil
}

type ListAppsFilter struct {
	Platform model.AppPlatform
	Limit    int
	Offset   int
}

func (s *ApplicationService) ListApps(filter *ListAppsFilter) ([]model.Application, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	var apps []model.Application
	var total int64
	var err error

	if filter.Platform != "" {
		apps, err = s.appRepo.ListByPlatform(filter.Platform, filter.Limit, filter.Offset)
		if err != nil {
			return nil, 0, err
		}
		total, err = s.appRepo.CountByPlatform(filter.Platform)
	} else {
		apps, err = s.appRepo.List(filter.Limit, filter.Offset)
		if err != nil {
			return nil, 0, err
		}
		total, err = s.appRepo.Count()
	}

	if err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

func (s *ApplicationService) GetApp(id uint) (*model.Application, error) {
	app, err := s.appRepo.GetByID(id)
	if err != nil {
		return nil, ErrAppNotFound
	}
	return app, nil
}

type UpdateAppRequest struct {
	Name        string           `json:"name"`
	BundleID    string           `json:"bundle_id"`
	Platform    model.AppPlatform `json:"platform"`
	Version     string           `json:"version"`
	FileURL     string           `json:"file_url"`
	FileSize    int64            `json:"file_size"`
	IconURL     string           `json:"icon_url"`
	Description string           `json:"description"`
	Status      model.AppStatus  `json:"status"`
}

func (s *ApplicationService) UpdateApp(id uint, req *UpdateAppRequest) (*model.Application, error) {
	app, err := s.appRepo.GetByID(id)
	if err != nil {
		return nil, ErrAppNotFound
	}

	if req.Name != "" {
		app.Name = req.Name
	}
	if req.BundleID != "" {
		app.BundleID = req.BundleID
	}
	if req.Platform != "" {
		app.Platform = req.Platform
	}
	if req.Version != "" {
		app.Version = req.Version
	}
	if req.FileURL != "" {
		app.FileURL = req.FileURL
	}
	if req.FileSize > 0 {
		app.FileSize = req.FileSize
	}
	if req.IconURL != "" {
		app.IconURL = req.IconURL
	}
	if req.Description != "" {
		app.Description = req.Description
	}
	if req.Status != "" {
		app.Status = req.Status
	}

	if err := s.appRepo.Update(app); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *ApplicationService) DeleteApp(id uint) error {
	_, err := s.appRepo.GetByID(id)
	if err != nil {
		return ErrAppNotFound
	}
	return s.appRepo.Delete(id)
}

type InstallAppRequest struct {
	ApplicationID uint `json:"application_id"`
}

func (s *ApplicationService) InstallApp(deviceID, appID uint) (*model.DeviceApplication, error) {
	// 检查设备是否存在
	_, err := s.deviceRepo.GetByID(uint64(deviceID))
	if err != nil {
		return nil, ErrAppDeviceNotFound
	}

	// 检查应用是否存在
	_, err = s.appRepo.GetByID(appID)
	if err != nil {
		return nil, ErrAppNotFound
	}

	// 检查是否已安装
	existing, err := s.appRepo.GetDeviceApp(uint64(deviceID), appID)
	if err == nil && existing != nil {
		// 已存在则更新状态为pending
		now := time.Now()
		existing.InstallStatus = model.InstallStatusPending
		existing.InstalledAt = &now
		if err := s.appRepo.UpdateDeviceApp(existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// 创建新的安装记录
	now := time.Now()
	deviceApp := &model.DeviceApplication{
		DeviceID:      deviceID,
		ApplicationID: appID,
		InstallStatus: model.InstallStatusPending,
		InstalledAt:   &now,
	}

	if err := s.appRepo.CreateDeviceApp(deviceApp); err != nil {
		return nil, err
	}

	return deviceApp, nil
}

func (s *ApplicationService) ListDeviceApps(deviceID uint) ([]model.DeviceApplication, error) {
	return s.appRepo.ListByDevice(uint64(deviceID))
}

type AppStats struct {
	TotalApps       int64             `json:"total_apps"`
	IOSApps         int64             `json:"ios_apps"`
	AndroidApps     int64             `json:"android_apps"`
	InstallStats    map[string]int64  `json:"install_stats"`
}

func (s *ApplicationService) GetAppStats() (*AppStats, error) {
	total, _ := s.appRepo.Count()
	iosCount, _ := s.appRepo.CountByPlatform(model.AppPlatformIOS)
	androidCount, _ := s.appRepo.CountByPlatform(model.AppPlatformAndroid)
	installStats, _ := s.appRepo.CountByStatus()

	stats := &AppStats{
		TotalApps:    total,
		IOSApps:      iosCount,
		AndroidApps:  androidCount,
		InstallStats: make(map[string]int64),
	}

	for k, v := range installStats {
		stats.InstallStats[string(k)] = v
	}

	return stats, nil
}
