package utils

import (
	"fmt"
	"go-fiber-starter/constants"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func UploadFileToStorage(header *multipart.FileHeader) (string, error) {
	file, err := header.Open()
	if err != nil {
		Logger.Error("error Retrieving the File : " + err.Error())
		return "", err
	}
	defer file.Close()

	Logger.Info(fmt.Sprintf("Uploaded File : %+v\n", header.Filename))
	Logger.Info(fmt.Sprintf("File Size : %+v kb\n", float32(header.Size)/(1<<10)))
	Logger.Info(fmt.Sprintf("MIME Header : %+v\n", header.Header))

	fileExt := filepath.Ext(header.Filename)
	randString := RandStringRunes(12)
	fileName := TimestampString() + "_" + randString + fileExt
	imagePath := constants.UploadDir + "/" + fileName
	dst, err := os.Create(imagePath)
	if err != nil {
		Logger.Error("error creating file : " + err.Error())
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		Logger.Error("error copy file : " + err.Error())
	}

	Logger.Info("Successfully uploaded file to " + imagePath)
	return imagePath, nil
}
