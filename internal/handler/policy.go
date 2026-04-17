package handler

import (
	"net/http"
	"openmdm/internal/model"
	"openmdm/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PolicyHandler struct {
	policySvc *service.PolicyService
}

func NewPolicyHandler(policySvc *service.PolicyService) *PolicyHandler {
	return &PolicyHandler{policySvc: policySvc}
}

// Create 创建策略
func (h *PolicyHandler) Create(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Description string `json:"description"`
		Rules       string `json:"rules"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证规则格式
	if err := h.policySvc.ValidateRules(input.Type, input.Rules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	policy := &model.SecurityPolicy{
		Name:        input.Name,
		Type:        input.Type,
		Description: input.Description,
		Rules:       input.Rules,
		Status:      "active",
	}

	if err := h.policySvc.CreatePolicy(policy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "policy created",
		"data":    policy,
	})
}

// List 列出策略
func (h *PolicyHandler) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	policies, total, err := h.policySvc.ListPolicies(offset, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  policies,
		"total": total,
	})
}

// Get 获取策略详情
func (h *PolicyHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}

	policy, err := h.policySvc.GetPolicy(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": policy})
}

// Update 更新策略
func (h *PolicyHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证规则格式（如果更新了rules）
	if rules, ok := input["rules"].(string); ok {
		policy, _ := h.policySvc.GetPolicy(id)
		if policy != nil {
			policyType := input["type"]
			if policyType == "" {
				policyType = policy.Type
			}
			if err := h.policySvc.ValidateRules(policyType.(string), rules); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
	}

	if err := h.policySvc.UpdatePolicy(id, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "policy updated"})
}

// Delete 删除策略
func (h *PolicyHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}

	if err := h.policySvc.DeletePolicy(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "policy deleted"})
}

// Assign 分配策略到设备
func (h *PolicyHandler) Assign(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	var input struct {
		PolicyID uint64 `json:"policy_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.policySvc.AssignPolicy(deviceID, input.PolicyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "policy assigned to device"})
}

// Remove 从设备移除策略
func (h *PolicyHandler) Remove(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	policyID, err := strconv.ParseUint(c.Query("policy_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}

	if err := h.policySvc.RemovePolicy(deviceID, policyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "policy removed from device"})
}

// ListByDevice 获取设备关联的策略
func (h *PolicyHandler) ListByDevice(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	policies, err := h.policySvc.ListDevicePolicies(deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": policies})
}

// Apply 应用策略到设备
func (h *PolicyHandler) Apply(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	policyID, err := strconv.ParseUint(c.Param("policyId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}

	if err := h.policySvc.ApplyPolicy(deviceID, policyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "policy applied to device"})
}

// Stats 获取策略统计
func (h *PolicyHandler) Stats(c *gin.Context) {
	stats, err := h.policySvc.GetPolicyStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
