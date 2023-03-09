package repository

import (
	"go-fiber-starter/app/models"
)

type ImportLogRepository struct{}

func (r *ImportLogRepository) GetAll() (importLogs []*models.ImportLog, err error) {
	columns := "guid, user_guid, exec_time, total_success, total_error, created_at, updated_at"
	if err := DB.Select(columns).Find(&importLogs).Order("updated_at DESC").Error; err != nil {
		return nil, err
	}

	return importLogs, nil
}

func (r *ImportLogRepository) FindByGUID(guid string) (importLog models.ImportLog, err error) {
	if err := DB.First(&importLog, "guid = ?", guid).Error; err != nil {
		return importLog, err
	}

	return importLog, nil
}

func (r *ImportLogRepository) Insert(importLog models.ImportLog) (models.ImportLog, error) {
	if err := DB.Create(&importLog).Error; err != nil {
		return importLog, err
	}

	return importLog, nil
}
