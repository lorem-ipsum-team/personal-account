package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"

	"github.com/kerilOvs/profile_sevice/api"
	config "github.com/kerilOvs/profile_sevice/internal/config"
	"github.com/labstack/echo/v4/middleware"
)

// This is the main function of the program
func main() {
	var f api.UserGender = api.UserGender(api.UserCreateGenderMALE)

	fmt.Println(f)

	config, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}
	api.CreateDatabaseIfNotExists(config.Database.User, config.Database.Password, config.Database.Host, strconv.Itoa(config.Database.Port), config.Database.Dbname)
	fmt.Println("db user :", config.Database.User)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Database.Host,
		config.Database.User,
		config.Database.Password,
		config.Database.Dbname,
		strconv.Itoa(config.Database.Port),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 2. Автомиграция (создание таблиц)
	if err := db.AutoMigrate(
		&api.UserModel{},
		&api.UserPhotoModel{},
		&api.UserTagModel{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 3. Инициализация сервера
	server := api.NewUserServer(db)

	// 4. Запуск Echo-сервера
	e := echo.New() // а не, библиотека нужна
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"}, // Разрешить фронтенд
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	api.RegisterHandlers(e, server) //  реализовано

	log.Printf("Server started on :%s", strconv.Itoa(config.Server.Port))
	log.Fatal(e.Start(":" + strconv.Itoa(config.Server.Port)))
}
