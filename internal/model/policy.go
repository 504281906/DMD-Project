package model

import (
	"time"
)

// SecurityPolicy 安全策略
type SecurityPolicy struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`                  // 策略名称
	Type        string    `json:"type"`                  // password/encryption/compliance/wipe
	Description string    `json:"description"`           // 策略描述
	Rules       string    `json:"rules"`                 // JSON格式的策略规则
	Status      string    `json:"status"`                // active/inactive
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DevicePolicy 设备与策略的关联
type DevicePolicy struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	DeviceID       uint      `json:"device_id"`
	PolicyID       uint      `json:"policy_id"`
	AssignStatus   string    `json:"assign_status"`      // pending/applied/failed
	AppliedAt      *time.Time `json:"applied_at"`        // 策略应用时间
	CreatedAt      time.Time `json:"created_at"`
}

// PolicyStats 策略统计
type PolicyStats struct {
	TotalPolicies   int64 `json:"total_policies"`
	ActivePolicies  int64 `json:"active_policies"`
	TotalDevices    int64 `json:"total_devices"`
	AppliedPolicies int64 `json:"applied_policies"`
}
