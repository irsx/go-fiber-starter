package models

import (
	"time"

	"github.com/google/uuid"
)

type ImportLog struct {
	GUID         uuid.UUID `json:"guid" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserGUID     string    `json:"user_guid" gorm:"type:uuid"`
	Source       string    `json:"source" gorm:"type:varchar"`
	ExecTime     float64   `json:"exec_time"`
	TotalSuccess int       `json:"total_success"`
	TotalError   int       `json:"total_error"`
	Errors       string    `json:"errors" gorm:"default:null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	User         User      `json:"-" gorm:"foreignKey:UserGUID"`
}

func (m *ImportLog) TableName() string {
	return "import_log"
}
