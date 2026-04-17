package service

import (
	"encoding/json"
	"errors"
	"openmdm/internal/model"
	"openmdm/internal/repository"
)

type PolicyService struct {
	policyRepo *repository.PolicyRepository
	deviceRepo *repository.DeviceRepository
}

func NewPolicyService(policyRepo *repository.PolicyRepository, deviceRepo *repository.DeviceRepository) *PolicyService {
	return &PolicyService{
		policyRepo: policyRepo,
		deviceRepo: deviceRepo,
	}
}

// CreatePolicy 创建策略
func (s *PolicyService) CreatePolicy(policy *model.SecurityPolicy) error {
	if policy.Name == "" {
		return errors.New("policy name is required")
	}
	if policy.Type == "" {
		return errors.New("policy type is required")
	}
	policy.Status = "active"
	return s.policyRepo.Create(policy)
}

// GetPolicy 获取策略详情
func (s *PolicyService) GetPolicy(id uint) (*model.SecurityPolicy, error) {
	return s.policyRepo.GetByID(id)
}

// ListPolicies 列出策略
func (s *PolicyService) ListPolicies(offset, limit int, status string) ([]model.SecurityPolicy, int64, error) {
	return s.policyRepo.List(offset, limit, status)
}

// UpdatePolicy 更新策略
func (s *PolicyService) UpdatePolicy(id uint, updates map[string]interface{}) error {
	// 验证策略存在
	_, err := s.policyRepo.GetByID(id)
	if err != nil {
		return err
	}
	return s.policyRepo.Update(id, updates)
}

// DeletePolicy 删除策略
func (s *PolicyService) DeletePolicy(id uint) error {
	_, err := s.policyRepo.GetByID(id)
	if err != nil {
		return err
	}
	return s.policyRepo.Delete(id)
}

// AssignPolicy 分配策略到设备
func (s *PolicyService) AssignPolicy(deviceID, policyID uint) error {
	// 验证设备存在
	_, err := s.deviceRepo.GetByID(deviceID)
	if err != nil {
		return errors.New("device not found")
	}
	// 验证策略存在
	_, err = s.policyRepo.GetByID(policyID)
	if err != nil {
		return errors.New("policy not found")
	}
	return s.policyRepo.AssignToDevice(deviceID, policyID)
}

// RemovePolicy 从设备移除策略
func (s *PolicyService) RemovePolicy(deviceID, policyID uint) error {
	return s.policyRepo.RemoveFromDevice(deviceID, policyID)
}

// ListDevicePolicies 获取设备关联的策略
func (s *PolicyService) ListDevicePolicies(deviceID uint) ([]model.DevicePolicy, error) {
	return s.policyRepo.ListByDevice(deviceID)
}

// GetPolicyStats 获取策略统计
func (s *PolicyService) GetPolicyStats() (*model.PolicyStats, error) {
	return s.policyRepo.GetStats()
}

// ApplyPolicy 应用策略到设备（模拟）
func (s *PolicyService) ApplyPolicy(deviceID, policyID uint) error {
	// 验证设备存在
	device, err := s.deviceRepo.GetByID(deviceID)
	if err != nil {
		return errors.New("device not found")
	}
	if device.Status != "active" {
		return errors.New("device is not active")
	}

	// 验证策略存在且启用
	policy, err := s.policyRepo.GetByID(policyID)
	if err != nil {
		return errors.New("policy not found")
	}
	if policy.Status != "active" {
		return errors.New("policy is not active")
	}

	// 解析策略规则
	var rules map[string]interface{}
	if err := json.Unmarshal([]byte(policy.Rules), &rules); err != nil {
		rules = make(map[string]interface{})
	}

	// 这里可以添加实际的策略应用逻辑
	// 例如：发送推送通知到设备，让设备执行相应的安全操作
	_ = rules // 使用rules变量

	return nil
}

// ValidateRules 验证策略规则
func (s *PolicyService) ValidateRules(policyType string, rules string) error {
	if rules == "" {
		return errors.New("rules cannot be empty")
	}

	var rulesMap map[string]interface{}
	if err := json.Unmarshal([]byte(rules), &rulesMap); err != nil {
		return errors.New("invalid rules JSON format")
	}

	switch policyType {
	case "password":
		// 验证密码策略规则
		if _, ok := rulesMap["min_length"]; !ok {
			return errors.New("password policy requires min_length")
		}
	case "encryption":
		// 验证加密策略规则
		if _, ok := rulesMap["encryption_enabled"]; !ok {
			return errors.New("encryption policy requires encryption_enabled")
		}
	case "compliance":
		// 合规策略验证
		if _, ok := rulesMap["check_items"]; !ok {
			return errors.New("compliance policy requires check_items")
		}
	}
	return nil
}
