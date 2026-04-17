package repository

import (
	"openmdm/internal/model"

	"gorm.io/gorm"
)

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) Create(app *model.Application) error {
	return r.db.Create(app).Error
}

func (r *ApplicationRepository) GetByID(id uint) (*model.Application, error) {
	var app model.Application
	if err := r.db.First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *ApplicationRepository) List(limit, offset int) ([]model.Application, error) {
	var apps []model.Application
	err := r.db.Order("created_at desc").Limit(limit).Offset(offset).Find(&apps).Error
	return apps, err
}

func (r *ApplicationRepository) ListByPlatform(platform model.AppPlatform, limit, offset int) ([]model.Application, error) {
	var apps []model.Application
	err := r.db.Where("platform = ?", platform).Order("created_at desc").Limit(limit).Offset(offset).Find(&apps).Error
	return apps, err
}

func (r *ApplicationRepository) Update(app *model.Application) error {
	return r.db.Save(app).Error
}

func (r *ApplicationRepository) Delete(id uint) error {
	return r.db.Delete(&model.Application{}, id).Error
}

func (r *ApplicationRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Application{}).Count(&count).Error
	return count, err
}

func (r *ApplicationRepository) CountByPlatform(platform model.AppPlatform) (int64, error) {
	var count int64
	err := r.db.Model(&model.Application{}).Where("platform = ?", platform).Count(&count).Error
	return count, err
}

// ListByDevice 获取设备已安装的应用列表
func (r *ApplicationRepository) ListByDevice(deviceID uint64) ([]model.DeviceApplication, error) {
	var deviceApps []model.DeviceApplication
	err := r.db.Where("device_id = ?", deviceID).Order("created_at desc").Find(&deviceApps).Error
	return deviceApps, err
}

// GetDeviceApp 获取设备应用安装记录
func (r *ApplicationRepository) GetDeviceApp(deviceID uint64, appID uint) (*model.DeviceApplication, error) {
	var deviceApp model.DeviceApplication
	err := r.db.Where("device_id = ? AND application_id = ?", deviceID, appID).First(&deviceApp).Error
	if err != nil {
		return nil, err
	}
	return &deviceApp, nil
}

// CreateDeviceApp 创建设备应用安装记录
func (r *ApplicationRepository) CreateDeviceApp(deviceApp *model.DeviceApplication) error {
	return r.db.Create(deviceApp).Error
}

// UpdateDeviceApp 更新设备应用安装记录
func (r *ApplicationRepository) UpdateDeviceApp(deviceApp *model.DeviceApplication) error {
	return r.db.Save(deviceApp).Error
}

// CountByStatus 统计各状态安装数量
func (r *ApplicationRepository) CountByStatus() (map[model.InstallStatus]int64, error) {
	type Result struct {
		InstallStatus string
		Count         int64
	}
	var results []Result
	err := r.db.Model(&model.DeviceApplication{}).
		Select("install_status, count(*) as count").
		Group("install_status").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	stats := make(map[model.InstallStatus]int64)
	for _, r := range results {
		stats[model.InstallStatus(r.InstallStatus)] = r.Count
	}
	return stats, nil
}
