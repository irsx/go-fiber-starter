package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-fiber-starter/app/dto"
	"go-fiber-starter/app/jobs"
	"go-fiber-starter/app/models"
	"go-fiber-starter/app/repository"
	"go-fiber-starter/app/transformer"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/sse"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
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
	rows, err := s.readUploadedExcelFile(ctx, "product_list")
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_IMPORT_READ")
	}

	go s.importProductExcelRows(*rows, userGUID)

	return utils.JsonSuccess(ctx, fiber.Map{
		"total": len(*rows) - 1,
	})
}

func (s *ImportService) readUploadedExcelFile(ctx *fiber.Ctx, sheetName string) (*[][]string, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, err
	}

	opened, err := file.Open()
	if err != nil {
		return nil, err
	}

	f, err := excelize.OpenReader(opened)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			utils.Logger.Error("❌ err file close :" + err.Error())
		}
	}()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	return &rows, nil
}

func (s *ImportService) storeExcelProduct(row []string, index int) dto.ImportResultDTO {
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
			return errors.New("empty row at index #" + indexStr)
		}

		storeData := s.storeProductFromExcel(row)
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
	jobs.ProductStreamJob(product)

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
			result := s.storeExcelProduct(row, index)
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
	s.storeImportLog("Product", userGUID, execTime, totalSuccess, resultErrors)

	utils.Logger.Info(fmt.Sprintf("✅ import excel time: %f seconds", execTime))
}

func (s *ImportService) storeProductFromExcel(row []string) (Product models.Product) {
	Product.CreatedAt = time.Now()

	Product.Name = row[0]
	Product.SKU = row[1]
	Product.Description = row[2]

	if row[3] != "" {
		imageList := strings.Split(row[3], ",")
		imageJson, _ := json.Marshal(imageList)
		Product.Image = string(imageJson)
	}

	stock, _ := strconv.Atoi(row[4])
	Product.Stock = stock

	sellPrice, _ := strconv.ParseFloat(row[5], 64)
	Product.SellPrice = sellPrice

	buyPrice, _ := strconv.ParseFloat(row[6], 64)
	Product.BuyPrice = buyPrice

	if len(row) == 7 {
		Product.ExpiredAt = &row[7]
	} else {
		Product.ExpiredAt = nil
	}

	return Product
}

func (s *ImportService) storeImportLog(source string, userGUID string, execTime float64, totalSuccess int, resultErrors []dto.ImportResultDTO) {
	resultErrorBytes, _ := json.Marshal(resultErrors)
	storeImportLog := models.ImportLog{
		UserGUID:     userGUID,
		Source:       source,
		ExecTime:     execTime,
		TotalSuccess: totalSuccess,
		TotalError:   len(resultErrors),
		Errors:       string(resultErrorBytes),
	}

	_, err := s.importLog().Insert(storeImportLog)
	if err != nil {
		utils.Logger.Info("E_IMPORT_LOG : " + err.Error())
	}
}
