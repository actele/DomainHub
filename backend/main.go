package main

import (
	"context"
	"domain-manager/config"
	"domain-manager/controller"
	"domain-manager/ent"
	"domain-manager/middleware"
	"domain-manager/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("config/config.json"); err != nil {
		log.Fatal("加载配置失败:", err)
	}

	// 初始化数据库客户端
	client, err := ent.Open("sqlite3", config.GlobalConfig.Database.Path)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer client.Close()

	// 自动迁移数据库schema
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatal("创建schema失败:", err)
	}

	// 初始化服务
	userService := service.NewUserService(client)
	domainService := service.NewDomainService(client)
	providerService := service.NewProviderService(client)

	if err := userService.EnsureUser(context.Background(), "root", "Tele.Box"); err != nil {
		log.Fatal("初始化默认用户失败:", err)
	}
	if err := providerService.EnsureDefaultProviders(context.Background()); err != nil {
		log.Fatal("初始化默认服务商失败:", err)
	}

	// 初始化控制器
	userController := controller.NewUserController(userService)
	domainController := controller.NewDomainController(domainService)
	providerKeyController := controller.NewProviderKeyController(domainService)
	adminUserController := controller.NewAdminUserController(userService)
	adminProviderController := controller.NewAdminProviderController(providerService)

	// 初始化Gin引擎
	r := gin.Default()

	// 配置CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 初始化路由
	v1 := r.Group("/api/v1")
	{
		// 用户相关路由
		v1.POST("/users/register", userController.Register)
		v1.POST("/users/login", userController.Login)

		// 需要认证的路由
		auth := v1.Group("/", middleware.AuthMiddleware())
		{
			// 用户相关
			auth.POST("/users/change-password", userController.ChangePassword)

			// 域名相关
			auth.GET("/domains", domainController.GetDomains)
			auth.POST("/domains", domainController.AddDomain)
			auth.GET("/domains/:id", domainController.GetDomain)
			auth.DELETE("/domains/:id", domainController.DeleteDomain)
			auth.GET("/domains/available", domainController.GetAvailableDomains)

			// 服务商密钥相关
			auth.GET("/provider-keys", providerKeyController.GetProviderKeys)
			auth.POST("/provider-keys", providerKeyController.AddProviderKey)
			auth.PUT("/provider-keys", providerKeyController.UpdateProviderKey)
			auth.PUT("/provider-keys/:id/name", providerKeyController.UpdateProviderKeyName)
			auth.DELETE("/provider-keys", providerKeyController.DeleteProviderKey)

			// DNS记录相关
			auth.GET("/domains/:id/records", domainController.GetDomainRecords)
			auth.POST("/domains/:id/records", domainController.AddDomainRecord)
			auth.PUT("/domains/:id/records/:record_id", domainController.UpdateDomainRecord)
			auth.DELETE("/domains/:id/records/:record_id", domainController.DeleteDomainRecord)

			// 已启用服务商（所有登录用户可访问）
			auth.GET("/providers", adminProviderController.ListEnabledProviders)

			// 管理员路由
			admin := auth.Group("/admin", middleware.AdminMiddleware())
			{
				// 用户管理
				admin.GET("/users", adminUserController.ListUsers)
				admin.POST("/users", adminUserController.CreateUser)
				admin.PUT("/users/:id/status", adminUserController.UpdateUserStatus)
				admin.PUT("/users/:id/role", adminUserController.UpdateUserRole)
				admin.PUT("/users/:id/password", adminUserController.ResetPassword)

				// 服务商管理
				admin.GET("/providers", adminProviderController.ListProviders)
				admin.POST("/providers", adminProviderController.CreateProvider)
				admin.PUT("/providers/:id/status", adminProviderController.UpdateProviderStatus)
				admin.DELETE("/providers/:id", adminProviderController.DeleteProvider)
			}
		}
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 在goroutine中启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 设置5秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
