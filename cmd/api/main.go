package main

import (
	"fmt"
	"log"
	"openmdm/internal/config"
	"openmdm/internal/handler"
	"openmdm/internal/middleware"
	"openmdm/internal/model"
	"openmdm/internal/repository"
	"openmdm/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Load config
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to SQLite database
	db, err := gorm.Open(sqlite.Open("openmdm.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Auto migrate tables
	if err := db.AutoMigrate(
		&model.Device{},
		&model.DeviceGroup{},
		&model.DeviceGroupRelation{},
		&model.Application{},
		&model.DeviceApplication{},
		&model.SecurityPolicy{},
		&model.DevicePolicy{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize dependencies
	deviceRepo := repository.NewDeviceRepository(db)
	deviceSvc := service.NewDeviceService(deviceRepo)
	deviceHandler := handler.NewDeviceHandler(deviceSvc)

	appRepo := repository.NewApplicationRepository(db)
	appSvc := service.NewApplicationService(appRepo, deviceRepo)
	appHandler := handler.NewApplicationHandler(appSvc)

	policyRepo := repository.NewPolicyRepository(db)
	policySvc := service.NewPolicyService(policyRepo, deviceRepo)
	policyHandler := handler.NewPolicyHandler(policySvc)

	// Initialize Gin
	if cfg.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"version": cfg.App.Version,
		})
	})

	// API routes
	v1 := r.Group("/api/v1")
	{
		devices := v1.Group("/devices")
		{
			devices.POST("", deviceHandler.Create)
			devices.GET("", deviceHandler.List)
			devices.GET("/:id", deviceHandler.Get)
			devices.PUT("/:id", deviceHandler.Update)
			devices.DELETE("/:id", deviceHandler.Delete)
			devices.POST("/:id/lock", deviceHandler.Lock)
			devices.POST("/:id/wipe", deviceHandler.Wipe)
			devices.POST("/:id/policies", policyHandler.Assign)
			devices.DELETE("/:id/policies/:policyId", policyHandler.Remove)
			devices.GET("/:id/policies", policyHandler.ListByDevice)
			devices.POST("/:id/policies/:policyId/apply", policyHandler.Apply)
		}

		apps := v1.Group("/apps")
		{
			apps.POST("", appHandler.Create)
			apps.GET("", appHandler.List)
			apps.GET("/:id", appHandler.Get)
			apps.PUT("/:id", appHandler.Update)
			apps.DELETE("/:id", appHandler.Delete)
		}

		deviceApps := v1.Group("/devices/:id/apps")
		{
			deviceApps.POST("", appHandler.Install)
			deviceApps.GET("", appHandler.ListByDevice)
		}

		policies := v1.Group("/policies")
		{
			policies.POST("", policyHandler.Create)
			policies.GET("", policyHandler.List)
			policies.GET("/stats", policyHandler.Stats)
			policies.GET("/:id", policyHandler.Get)
			policies.PUT("/:id", policyHandler.Update)
			policies.DELETE("/:id", policyHandler.Delete)
		}
	}

	// Start server
	addr := cfg.Server.Address()
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	fmt.Println("OpenMDM server started successfully!")
}
