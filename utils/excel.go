package utils

import (
	"mime/multipart"

	"github.com/xuri/excelize/v2"
)

func ReadExcelFile(file *multipart.FileHeader, sheetName string) (*[][]string, error) {
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
			Logger.Error("❌ excel file close :" + err.Error())
		}
	}()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	return &rows, nil
}
