package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kerilOvs/profile_sevice/internal/models"
	"github.com/kerilOvs/profile_sevice/internal/storage"
)

type UserService struct {
	storage storage.UserStorage
}

func NewUserService(storage storage.UserStorage) *UserService {
	return &UserService{storage: storage}
}

func (s *UserService) CreateUser(id uuid.UUID, name, surname string, aboutMyself *string, gender *models.UserGender) (*models.User, error) {
	if name == "" || surname == "" {
		return nil, errors.New("name and surname are required")
	}

	user := &models.User{
		ID:          id,
		Name:        name,
		Surname:     surname,
		AboutMyself: aboutMyself,
		Gender:      gender,
		CreatedAt:   time.Now(),
	}

	if err := s.storage.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := s.storage.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Загружаем связанные данные
	photos, err := s.storage.GetUserPhotos(id)
	if err != nil {
		return nil, err
	}

	tags, err := s.storage.GetUserTags(id)
	if err != nil {
		return nil, err
	}

	// Преобразуем []*UserPhoto → []UserPhoto
	photoVals := make([]models.UserPhoto, len(photos))
	for i, p := range photos {
		if p != nil {
			photoVals[i] = *p
		}
	}

	// Преобразуем []*UserTag → []UserTag
	tagVals := make([]models.UserTag, len(tags))
	for i, t := range tags {
		if t != nil {
			tagVals[i] = *t
		}
	}

	user.Photos = photoVals
	user.Tags = tagVals

	return user, nil
}

func (s *UserService) UpdateUserProfile(id uuid.UUID, updates models.UserProfileUpdate) error {
	updateFields := make(map[string]interface{})

	if updates.Name != nil {
		if *updates.Name == "" {
			return errors.New("name cannot be empty")
		}
		updateFields["name"] = *updates.Name
	}

	if updates.Surname != nil {
		if *updates.Surname == "" {
			return errors.New("surname cannot be empty")
		}
		updateFields["surname"] = *updates.Surname
	}

	if updates.AboutMyself != nil {
		updateFields["about_myself"] = *updates.AboutMyself
	}

	if updates.Gender != nil {
		updateFields["gender"] = *updates.Gender
	}

	if updates.BirthDate != nil {
		updateFields["birth_date"] = *updates.BirthDate
	}

	if updates.JungResult != nil {
		if !isValidJungType(*updates.JungResult) {
			return errors.New("invalid Jung personality type")
		}
		updateFields["jung_result"] = *updates.JungResult
		now := time.Now()
		updateFields["jung_last_attempt"] = now
	}

	if len(updateFields) > 0 {
		return s.storage.UpdateUser(id, updateFields)
	}

	return nil
}

func (s *UserService) AddUserPhoto(userID uuid.UUID, photoURL string) (*models.UserPhoto, error) {
	if photoURL == "" {
		return nil, errors.New("photo URL cannot be empty")
	}

	photo := &models.UserPhoto{
		ID:     uuid.New(),
		UserID: userID,
		URL:    photoURL,
	}

	if err := s.storage.AddPhoto(photo); err != nil {
		return nil, err
	}

	return photo, nil
}

func (s *UserService) SetPrimaryPhoto(userID, photoID uuid.UUID) error {
	photos, err := s.storage.GetUserPhotos(userID)
	if err != nil {
		return err
	}

	var photoURL string
	for _, photo := range photos {
		if photo.ID == photoID {
			photoURL = photo.URL
			break
		}
	}

	if photoURL == "" {
		return errors.New("photo not found")
	}

	return s.storage.SetPrimaryPhoto(userID, photoURL)
}

func (s *UserService) AddUserTag(userID uuid.UUID, tagValue string) (*models.UserTag, error) {
	if tagValue == "" {
		return nil, errors.New("tag cannot be empty")
	}

	tag := &models.UserTag{
		ID:     uuid.New(),
		UserID: userID,
		Value:  tagValue,
	}

	if err := s.storage.AddTag(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

func isValidJungType(jungType string) bool {
	validTypes := map[string]bool{
		"INTJ": true, "INTP": true, "ENTJ": true, "ENTP": true,
		"INFJ": true, "INFP": true, "ENFJ": true, "ENFP": true,
		"ISTJ": true, "ISFJ": true, "ESTJ": true, "ESFJ": true,
		"ISTP": true, "ISFP": true, "ESTP": true, "ESFP": true,
	}
	return validTypes[jungType]
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.storage.DeleteUser(id)
}

func (s *UserService) UpdateUserAbout(id uuid.UUID, about string) error {
	return s.storage.UpdateUserAbout(id, about)
}

func (s *UserService) UpdateUserName(id uuid.UUID, name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	return s.storage.UpdateUserName(id, name)
}

func (s *UserService) UpdateUserSurname(id uuid.UUID, surname string) error {
	if surname == "" {
		return errors.New("surname cannot be empty")
	}
	return s.storage.UpdateUserSurname(id, surname)
}

func (s *UserService) GetUserPhotos(userID uuid.UUID) ([]*models.UserPhoto, error) {
	return s.storage.GetUserPhotos(userID)
}

func (s *UserService) RemoveUserPhoto(userID, photoID uuid.UUID) error {
	return s.storage.RemovePhoto(userID, photoID)
}

func (s *UserService) GetUserTags(userID uuid.UUID) ([]*models.UserTag, error) {
	return s.storage.GetUserTags(userID)
}

func (s *UserService) RemoveUserTag(userID, tagID uuid.UUID) error {
	return s.storage.RemoveTag(userID, tagID)
}
