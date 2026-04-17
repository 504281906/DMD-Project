package service

import (
	"errors"
	"openmdm/internal/model"
	"openmdm/internal/repository"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDeviceNotFound = errors.New("device not found")
	ErrDeviceExists   = errors.New("device already exists")
)

type DeviceService struct {
	repo *repository.DeviceRepository
}

func NewDeviceService(repo *repository.DeviceRepository) *DeviceService {
	return &DeviceService{repo: repo}
}

type CreateDeviceRequest struct {
	Name     string        `json:"name"`
	Platform model.Platform `json:"platform"`
}

func (s *DeviceService) Create(req *CreateDeviceRequest) (*model.Device, error) {
	// 生成唯一UUID
	deviceUUID := uuid.New().String()
	now := time.Now()
	
	device := &model.Device{
		UUID:     deviceUUID,
		Name:     req.Name,
		Platform: req.Platform,
		Status:   model.DeviceStatusActive,
		LastSeen: &now,
	}
	
	if err := s.repo.Create(device); err != nil {
		return nil, err
	}
	
	return device, nil
}

func (s *DeviceService) GetByID(id uint64) (*model.Device, error) {
	device, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrDeviceNotFound
	}
	return device, nil
}

func (s *DeviceService) GetByUUID(uuid string) (*model.Device, error) {
	device, err := s.repo.GetByUUID(uuid)
	if err != nil {
		return nil, ErrDeviceNotFound
	}
	return device, nil
}

type ListDevicesRequest struct {
	Limit  int
	Offset int
}

func (s *DeviceService) List(req *ListDevicesRequest) ([]model.Device, int64, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}
	
	devices, err := s.repo.List(req.Limit, req.Offset)
	if err != nil {
		return nil, 0, err
	}
	
	count, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	
	return devices, count, nil
}

type UpdateDeviceRequest struct {
	Name     string
	Status   model.DeviceStatus
}

func (s *DeviceService) Update(id uint64, req *UpdateDeviceRequest) (*model.Device, error) {
	device, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrDeviceNotFound
	}
	
	if req.Name != "" {
		device.Name = req.Name
	}
	if req.Status != "" {
		device.Status = req.Status
	}
	
	if err := s.repo.Update(device); err != nil {
		return nil, err
	}
	
	return device, nil
}

func (s *DeviceService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

// 远程锁定设备
func (s *DeviceService) Lock(id uint64) error {
	device, err := s.repo.GetByID(id)
	if err != nil {
		return ErrDeviceNotFound
	}
	
	device.Status = model.DeviceStatusBlocked
	return s.repo.Update(device)
}

// 远程擦除设备
func (s *DeviceService) Wipe(id uint64) error {
	device, err := s.repo.GetByID(id)
	if err != nil {
		return ErrDeviceNotFound
	}
	
	// 更新状态，实际擦除操作由客户端执行
	device.Status = model.DeviceStatusLost
	return s.repo.Update(device)
}