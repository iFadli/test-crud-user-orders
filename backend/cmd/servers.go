package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"test-crud-user-orders/internal/handler"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"test-crud-user-orders/internal/config"
	"test-crud-user-orders/internal/repository"
	"test-crud-user-orders/internal/usecase"
)

type Server struct {
	e *echo.Echo
}

func NewServer() *Server {
	// Load environment variables from .env file
	loadConfig := config.LoadEnv()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}

	// Setup Cache (redis) & Database (MariaDB)
	cache := config.SetupCache(loadConfig)
	db, errDB := config.SetupDatabase(loadConfig)
	if errDB != nil {
		log.Fatalf("error connecting to database: %s", errDB.Error())
	}
	// Do AutoMigrate of Database
	if errMigrate := config.AutoMigrate(db); errMigrate != nil {
		log.Fatalf("error initializing table: %s", errMigrate.Error())
	}

	// init Repository, UseCase, and Handler of User table
	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase)

	// init Repository, UseCase, and Handler of Order Item table
	orderItemRepo := repository.NewOrderItemRepository(db)
	orderItemUseCase := usecase.NewOrderItemUseCase(orderItemRepo, cache)
	orderItemHandler := handler.NewOrderItemHandler(orderItemUseCase)

	// init Repository, UseCase, and Handler of Order History table
	orderHistoryRepo := repository.NewOrderHistoryRepository(db)
	orderHistoryUseCase := usecase.NewOrderHistoryUseCase(orderHistoryRepo, orderItemRepo, userRepo)
	orderHistoryHandler := handler.NewOrderHistoryHandler(orderHistoryUseCase)

	// init Path of User Table
	pathUser := e.Group("/users")
	pathUser.POST("/", userHandler.Create)
	pathUser.GET("/", userHandler.GetAllPagination)
	pathUser.GET("/:id", userHandler.GetByID)
	pathUser.PUT("/:id", userHandler.Update)
	pathUser.DELETE("/:id", userHandler.Delete)

	pathUser.GET("/:id/order-histories", orderHistoryHandler.GetHistoryByUserID)

	// init Path of OrderItem Table
	pathOrderItems := e.Group("/order-items")
	pathOrderItems.POST("/", orderItemHandler.Create)
	pathOrderItems.GET("/", orderItemHandler.GetAllPagination)
	pathOrderItems.GET("/:id", orderItemHandler.GetByID)
	pathOrderItems.PUT("/:id", orderItemHandler.Update)
	pathOrderItems.DELETE("/:id", orderItemHandler.Delete)

	// init Path of OrderHistory Table
	pathOrderHistory := e.Group("/order-histories")
	pathOrderHistory.POST("/", orderHistoryHandler.Create)
	pathOrderHistory.GET("/", orderHistoryHandler.GetAllPagination)
	pathOrderHistory.GET("/:id", orderHistoryHandler.GetByID)
	pathOrderHistory.PUT("/:id", orderHistoryHandler.Update)
	pathOrderHistory.DELETE("/:id", orderHistoryHandler.Delete)

	return &Server{e}
}

func (s *Server) Start() {
	loadConfig := config.LoadEnv()

	// init Logger File
	fileLog, errLog := config.SetupFileLog(loadConfig.Log.File)
	if errLog != nil {
		log.Fatal(errLog)
	}
	defer func(fileLog *os.File) {
		err := fileLog.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(fileLog)
	// Log all Pre-Request and Post-Request to the Logger File
	s.e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: fileLog,
	}))

	addr := fmt.Sprintf(":%s", loadConfig.Service.Port)

	go func() {
		if err := s.e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.e.Shutdown(ctx); err != nil {
		log.Fatalf("server error: %s\n", err.Error())
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
