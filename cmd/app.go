package main

import (
	"fmt"
	"go-backend/ent"
	"go-backend/internal/common/ent_client"
	"go-backend/internal/common/env"
	"go-backend/internal/common/middlewares"
	"go-backend/internal/common/rabbitmq"
	"go-backend/internal/common/response"
	"go-backend/internal/common/swagger"
	dependency "go-backend/internal/di"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	ginEngine *gin.Engine
	env       *env.Env
	entClient *ent.Client
	rabbitmq  *rabbitmq.RabbitMQ
}

func NewApp() *App {
	env := env.New()
	// Create a Gin router with default middleware (logger and recovery)
	ginEngine := gin.New()
	ginEngine.Use(gin.Logger())
	ginEngine.Use(middlewares.ErrorHandler())

	ginEngine.Use(gin.CustomRecovery(func(ctx *gin.Context, err any) {
		ctx.Error(response.NewInternalServerErrorException())
		ctx.Abort()
	}))

	// relativePath: bí danh dùng ở web, thay thế cho tên folder
	// root: folder muốn public ra
	// ginEngine.Static("go-backend", "") lộ code
	ginEngine.Static("images", "./public/images")

	// ginEngine.Use(func(ctx *gin.Context){
	// 	if ctx.Request.Method == "OPTIONS"{
	// 		ctx.Header("access-control-allow-headers")
	// 	}
	// })

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
	}
	ginEngine.Use(cors.New(corsConfig))

	// ginEngine.Use(middlewares.A)
	// ginEngine.Use(middlewares.B)
	// ginEngine.Use(middlewares.C)

	swagger.Start(ginEngine, env)

	entClient := ent_client.New(env)
	rabbitmq := rabbitmq.NewRabbitMQ(env)
	dependency.Injection(ginEngine, entClient, env, corsConfig.AllowOrigins, rabbitmq)

	return &App{
		ginEngine: ginEngine,
		env:       env,
		entClient: entClient,
		rabbitmq:  rabbitmq,
	}
}

func (a *App) Start() {
	addr := fmt.Sprintf("%s:%s", a.env.Host, a.env.Port)
	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	a.ginEngine.Run(addr)
}
