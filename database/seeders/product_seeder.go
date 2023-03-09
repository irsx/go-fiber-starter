package seeders

import (
	"database/sql"
	"go-fiber-starter/app/models"
	"go-fiber-starter/app/repository"
	"go-fiber-starter/utils"
	"math/rand"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductSeeder struct{}

func (s *ProductSeeder) Seed(db *gorm.DB) error {
	utils.Logger.Info("✅ seed data from ProductSeeder")

	userRepo := new(repository.UserRepository)
	usersGUID := userRepo.GetListGUID(db)
	maxSize := 1000
	batchSize := 100

	var productsData []models.Product
	for i := 1; i <= maxSize; i++ {
		buyPrices := []float64{10000, 20000, 30000, 40000, 50000}
		sellPrices := []float64{15000, 23000, 35000, 42000, 51000}
		expiredTime := gofakeit.DateRange(time.Now(), time.Now().AddDate(5, 0, 0))
		randIndex := gofakeit.RandomInt([]int{0, 1, 2, 3})

		data := models.Product{
			GUID:        uuid.MustParse(gofakeit.UUID()),
			UserGUID:    usersGUID[rand.Intn(len(usersGUID))],
			Name:        gofakeit.Fruit(),
			BarcodeID:   gofakeit.Zip() + strconv.Itoa(i),
			SKU:         "SKU-" + strconv.Itoa(i),
			Description: gofakeit.Lunch(),
			Image:       gofakeit.ImageURL(300, 300),
			Stock:       gofakeit.Number(100, 1000),
			BuyPrice:    buyPrices[randIndex],
			SellPrice:   sellPrices[randIndex],
			ExpiredAt:   sql.NullTime{Time: expiredTime, Valid: gofakeit.Bool()},
		}

		productsData = append(productsData, data)
	}

	productRepo := new(repository.ProductRepository)
	return productRepo.InsertMany(db, productsData, batchSize)
}
