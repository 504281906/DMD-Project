package model

// PolicyStats 策略统计
type PolicyStats struct {
	TotalPolicies   int64 `json:"total_policies"`
	ActivePolicies  int64 `json:"active_policies"`
	TotalDevices    int64 `json:"total_devices"`
	AppliedPolicies int64 `json:"applied_policies"`
}
