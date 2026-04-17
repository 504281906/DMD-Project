package handler

import (
	"net/http"
	"openmdm/internal/model"
	"openmdm/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct {
	svc *service.ApplicationService
}

func NewApplicationHandler(svc *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{svc: svc}
}

// Create 创建应用
// POST /api/v1/apps
func (h *ApplicationHandler) Create(c *gin.Context) {
	var req service.CreateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app, err := h.svc.CreateApp(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, app)
}

// List 列出应用
// GET /api/v1/apps
func (h *ApplicationHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	platform := c.Query("platform")

	filter := &service.ListAppsFilter{
		Limit:  limit,
		Offset: offset,
	}

	if platform != "" {
		filter.Platform = model.AppPlatform(platform)
	}

	apps, total, err := h.svc.ListApps(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   apps,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// Get 获取应用详情
// GET /api/v1/apps/:id
func (h *ApplicationHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	app, err := h.svc.GetApp(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	c.JSON(http.StatusOK, app)
}

// Update 更新应用
// PUT /api/v1/apps/:id
func (h *ApplicationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req service.UpdateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app, err := h.svc.UpdateApp(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, app)
}

// Delete 删除应用
// DELETE /api/v1/apps/:id
func (h *ApplicationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.DeleteApp(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Install 设备安装应用
// POST /api/v1/devices/:id/apps
func (h *ApplicationHandler) Install(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	var req service.InstallAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deviceApp, err := h.svc.InstallApp(uint(deviceID), req.ApplicationID)
	if err != nil {
		if err == service.ErrDeviceNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}
		if err == service.ErrAppNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, deviceApp)
}

// ListByDevice 列出设备已安装应用
// GET /api/v1/devices/:id/apps
func (h *ApplicationHandler) ListByDevice(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	deviceApps, err := h.svc.ListDeviceApps(uint(deviceID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": deviceApps,
	})
}
