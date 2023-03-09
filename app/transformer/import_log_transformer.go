package transformer

import (
	"encoding/json"
	"go-fiber-starter/app/dto"
	"go-fiber-starter/app/models"
	"time"
)

type ImportLogResponse struct {
	GUID         string    `json:"guid"`
	UserGUID     string    `json:"user_guid"`
	ExecTime     float64   `json:"exec_time"`
	TotalSuccess int       `json:"total_success"`
	TotalError   int       `json:"total_error"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ImportLogTransformer(importLogs []*models.ImportLog) (rows []*ImportLogResponse) {
	for _, importLog := range importLogs {
		var row *ImportLogResponse
		jsonResponse, _ := json.Marshal(importLog)
		json.Unmarshal(jsonResponse, &row)
		rows = append(rows, row)
	}

	return rows
}

type ImportLogDetailResponse struct {
	GUID       string                `json:"guid"`
	TotalError int                   `json:"total_error"`
	ErrorLogs  []dto.ImportResultDTO `json:"error_logs"`
}

func ImportLogDetailTransformer(importLog models.ImportLog) (row ImportLogDetailResponse) {
	jsonResponse, _ := json.Marshal(importLog)
	json.Unmarshal(jsonResponse, &row)

	errorLogs := make([]dto.ImportResultDTO, 0)
	if importLog.Errors.Valid {
		json.Unmarshal([]byte(importLog.Errors.String), &errorLogs)
	}

	row.ErrorLogs = errorLogs
	return row
}
