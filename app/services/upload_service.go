package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-fiber-starter/utils"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UploadService struct{}

type CDNResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Imageurl string `json:"imageurl"`
}

func (s *UploadService) SendCallback(ctx *fiber.Ctx, payloads interface{}) error {
	type PipedreamResponse struct {
		Success bool `json:"success"`
	}

	callbackUrl := os.Getenv("CALLBACK_URL")
	if callbackUrl == "" {
		return utils.JsonErrorEnvironment(ctx, "CALLBACK_URL")
	}

	utils.Logger.Info("🔥 h2h url : " + callbackUrl + "/callback")
	body, _ := json.Marshal(payloads)
	RestClient := utils.RestClient{
		Url:             callbackUrl + "/callback",
		Method:          fiber.MethodPost,
		Headers:         nil,
		Payload:         body,
		ResponseSuccess: &PipedreamResponse{},
		ResponseError:   &PipedreamResponse{},
		ErrorPrefix:     "CALLBACK",
		LogRequest:      true,
	}

	clientResp, err := RestClient.Send(ctx)
	if err != nil {
		return err
	}

	return utils.JsonSuccess(ctx, clientResp)
}

func (s *UploadService) UploadToCDN(ctx *fiber.Ctx, imagePath string) (resp CDNResponse, err error) {
	cdnUrl := os.Getenv("CDN_URL")
	if cdnUrl == "" {
		return resp, utils.JsonErrorEnvironment(ctx, "CDN_URL")
	}

	// custom headers
	body, contentType, err := s.customCDNContentType("file", imagePath)
	if err != nil {
		return resp, utils.JsonErrorInternal(ctx, err, "E_CDN_CONTENT_TYPE")
	}

	var headers []utils.RestClientHeaders
	headers = append(headers, utils.RestClientHeaders{
		Key:   "Content-Type",
		Value: contentType,
	})

	utils.Logger.Info("🔥 h2h url : " + cdnUrl)
	RestClient := utils.RestClient{
		Url:             cdnUrl,
		Method:          fiber.MethodPost,
		Headers:         headers,
		Payload:         body,
		ResponseSuccess: &CDNResponse{},
		ResponseError:   &CDNResponse{},
		ErrorPrefix:     "CDN",
		LogRequest:      false,
	}

	clientResp, err := RestClient.Send(ctx)
	os.Remove(imagePath) // remove uploaded image before return
	if err != nil {
		return resp, utils.JsonErrorInternal(ctx, err, "E_CDN_CLIENT")
	}

	return clientResp.(CDNResponse), nil
}

func (s *UploadService) customCDNContentType(key string, path string) (bodyByte []byte, contentType string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return bodyByte, contentType, err
	}
	defer file.Close()

	// custom content type with boundary
	var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
	key = quoteEscaper.Replace(key)
	path = quoteEscaper.Replace(filepath.Base(path))
	contentDisposition := fmt.Sprintf(`form-data; name="%s"; filename="%s"`, key, path)

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", contentDisposition)
	header.Set("Content-Type", "image/png")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreatePart(header)
	if err != nil {
		return bodyByte, contentType, err
	}

	_, err = io.Copy(fw, file)
	if err != nil {
		return bodyByte, contentType, err
	}

	writer.Close()

	return body.Bytes(), writer.FormDataContentType(), nil
}
