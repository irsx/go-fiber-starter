package models

import (
	"time"

	"github.com/google/uuid"
)

type News struct {
	GUID        uuid.UUID `json:"guid" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserGUID    string    `json:"user_guid" gorm:"type:uuid"`
	Title       string    `json:"title" gorm:"type:varchar"`
	Description string    `json:"description"`
	Image       string    `json:"image" gorm:"type:varchar"`
	HyperLink   string    `json:"hyperlink" gorm:"type:varchar"`
	Status      string    `json:"status" gorm:"type:smallint;default:1"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `json:"-"`
}

func (m *News) TableName() string {
	return "news"
}
