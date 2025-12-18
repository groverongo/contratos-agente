package common

import (
	"fmt"
	"os"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", d.Host, d.Port, d.User, d.Password, d.Name)
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

type AuthConfig struct {
	SecretKey string
	ProjectId string
}

type ServerConfig struct {
	Database    DatabaseConfig
	Minio       MinioConfig
	AuthService AuthConfig
}

func LoadConfig() *ServerConfig {
	config := &ServerConfig{
		Database: DatabaseConfig{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			Name:     os.Getenv("DATABASE_NAME"),
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
		},
		Minio: MinioConfig{
			Endpoint:  os.Getenv("MINIO_ENDPOINT"),
			AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
			SecretKey: os.Getenv("MINIO_SECRET_KEY"),
		},
		AuthService: AuthConfig{
			SecretKey: os.Getenv("STACK_SECRET_SERVER_KEY"),
			ProjectId: os.Getenv("STACK_PROJECT_ID"),
		},
	}

	return config
}
