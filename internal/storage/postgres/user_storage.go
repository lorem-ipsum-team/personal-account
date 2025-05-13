package postgres

import (
	"errors"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/kerilOvs/profile_sevice/internal/models"
)

type UserPostgresStorage struct {
	db *gorm.DB
}

func NewUserPostgresStorage(db *gorm.DB) *UserPostgresStorage {
	return &UserPostgresStorage{db: db}
}

func (s *UserPostgresStorage) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}

func (s *UserPostgresStorage) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserPostgresStorage) UpdateUser(id uuid.UUID, updates map[string]interface{}) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

func (s *UserPostgresStorage) DeleteUser(id uuid.UUID) error {
	return s.db.Where("id = ?", id).Delete(&models.User{}).Error
}

func (s *UserPostgresStorage) AddPhoto(photo *models.UserPhoto) error {
	return s.db.Create(photo).Error
}

func (s *UserPostgresStorage) GetUserPhotos(userID uuid.UUID) ([]*models.UserPhoto, error) {
	var photos []*models.UserPhoto
	err := s.db.Where("user_id = ?", userID).Find(&photos).Error
	return photos, err
}

func (s *UserPostgresStorage) RemovePhoto(userID, photoID uuid.UUID) error {
	return s.db.Where("id = ? AND user_id = ?", photoID, userID).Delete(&models.UserPhoto{}).Error
}

func (s *UserPostgresStorage) SetPrimaryPhoto(userID uuid.UUID, photoURL string) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Update("primary_photo", photoURL).Error
}

func (s *UserPostgresStorage) AddTag(tag *models.UserTag) error {
	return s.db.Create(tag).Error
}

func (s *UserPostgresStorage) GetUserTags(userID uuid.UUID) ([]*models.UserTag, error) {
	var tags []*models.UserTag
	err := s.db.Where("user_id = ?", userID).Find(&tags).Error
	return tags, err
}

func (s *UserPostgresStorage) RemoveTag(userID, tagID uuid.UUID) error {
	return s.db.Where("id = ? AND user_id = ?", tagID, userID).Delete(&models.UserTag{}).Error
}

func (s *UserPostgresStorage) UpdateUserAbout(id uuid.UUID, about string) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("about_myself", about).Error
}

func (s *UserPostgresStorage) UpdateUserName(id uuid.UUID, name string) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("name", name).Error
}

func (s *UserPostgresStorage) UpdateUserSurname(id uuid.UUID, surname string) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("surname", surname).Error
}
