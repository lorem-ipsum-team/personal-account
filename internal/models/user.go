package models

import (
	"time"

	"github.com/google/uuid"
)

type UserGender string

const (
	GenderMale   UserGender = "MALE"
	GenderFemale UserGender = "FEMALE"
)

type User struct {
	ID              uuid.UUID   `json:"id"`
	Name            string      `json:"name"`
	Surname         string      `json:"surname"`
	AboutMyself     *string     `json:"about_myself,omitempty"`
	Gender          *UserGender `json:"gender,omitempty"`
	BirthDate       *time.Time  `json:"birth_date,omitempty"`
	JungResult      *string     `json:"jung_result,omitempty"`
	JungLastAttempt *time.Time  `json:"jung_last_attempt,omitempty"`
	PrimaryPhoto    *string     `json:"primary_photo,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	Photos          []UserPhoto `gorm:"foreignKey:UserID" json:"photos,omitempty"`
	Tags            []UserTag   `gorm:"foreignKey:UserID" json:"tags,omitempty"`
}

type UserPhoto struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	URL    string    `json:"url"`
}

type UserTag struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Value  string    `json:"value"`
}

type UserProfileUpdate struct {
	Name            *string     `json:"name,omitempty"`
	Surname         *string     `json:"surname,omitempty"`
	AboutMyself     *string     `json:"about_myself,omitempty"`
	Gender          *UserGender `json:"gender,omitempty"`
	BirthDate       *time.Time  `json:"birth_date,omitempty"`
	JungResult      *string     `json:"jung_result,omitempty"`
	JungLastAttempt *time.Time  `json:"jung_last_attempt,omitempty"` // Новое поле
}
