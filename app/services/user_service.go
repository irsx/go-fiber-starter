package services

import (
	"encoding/json"
	"go-fiber-starter/app/dto"
	"go-fiber-starter/app/models"
	"go-fiber-starter/app/repository"
	"go-fiber-starter/constants"
	"go-fiber-starter/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func (s *UserService) userRepo() *repository.UserRepository {
	return new(repository.UserRepository)
}

func (s *UserService) List(ctx *fiber.Ctx, paginate utils.Pagination) error {
	paginationData, err := s.userRepo().GetAll(paginate)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_USER_LIST")
	}

	return utils.JsonPagination(ctx, paginationData)
}

func (s *UserService) Add(ctx *fiber.Ctx, req dto.UserRequestDTO) (err error) {
	if err := req.Validate(); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	userRepo := s.userRepo()
	if userRepo.IsExist(req.Email) {
		return utils.JsonErrorValidation(ctx, constants.ErrEmailExist)
	}

	storeData := s.StoreFromRequest(models.User{}, req)
	user, err := userRepo.Insert(repository.DB, storeData)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_USER_ADD")
	}

	return utils.JsonSuccess(ctx, fiber.Map{"guid": user.GUID.String()})
}

func (s *UserService) Update(ctx *fiber.Ctx, guid string, req dto.UserRequestDTO) (err error) {
	if err := req.Validate(); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	user, err := s.userRepo().FindByGUID(guid)
	if err != nil {
		return utils.JsonErrorNotFound(ctx, err)
	}

	storeData := s.StoreFromRequest(user, req)
	_, err = s.userRepo().UpdateByGUID(repository.DB, guid, storeData)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_USER_UPDATE")
	}

	return utils.JsonSuccess(ctx, fiber.Map{"guid": guid})
}

func (s *UserService) Delete(ctx *fiber.Ctx, guid string) error {
	if err := s.userRepo().DeleteByGUID(guid); err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_USER_DELETE")
	}

	return utils.JsonSuccess(ctx, fiber.Map{"guid": guid})
}

func (s *UserService) StoreFromRequest(user models.User, req dto.UserRequestDTO) models.User {
	pw := []byte(req.Password)
	var pwgenerated []byte
	if req.Password != "" {
		pwgenerated, _ = bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	}

	user = models.User{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    string(pwgenerated),
	}
	return user
}

func (s *UserService) UpdateUserFromConsumer(data []byte) error {
	var UserRegisterConsumerDTO dto.UserRegisterConsumerDTO

	payload := &utils.DefaultJobPayloads{}
	json.Unmarshal(data, &payload)

	if err := mapstructure.Decode(payload.Data, &UserRegisterConsumerDTO); err != nil {
		return err
	}

	utils.Logger.Info("✅ Update data user by email : " + UserRegisterConsumerDTO.Email)

	user := models.User{
		Email: UserRegisterConsumerDTO.Email,
		Name:  UserRegisterConsumerDTO.Name,
	}

	if err := s.userRepo().UpdateByEmail(UserRegisterConsumerDTO.Email, user); err != nil {
		return err
	}

	return nil
}
