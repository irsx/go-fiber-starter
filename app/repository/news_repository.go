package repository

import (
	"go-fiber-starter/app/models"

	"gorm.io/gorm"
)

type NewsRepository struct{}

func (r *NewsRepository) GetAll(status int) (news []*models.News, err error) {
	if err := DB.Find(&news, "status = ?", status).Order("updated_at DESC").Error; err != nil {
		return nil, err
	}

	return news, nil
}

func (r *NewsRepository) FindByGUID(guid string) (news models.News, err error) {
	if err := DB.First(&news, "guid = ?", guid).Error; err != nil {
		return news, err
	}

	return news, nil
}

func (r *NewsRepository) Insert(news models.News) (models.News, error) {
	if err := DB.Create(&news).Error; err != nil {
		return news, err
	}

	DB.First(&news, "guid = ?", news.GUID.String())
	return news, nil
}

func (r *NewsRepository) InsertMany(tx *gorm.DB, news []models.News, batchSize int) error {
	if err := DB.CreateInBatches(&news, batchSize).Error; err != nil {
		return err
	}

	return nil
}

func (r *NewsRepository) UpdateByGUID(guid string, storeData models.News) (models.News, error) {
	var news models.News
	if err := DB.Model(&news).Where("guid = ?", guid).Updates(&storeData).Error; err != nil {
		return news, err
	}

	DB.First(&news, "guid = ?", news.GUID.String())
	return news, nil
}

func (r *NewsRepository) DeleteByGUID(guid string) error {
	var news models.News
	if err := DB.Where("guid = ?", guid).Delete(&news).Error; err != nil {
		return err
	}

	return nil
}
