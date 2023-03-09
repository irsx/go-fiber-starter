package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	GUID        uuid.UUID      `json:"guid" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	Password    string         `json:"password"`
	PhoneNumber string         `json:"phone_number"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (m *User) TableName() string {
	return "user"
}
