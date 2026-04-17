package model

import (
	"time"
)

// AppPlatform 应用平台类型
type AppPlatform string

const (
	AppPlatformIOS     AppPlatform = "ios"
	AppPlatformAndroid AppPlatform = "android"
)

// AppStatus 应用状态
type AppStatus string

const (
	AppStatusActive   AppStatus = "active"
	AppStatusInactive AppStatus = "inactive"
)

// InstallStatus 应用安装状态
type InstallStatus string

const (
	InstallStatusPending  InstallStatus = "pending"
	InstallStatusInstalled InstallStatus = "installed"
	InstallStatusFailed    InstallStatus = "failed"
)

// Application 应用模型
type Application struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name"`
	BundleID    string     `json:"bundle_id"` // iOS包名/Android包名
	Platform    AppPlatform `json:"platform"`  // ios/android
	Version     string     `json:"version"`
	FileURL     string     `json:"file_url"`  // 安装包下载地址
	FileSize    int64      `json:"file_size"` // 文件大小字节
	IconURL     string     `json:"icon_url"`  // 应用图标
	Description string     `json:"description"`
	Status      AppStatus  `json:"status"`    // active/inactive
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// DeviceApplication 设备应用安装记录
type DeviceApplication struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	DeviceID      uint          `json:"device_id"`
	ApplicationID uint          `json:"application_id"`
	InstallStatus InstallStatus `json:"install_status"` // pending/installed/failed
	InstalledAt   *time.Time    `json:"installed_at"`
	CreatedAt     time.Time     `json:"created_at"`
}
