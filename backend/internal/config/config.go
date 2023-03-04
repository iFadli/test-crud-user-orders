package config

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
	"test-crud-user-orders/internal/entity"
	"time"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
	}
	Log struct {
		File string
	}
	Service struct {
		Port string
	}
}

func LoadEnv() *Config {
	cfg := &Config{}

	// Database
	cfg.Database.Name = os.Getenv("DB_NAME")
	cfg.Database.Host = os.Getenv("DB_HOST")
	cfg.Database.Port = os.Getenv("DB_PORT")
	cfg.Database.User = os.Getenv("DB_USER")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")

	// Redis
	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Port = os.Getenv("REDIS_PORT")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

	// Log
	cfg.Log.File = os.Getenv("LOG_FILE")

	// Service
	cfg.Service.Port = os.Getenv("SERVICE_PORT")

	return cfg
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&entity.User{}, &entity.OrderItem{}, &entity.OrderHistory{})
}

func SetupDatabase(cfg *Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	maxRetry := 10
	config := cfg.Database
	// Initialize database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Name)
	for i := 0; i < maxRetry; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if i > 0 {
			fmt.Println("DB Connection : Retry Mechanism [" + strconv.Itoa(i) + "x]")
		}
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	return db, err
}

func SetupCache(cfg *Config) *redis.Client {
	config := cfg.Redis
	// Initialize Redis client
	return redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Password: config.Password,
		DB:       0,
	})
}

func SetupFileLog(path string) (*os.File, error) {
	now := time.Now()
	// Format the time in YYYY-MM-DD format
	dateString := now.Format("2006-01-02")
	// Replace "dateformat" with the date in the filename
	newFilename := strings.Replace(path, "dateformat", dateString, 1)

	// Create a log file
	return os.Create(newFilename)
}
