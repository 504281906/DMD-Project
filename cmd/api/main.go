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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动迁移表
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

	// 初始化依赖
	deviceRepo := repository.NewDeviceRepository(db)
	deviceSvc := service.NewDeviceService(deviceRepo)
	deviceHandler := handler.NewDeviceHandler(deviceSvc)

	// 应用分发模块
	appRepo := repository.NewApplicationRepository(db)
	appSvc := service.NewApplicationService(appRepo, deviceRepo)
	appHandler := handler.NewApplicationHandler(appSvc)

	// 初始化 Gin
	if cfg.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"version": cfg.App.Version,
		})
	})

	// API 路由
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
		}

		// 应用管理
		apps := v1.Group("/apps")
		{
			apps.POST("", appHandler.Create)
			apps.GET("", appHandler.List)
			apps.GET("/:id", appHandler.Get)
			apps.PUT("/:id", appHandler.Update)
			apps.DELETE("/:id", appHandler.Delete)
		}

		// 设备应用安装
		deviceApps := v1.Group("/devices/:id/apps")
		{
			deviceApps.POST("", appHandler.Install)
			deviceApps.GET("", appHandler.ListByDevice)
		}
	}

	// 启动服务器
	addr := cfg.Server.Address()
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	fmt.Println("OpenMDM server started successfully!")
}
