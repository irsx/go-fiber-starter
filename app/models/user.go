package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	GUID        uuid.UUID      `json:"guid" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"type:varchar"`
	Email       string         `json:"email" gorm:"type:varchar;uniqueIndex:unique_idx_user"`
	Password    string         `json:"password" gorm:"type:varchar"`
	PhoneNumber string         `json:"phone_number" gorm:"type:varchar"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (m *User) TableName() string {
	return "user"
}
