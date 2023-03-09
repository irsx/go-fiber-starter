package seeders

import (
	"go-fiber-starter/app/models"
	"go-fiber-starter/app/repository"
	"go-fiber-starter/utils"
	"math/rand"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NewsSeeder struct{}

func (s *NewsSeeder) Seed(db *gorm.DB) error {
	utils.Logger.Info("✅ seed data from NewsSeeder")

	userRepo := new(repository.UserRepository)
	usersGUID := userRepo.GetListGUID(db)
	maxSize := 200
	batchSize := 50

	var newssData []models.News
	for i := 1; i <= maxSize; i++ {
		data := models.News{
			GUID:        uuid.MustParse(gofakeit.UUID()),
			UserGUID:    usersGUID[rand.Intn(len(usersGUID))],
			Title:       gofakeit.Sentence(5),
			Description: gofakeit.Paragraph(1, 20, 50, " "),
			Image:       gofakeit.ImageURL(400, 400),
			HyperLink:   gofakeit.URL(),
			Status:      strconv.Itoa(gofakeit.RandomInt([]int{0, 1})),
		}

		newssData = append(newssData, data)
	}

	newsRepo := new(repository.NewsRepository)
	return newsRepo.InsertMany(db, newssData, batchSize)
}
