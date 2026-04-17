package repository

import (
	"openmdm/internal/model"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Create(device *model.Device) error {
	return r.db.Create(device).Error
}

func (r *DeviceRepository) GetByID(id uint64) (*model.Device, error) {
	var device model.Device
	if err := r.db.First(&device, id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) GetByUUID(uuid string) (*model.Device, error) {
	var device model.Device
	if err := r.db.Where("uuid = ?", uuid).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) List(limit, offset int) ([]model.Device, error) {
	var devices []model.Device
	err := r.db.Order("created_at desc").Limit(limit).Offset(offset).Find(&devices).Error
	return devices, err
}

func (r *DeviceRepository) Update(device *model.Device) error {
	return r.db.Save(device).Error
}

func (r *DeviceRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Device{}, id).Error
}

func (r *DeviceRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Device{}).Count(&count).Error
	return count, err
}
