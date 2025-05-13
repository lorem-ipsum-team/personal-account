package storage

import (
	"github.com/google/uuid"
	"github.com/kerilOvs/profile_sevice/internal/models"
)

type UserStorage interface {
	// Основные операции с пользователем
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	UpdateUser(id uuid.UUID, updates map[string]interface{}) error
	DeleteUser(id uuid.UUID) error

	// Фото пользователя
	AddPhoto(photo *models.UserPhoto) error
	GetUserPhotos(userID uuid.UUID) ([]*models.UserPhoto, error)
	RemovePhoto(userID, photoID uuid.UUID) error
	SetPrimaryPhoto(userID uuid.UUID, photoURL string) error

	// Теги пользователя
	AddTag(tag *models.UserTag) error
	GetUserTags(userID uuid.UUID) ([]*models.UserTag, error)
	RemoveTag(userID, tagID uuid.UUID) error

	// Специальные методы
	UpdateUserAbout(id uuid.UUID, about string) error
	UpdateUserName(id uuid.UUID, name string) error
	UpdateUserSurname(id uuid.UUID, surname string) error
}
