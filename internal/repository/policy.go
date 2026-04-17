package repository

import (
	"openmdm/internal/model"

	"gorm.io/gorm"
)

type PolicyRepository struct {
	db *gorm.DB
}

func NewPolicyRepository(db *gorm.DB) *PolicyRepository {
	return &PolicyRepository{db: db}
}

// Create 创建策略
func (r *PolicyRepository) Create(policy *model.SecurityPolicy) error {
	return r.db.Create(policy).Error
}

// GetByID 根据ID获取策略
func (r *PolicyRepository) GetByID(id uint) (*model.SecurityPolicy, error) {
	var policy model.SecurityPolicy
	if err := r.db.First(&policy, id).Error; err != nil {
		return nil, err
	}
	return &policy, nil
}

// List 列出所有策略
func (r *PolicyRepository) List(offset, limit int, status string) ([]model.SecurityPolicy, int64, error) {
	var policies []model.SecurityPolicy
	var total int64

	query := r.db.Model(&model.SecurityPolicy{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&policies).Error; err != nil {
		return nil, 0, err
	}

	return policies, total, nil
}

// Update 更新策略
func (r *PolicyRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.SecurityPolicy{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除策略
func (r *PolicyRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先删除关联的设备策略
		if err := tx.Where("policy_id = ?", id).Delete(&model.DevicePolicy{}).Error; err != nil {
			return err
		}
		// 删除策略
		return tx.Delete(&model.SecurityPolicy{}, id).Error
	})
}

// AssignToDevice 分配策略到设备
func (r *PolicyRepository) AssignToDevice(deviceID, policyID uint) error {
	dp := model.DevicePolicy{
		DeviceID:     deviceID,
		PolicyID:     policyID,
		AssignStatus: "pending",
	}
	return r.db.Create(&dp).Error
}

// RemoveFromDevice 从设备移除策略
func (r *PolicyRepository) RemoveFromDevice(deviceID, policyID uint) error {
	return r.db.Where("device_id = ? AND policy_id = ?", deviceID, policyID).Delete(&model.DevicePolicy{}).Error
}

// ListByDevice 获取设备关联的策略
func (r *PolicyRepository) ListByDevice(deviceID uint) ([]model.DevicePolicy, error) {
	var dps []model.DevicePolicy
	if err := r.db.Where("device_id = ?", deviceID).Find(&dps).Error; err != nil {
		return nil, err
	}
	return dps, nil
}

// GetStats 获取策略统计
func (r *PolicyRepository) GetStats() (*model.PolicyStats, error) {
	var stats model.PolicyStats

	r.db.Model(&model.SecurityPolicy{}).Count(&stats.TotalPolicies)
	r.db.Model(&model.SecurityPolicy{}).Where("status = ?", "active").Count(&stats.ActivePolicies)
	r.db.Model(&model.Device{}).Count(&stats.TotalDevices)
	r.db.Model(&model.DevicePolicy{}).Where("assign_status = ?", "applied").Count(&stats.AppliedPolicies)

	return &stats, nil
}
