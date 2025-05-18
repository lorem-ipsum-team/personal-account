package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	//"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/kerilOvs/profile_sevice/internal/config"
	"github.com/kerilOvs/profile_sevice/internal/handlers"
	"github.com/kerilOvs/profile_sevice/internal/models"
	"github.com/kerilOvs/profile_sevice/internal/service"
	"github.com/kerilOvs/profile_sevice/internal/storage/minio"
	"github.com/kerilOvs/profile_sevice/internal/storage/rabbit"

	//"github.com/kerilOvs/profile_sevice/logger"
	postgresstorage "github.com/kerilOvs/profile_sevice/internal/storage/postgres"
	"github.com/kerilOvs/profile_sevice/pkg/logger"
)

func main() {

	//logLevel := os.Getenv("LOG_LEVEL")
	//logFormat := os.Getenv("LOG_FORMAT")
	//log := logger.Init(logFormat, logLevel)
	// 1. Загрузка конфигурации
	log := logger.Init("text", "debug")
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Error("Failed to read config:", slog.Any("error", err))
	}
	log.Info("Read config", slog.Any("config", cfg))

	log.Info("Connecting to db...")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Dbname,
		strconv.Itoa(cfg.Database.Port),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect to database:", slog.Any("error", err))
	}

	log.Info("Running migrations...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.UserPhoto{},
		&models.UserTag{},
	); err != nil {
		log.Error("Failed to migrate database:", slog.Any("error", err))
	}

	// 4. Инициализация MinIO клиента
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("Connection to minio")
	minioClient, err := minio.New(ctx, cfg.Minio)
	if err != nil {
		log.Error("Failed to initialize MinIO client:", slog.Any("error", err))
	}

	log.Info("Connecting to Rabbit")
	rabbitRepo, err := rabbit.New(ctx, &cfg.Rabbit)
	if err != nil {
		log.Error("failed to create rabbit_repo", slog.Any("error", err))

		return
	}

	// 5. Инициализация слоев приложения
	userStorage := postgresstorage.NewUserPostgresStorage(db)
	userService := service.NewUserService(userStorage, rabbitRepo)

	// Инициализация фото сервиса
	photoService := service.NewPhotoService(minioClient.Client, cfg.Minio.Bucket)
	photoHandler := handlers.NewPhotoHandler(userService, photoService)

	// 6. Настройка Echo сервера
	e := echo.New()
	e.Use(handlers.Logging(log))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	// 7. Регистрация маршрутов
	userHandler := handlers.NewUserHandler(userService)
	registerRoutes(e, userHandler, photoHandler)

	// 8. Запуск сервера
	serverAddr := ":" + strconv.Itoa(cfg.Server.Port)
	log.Info("Server started", slog.String("port", serverAddr))
	log.Error("Server stopped", slog.Any("error", e.Start(serverAddr)))
}

func registerRoutes(e *echo.Echo, userHandler *handlers.UserHandler, photoHandler *handlers.PhotoHandler) {
	e.POST("/users", userHandler.CreateUser)                     // +
	e.DELETE("/users/:id", userHandler.DeleteUser)               // ?
	e.GET("/users/:id", userHandler.GetUserById)                 // +
	e.PATCH("/users/:id/profile", userHandler.UpdateUserProfile) // +
	e.PATCH("/users/:id/about", userHandler.UpdateUserAbout)     // depricated
	e.PATCH("/users/:id/name", userHandler.UpdateUserName)       // depricated
	e.PATCH("/users/:id/surname", userHandler.UpdateUserSurname) // depricated

	e.GET("/users/:id/photos", userHandler.GetUserPhotos) // +
	//e.PUT("/users/:id/photos", userHandler.AddUserPhoto)
	e.DELETE("/users/:id/photos/:photoId", userHandler.RemoveUserPhoto) // +
	e.PATCH("/users/:id/primary_photo", userHandler.UpdatePrimaryPhoto) // + (айди фото)

	e.PUT("/users/:id/tag", userHandler.AddUserTag)               // +
	e.GET("/users/:id/tags", userHandler.GetUserTags)             // +
	e.DELETE("/users/:id/tags/:tagId", userHandler.RemoveUserTag) // +

	// Фото маршруты
	e.POST("/users/:id/addphoto", photoHandler.UploadPhoto) // + по айди юзера добавляет фотку
	e.GET("/photos/:id", photoHandler.GetPhoto)             // + по айди !фото! отдает фотку
}
