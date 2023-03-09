package services

import (
	"errors"
	"go-fiber-starter/app/dto"
	"go-fiber-starter/app/middlewares"
	"go-fiber-starter/app/models"
	"go-fiber-starter/app/repository"
	"go-fiber-starter/app/transformer"
	"go-fiber-starter/constants"
	"go-fiber-starter/utils"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func (s *AuthService) userRepo() *repository.UserRepository {
	return new(repository.UserRepository)
}

func (s *AuthService) Authenticate(ctx *fiber.Ctx, req *dto.LoginRequestDTO) error {
	if err := req.Validate(); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	user, err := s.userRepo().FindByEmail(req.Email)
	if err != nil {
		err := errors.New("username or password is wrong")
		return utils.JsonErrorUnauthorized(ctx, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return utils.JsonErrorUnauthorized(ctx, constants.ErrInvalidAuth)
	}

	expireHour, _ := time.ParseDuration(os.Getenv("JWT_EXPIRES") + "h")
	expiresAt := time.Now().Add(time.Hour * expireHour).Unix()
	token, err := s.generateToken(user.GUID.String(), expiresAt)
	if err != nil {
		return utils.JsonErrorUnauthorized(ctx, err)
	}

	authTransformer := transformer.AuthLoginTransformer(user, token, expiresAt)
	return utils.JsonSuccess(ctx, authTransformer)
}

func (s *AuthService) Register(ctx *fiber.Ctx, req dto.RegisterRequestDTO) error {
	if err := req.Validate(); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	userRepo := s.userRepo()
	if userRepo.IsExist(req.Email) {
		return utils.JsonErrorValidation(ctx, constants.ErrEmailExist)
	}

	storeUser := s.storeUser(req)
	user, err := userRepo.Insert(repository.DB, storeUser)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_USER_CREATE")
	}

	return utils.JsonSuccess(ctx, fiber.Map{
		"user_guid": user.GUID.String(),
	})
}

func (s *AuthService) generateToken(userGUID string, expiresAt int64) (string, error) {
	claims := middlewares.JwtCustomClaims{
		Issuer: userGUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func (s *AuthService) storeUser(req dto.RegisterRequestDTO) models.User {
	password, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	return models.User{
		Name:        req.FirstName + " " + req.LastName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    string(password),
	}
}
