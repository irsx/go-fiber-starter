package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ImportLog struct {
	GUID         uuid.UUID      `json:"guid" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserGUID     string         `json:"user_guid"`
	ExecTime     float64        `json:"exec_time"`
	TotalSuccess int            `json:"total_success"`
	TotalError   int            `json:"total_error"`
	Errors       sql.NullString `json:"errors"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	User         User           `json:"-"`
}

func (m *ImportLog) TableName() string {
	return "import_log"
}
