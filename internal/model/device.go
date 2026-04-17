package model

import (
	"time"
)

// Platform 设备平台类型
type Platform string

const (
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
)

// DeviceStatus 设备状态
type DeviceStatus string

const (
	DeviceStatusActive   DeviceStatus = "active"
	DeviceStatusInactive DeviceStatus = "inactive"
	DeviceStatusLost     DeviceStatus = "lost"
	DeviceStatusBlocked  DeviceStatus = "blocked"
)

// Device 设备模型
type Device struct {
	ID          uint64       `json:"id" gorm:"primaryKey"`
	UUID        string       `json:"uuid" gorm:"uniqueIndex;not null"`
	Name        string       `json:"name"`
	Platform    Platform     `json:"platform" gorm:"type:varchar(20)"`
	Status      DeviceStatus `json:"status" gorm:"type:varchar(20);default:'active'"`
	DeviceToken string       `json:"device_token,omitempty"`
	LastSeen    *time.Time   `json:"last_seen"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// DeviceGroup 设备分组
type DeviceGroup struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// DeviceGroupRelation 设备分组关联
type DeviceGroupRelation struct {
	DeviceID uint64 `json:"device_id" gorm:"primaryKey"`
	GroupID  uint64 `json:"group_id" gorm:"primaryKey"`
}

// Application 应用模型
type Application struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	BundleID    string    `json:"bundle_id" gorm:"index"`
	PackageName string    `json:"package_name" gorm:"index"`
	Version     string    `json:"version"`
	FilePath    string    `json:"file_path"`
	FileSize    int64     `json:"file_size"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DeviceApplication 设备应用安装记录
type DeviceApplication struct {
	ID            uint64    `json:"id" gorm:"primaryKey"`
	DeviceID      uint64    `json:"device_id" gorm:"index"`
	ApplicationID uint64    `json:"application_id" gorm:"index"`
	InstalledAt   time.Time `json:"installed_at"`
	Status        string    `json:"status"` // installed, failed, removed
}

// SecurityPolicy 安全策略
type SecurityPolicy struct {
	ID                uint64    `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name" gorm:"not null"`
	Description       string    `json:"description"`
	RequirePassword   bool      `json:"require_password"`
	MinPasswordLength int       `json:"min_password_length"`
	EncryptionEnabled bool      `json:"encryption_enabled"`
	JailbreakDetected bool      `json:"jailbreak_detected"` // 是否检测越狱/root
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// DevicePolicy 设备策略关联
type DevicePolicy struct {
	DeviceID  uint64 `json:"device_id" gorm:"primaryKey"`
	PolicyID  uint64 `json:"policy_id" gorm:"primaryKey"`
	AppliedAt time.Time `json:"applied_at"`
}
