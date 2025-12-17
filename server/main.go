package main

import (
	"log"
	"os"

	"github.com/contract-ai/server/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=postgres user=postgres password=postgres dbname=contract_ai port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// MinIO connection
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "minio:9000"
	}
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := false

	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: useSSL,
	})
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
	// Auth middleware (mock/stub for StackAuth - assuming User ID comes in header or token)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// In a real app, verify StackAuth token here
			// For now, checks Authorization header exists
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				// Allow login/public routes if any (none here really)
				// return echo.ErrUnauthorized
			}
			return next(c)
		}
	})

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
