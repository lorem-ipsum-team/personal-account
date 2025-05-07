package api

import (
	"errors"
	"net/http"
	"time"

	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Полная реализация UserServer
func CreateDatabaseIfNotExists(dbUser, dbPassword, dbHost, dbPort, dbName string) error {
	// Подключаемся к системной БД postgres для создания новой БД
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=postgres sslmode=disable",
		dbHost, dbUser, dbPassword, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}

	// Проверяем существование БД
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", dbName).Scan(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	// Создаем БД если не существует
	if count == 0 {
		err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		fmt.Printf("Database '%s' created successfully\n", dbName)
	} else {
		fmt.Printf("Database '%s' already exists\n", dbName)
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()
	return nil
}

// UserModel представляет пользователя в базе данных
type UserModel struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name            string    `gorm:"not null"`
	Surname         string    `gorm:"not null"`
	AboutMyself     *string   `gorm:"type:text"`
	Gender          *string   `gorm:"type:varchar(10)"`
	JungResult      *string   `gorm:"type:text"`
	JungLastAttempt *time.Time
	PrimaryPhoto    *string   `gorm:"type:text"`
	CreatedAt       time.Time `gorm:"not null"`
	DeletedAt       gorm.DeletedAt
}

// UserPhotoModel представляет фото пользователя в базе данных
type UserPhotoModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	URL       string    `gorm:"not null;type:text"`
	CreatedAt time.Time `gorm:"not null"`
}

// UserTagModel представляет тег пользователя в базе данных
type UserTagModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Value     string    `gorm:"not null;type:varchar(100);index"`
	CreatedAt time.Time `gorm:"not null"`
}

type UserServer struct {
	db *gorm.DB
}

func NewUserServer(db *gorm.DB) *UserServer {
	return &UserServer{db: db}
}

// Вспомогательные функции
func errorResponse(msg string) ErrResponse {
	return ErrResponse{Error: &msg}
}

func (s *UserServer) CreateUser(ctx echo.Context, params CreateUserParams) error {
	var req UserCreate
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	// Валидация
	if req.Name == "" || req.Surname == "" {
		return ctx.JSON(http.StatusBadRequest, errorResponse("Name and surname are required"))
	}

	user := UserModel{
		ID:          uuid.New(),
		Name:        req.Name,
		Surname:     req.Surname,
		AboutMyself: req.AboutMyself,
		Gender:      (*string)(req.Gender),
		CreatedAt:   time.Now(),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse("Failed to create user"))
	}

	return ctx.JSON(http.StatusCreated, s.mapUserToResponse(user))
}

func (s *UserServer) DeleteUser(ctx echo.Context, id UserId, params DeleteUserParams) error {
	result := s.db.Where("id = ?", id).Delete(&UserModel{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse("Database error"))
	}
	if result.RowsAffected == 0 {
		return ctx.JSON(http.StatusNotFound, errorResponse("User not found"))
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (s *UserServer) GetUserById(ctx echo.Context, id UserId) error {
	var user UserModel
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, errorResponse("User not found"))
		}
		return ctx.JSON(http.StatusInternalServerError, errorResponse("Database error"))
	}

	return ctx.JSON(http.StatusOK, s.mapUserToResponse(user))
}

func (s *UserServer) UpdateUserAbout(ctx echo.Context, id UserId, params UpdateUserAboutParams) error {
	var req AboutUpdate
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	result := s.db.Model(&UserModel{}).Where("id = ?", id).Update("about_myself", req.AboutMyself)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse("Database error"))
	}
	if result.RowsAffected == 0 {
		return ctx.JSON(http.StatusNotFound, errorResponse("User not found"))
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (s *UserServer) UpdateUserName(ctx echo.Context, id UserId, params UpdateUserNameParams) error {
	var req NameUpdate
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	if req.Name == "" {
		return ctx.JSON(http.StatusBadRequest, errorResponse("Name cannot be empty"))
	}

	result := s.db.Model(&UserModel{}).Where("id = ?", id).Update("name", req.Name)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse("Database error"))
	}
	if result.RowsAffected == 0 {
		return ctx.JSON(http.StatusNotFound, errorResponse("User not found"))
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (s *UserServer) GetUserPhotos(ctx echo.Context, id UserId) error {
	var photos []UserPhotoModel
	if err := s.db.Where("user_id = ?", id).Find(&photos).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse("Database error"))
	}

	response := make([]UserPhoto, len(photos))
	for i, photo := range photos {
		response[i] = UserPhoto{
			Id:  photo.ID,
			Url: photo.URL,
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

func (s *UserServer) AddUserPhoto(ctx echo.Context, id UserId, params AddUserPhotoParams) error {
	var req struct {
		URL string `json:"url"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	if req.URL == "" {
		return ctx.JSON(http.StatusBadRequest, errorResponse("Photo URL cannot be empty"))
	}

	photo := UserPhotoModel{
		ID:     uuid.New(),
		UserID: id,
		URL:    req.URL,
	}

	if err := s.db.Create(&photo).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse("Failed to add photo"))
	}

	return ctx.JSON(http.StatusCreated, UserPhoto{
		Id:  photo.ID,
		Url: photo.URL,
	})
}

func (s *UserServer) RemoveUserPhoto(ctx echo.Context, id UserId, photoId PhotoId, params RemoveUserPhotoParams) error {
	result := s.db.Where("id = ? AND user_id = ?", photoId, id).Delete(&UserPhotoModel{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse("Database error"))
	}
	if result.RowsAffected == 0 {
		return ctx.JSON(http.StatusNotFound, errorResponse("Photo not found"))
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (s *UserServer) UpdatePrimaryPhoto(ctx echo.Context, id UserId, params UpdatePrimaryPhotoParams) error {
	var req PrimaryPhotoUpdate
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	if req.Id == nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("Photo ID is required"))
	}

	// Проверяем, что фото принадлежит пользователю
	var photo UserPhotoModel
	if err := s.db.Where("id = ? AND user_id = ?", req.Id, id).First(&photo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, errorResponse("Photo not found"))
		}
		return ctx.JSON(http.StatusInternalServerError, errorResponse("Database error"))
	}

	// Обновляем primary_photo у пользователя
	if err := s.db.Model(&UserModel{}).Where("id = ?", id).Update("primary_photo", photo.URL).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse("Failed to update primary photo"))
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (s *UserServer) UpdateUserSurname(ctx echo.Context, id UserId, params UpdateUserSurnameParams) error {
	var req SurnameUpdate
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	if req.Surname == "" {
		return ctx.JSON(http.StatusBadRequest, errorResponse("Surname cannot be empty"))
	}

	result := s.db.Model(&UserModel{}).Where("id = ?", id).Update("surname", req.Surname)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse("Database error"))
	}
	if result.RowsAffected == 0 {
		return ctx.JSON(http.StatusNotFound, errorResponse("User not found"))
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (s *UserServer) AddUserTag(ctx echo.Context, id UserId, params AddUserTagParams) error {
	var req TagAdd
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("Invalid request body"))
	}

	if req.Tag == "" {
		return ctx.JSON(http.StatusBadRequest, errorResponse("Tag cannot be empty"))
	}

	tag := UserTagModel{
		ID:     uuid.New(),
		UserID: id,
		Value:  req.Tag,
	}

	if err := s.db.Create(&tag).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse("Failed to add tag"))
	}

	return ctx.JSON(http.StatusCreated, UserTag{
		Id:    tag.ID,
		Value: tag.Value,
	})
}

func (s *UserServer) GetUserTags(ctx echo.Context, id UserId) error {
	var tags []UserTagModel
	if err := s.db.Where("user_id = ?", id).Find(&tags).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse("Database error"))
	}

	response := make([]UserTag, len(tags))
	for i, tag := range tags {
		response[i] = UserTag{
			Id:    tag.ID,
			Value: tag.Value,
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

func (s *UserServer) RemoveUserTag(ctx echo.Context, id UserId, tagId TagId, params RemoveUserTagParams) error {
	result := s.db.Where("id = ? AND user_id = ?", tagId, id).Delete(&UserTagModel{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse("Database error"))
	}
	if result.RowsAffected == 0 {
		return ctx.JSON(http.StatusNotFound, errorResponse("Tag not found"))
	}

	return ctx.NoContent(http.StatusNoContent)
}

// Вспомогательная функция для преобразования модели БД в ответ API
func (s *UserServer) mapUserToResponse(user UserModel) User {
	var photos []UserPhotoModel
	var tags []UserTagModel

	s.db.Where("user_id = ?", user.ID).Find(&photos)
	s.db.Where("user_id = ?", user.ID).Find(&tags)

	photoURLs := make([]string, len(photos))
	for i, photo := range photos {
		photoURLs[i] = photo.URL
	}

	tagValues := make([]string, len(tags))
	for i, tag := range tags {
		tagValues[i] = tag.Value
	}

	var primaryPhoto *string
	if len(photoURLs) > 0 {
		primaryPhoto = &photoURLs[0]
	}

	return User{
		Id:              user.ID,
		Name:            user.Name,
		Surname:         user.Surname,
		AboutMyself:     user.AboutMyself,
		Gender:          (*UserGender)(user.Gender),
		CreatedAt:       user.CreatedAt,
		JungLastAttempt: user.JungLastAttempt,
		JungResult:      user.JungResult,
		Photos:          &photoURLs,
		PrimaryPhoto:    primaryPhoto,
		Tags:            &tagValues,
	}
}
