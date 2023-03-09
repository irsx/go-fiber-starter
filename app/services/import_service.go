package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go-fiber-starter/app/dto"
	"go-fiber-starter/app/models"
	"go-fiber-starter/app/repository"
	"go-fiber-starter/app/transformer"
	"go-fiber-starter/constants"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/sse"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ImportService struct{}

func (s *ImportService) ProductRepo() *repository.ProductRepository {
	return new(repository.ProductRepository)
}

func (s *ImportService) importLog() *repository.ImportLogRepository {
	return new(repository.ImportLogRepository)
}

func (s *ImportService) GetLogs(ctx *fiber.Ctx) error {
	importLogs, err := s.importLog().GetAll()
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_IMPORT_LOG_ALL")
	}

	data := transformer.ImportLogTransformer(importLogs)
	return utils.JsonSuccess(ctx, data)
}

func (s *ImportService) GetLogDetail(ctx *fiber.Ctx, guid string) error {
	importLog, err := s.importLog().FindByGUID(guid)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_IMPORT_LOG_DETAIL")
	}

	row := transformer.ImportLogDetailTransformer(importLog)
	return utils.JsonSuccess(ctx, row)
}

func (s *ImportService) ImportProductExcel(ctx *fiber.Ctx, userGUID string) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_IMPORT_FILE")
	}

	rows, err := utils.ReadExcelFile(file, "products")
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_IMPORT_READ")
	}

	go s.importProductExcelRows(*rows, userGUID)

	return utils.JsonSuccess(ctx, fiber.Map{
		"total": len(*rows) - 1,
	})
}

func (s *ImportService) storeExcelProduct(row []string, index int, userGUID string) dto.ImportResultDTO {
	var (
		err       error
		errorCode string
		product   models.Product
		indexStr  string = "#" + strconv.Itoa(index)
	)

	productRepo := s.ProductRepo()
	err = repository.DB.Transaction(func(tx *gorm.DB) error {
		if len(row) == 0 {
			errorCode = "E_IMPORT_PRODUCT_EMPTY_ROW"
			return errors.New("empty row at index " + indexStr)
		}

		if row[3] == "" {
			errorCode = "E_IMPORT_PRODUCT_SKU"
			return errors.New("empty sku " + indexStr)
		}

		storeData := s.storeProductFromExcel(row)
		storeData.UserGUID = userGUID
		product, err = productRepo.Insert(tx, storeData)
		if err != nil {
			errorCode = "E_IMPORT_PRODUCT_ADD"
			return err
		}

		return nil
	})

	if err != nil {
		utils.Logger.Error(fmt.Sprintf("--> ❌ %s [%s] : %s", indexStr, errorCode, err.Error()))
		return dto.ImportResultDTO{
			Index:            index,
			ErrorCode:        errorCode,
			ErrorDescription: err.Error(),
			Data: dto.ImportResultProductDTO{
				Name: row[1],
				SKU:  row[3],
			},
		}
	}

	// SEND QUEUE
	utils.SendJobWithDefaultPayloads(constants.QueueNewProduct, product)

	// utils.Logger.Info("--> ✅ " + indexStr + " import success")
	return dto.ImportResultDTO{
		Index: index,
		GUID:  product.GUID.String(),
	}
}

func (s *ImportService) importProductExcelRows(rows [][]string, userGUID string) {
	start := time.Now()
	totalWorker, _ := strconv.Atoi(os.Getenv("APP_MAX_WORKER"))

	// init worker pool
	wp := workerpool.New(totalWorker)
	utils.Logger.Info(fmt.Sprintf("🔥 START WORKER POOL WITH %d CPU", totalWorker))

	var totalSuccess int = 0
	var resultErrors []dto.ImportResultDTO
	for index, row := range rows {
		if index == 0 {
			continue
		}

		row := row
		index := index
		wp.Submit(func() {
			result := s.storeExcelProduct(row, index, userGUID)
			if result.ErrorCode != "" {
				resultErrors = append(resultErrors, result)
			} else {
				totalSuccess++
			}

			response, _ := json.Marshal(result)
			sse.NotifierController.Append(string(response))
		})
	}

	// stops the worker pool
	wp.StopWait()

	execTime := math.Round(time.Since(start).Seconds()*100) / 100
	s.storeImportLog(userGUID, execTime, totalSuccess, resultErrors)

	utils.Logger.Info(fmt.Sprintf("✅ import excel time: %f seconds", execTime))
}

func (s *ImportService) storeProductFromExcel(row []string) (Product models.Product) {
	Product.Name = row[0]
	Product.SKU = row[1]
	Product.BarcodeID = row[2]
	Product.Description = row[3]

	if row[4] != "" {
		imageList := strings.Split(row[4], ",")
		imageJson, _ := json.Marshal(imageList)
		Product.Image = string(imageJson)
	}

	stock, _ := strconv.Atoi(row[5])
	Product.Stock = stock

	sellPrice, _ := strconv.ParseFloat(row[6], 64)
	Product.SellPrice = sellPrice

	buyPrice, _ := strconv.ParseFloat(row[7], 64)
	Product.BuyPrice = buyPrice

	if len(row) == 9 {
		parseExpiredAt, _ := time.Parse(constants.TimestampFormat, row[8])
		Product.ExpiredAt = sql.NullTime{Time: parseExpiredAt, Valid: true}
	} else {
		Product.ExpiredAt = sql.NullTime{Valid: false}
	}

	Product.CreatedAt = time.Now()
	return Product
}

func (s *ImportService) storeImportLog(userGUID string, execTime float64, totalSuccess int, resultErrors []dto.ImportResultDTO) {
	resultErrorBytes, _ := json.Marshal(resultErrors)
	storeImportLog := models.ImportLog{
		UserGUID:     userGUID,
		ExecTime:     execTime,
		TotalSuccess: totalSuccess,
		TotalError:   len(resultErrors),
		Errors:       sql.NullString{String: string(resultErrorBytes), Valid: true},
	}

	_, err := s.importLog().Insert(storeImportLog)
	if err != nil {
		utils.Logger.Info("E_IMPORT_LOG : " + err.Error())
	}
}
