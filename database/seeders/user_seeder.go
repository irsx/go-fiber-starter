package seeders

import (
	"go-fiber-starter/app/models"
	"go-fiber-starter/app/repository"
	"go-fiber-starter/constants"
	"go-fiber-starter/utils"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSeeder struct{}

func (s *UserSeeder) Seed(db *gorm.DB) error {
	utils.Logger.Info("✅ seed data from UserSeeder")
	var usersData []models.User
	password, _ := bcrypt.GenerateFromPassword([]byte(constants.UserDefaultPassword), 14)
	maxSize := 99
	batchSize := 100

	// user admin for demo
	usersData = append(usersData, models.User{
		GUID:        uuid.MustParse(gofakeit.UUID()),
		Name:        constants.UserDefaultName,
		Email:       constants.UserDefaultEmail,
		PhoneNumber: constants.UserDefaultContact,
		Password:    string(password),
	})

	for i := 1; i <= maxSize; i++ {
		data := models.User{
			GUID:        uuid.MustParse(gofakeit.UUID()),
			Name:        gofakeit.Name(),
			Email:       gofakeit.Email(),
			Password:    string(password),
			PhoneNumber: "628" + gofakeit.Phone(),
		}

		usersData = append(usersData, data)
	}

	userRepo := new(repository.UserRepository)
	return userRepo.InsertMany(db, usersData, batchSize)
}
