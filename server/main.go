package main

import (
	"log"

	"github.com/contract-ai/server/common"
	"github.com/contract-ai/server/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	config := common.LoadConfig()

	// Database connection
	db, err := common.ConnectDB(config)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// MinIO connection
	minioClient, err := common.ConnectMinio(config)
	if err != nil {
		log.Fatalf("failed to connect minio: %v", err)
	}

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize Handlers
	h := handlers.New(db, minioClient)

	// Routes
	// Auth middleware (Capture Bearer Token, and decode JWT)
	e.Use(common.StackAuthValidation)

	api := e.Group("/api")

	// Contracts
	api.POST("/contracts", h.CreateContract)
	api.GET("/contracts", h.ListContracts)
	api.GET("/contracts/:id/file", h.GetContractFileUrl)
	api.POST("/contracts/:id/recipients", h.UpdateRecipients)
	api.PUT("/contracts/:id/sign", h.SignContract) // User signs

	// Chat
	api.POST("/chat/ask", h.AskAI)

	e.Logger.Fatal(e.Start(":8080"))
}
