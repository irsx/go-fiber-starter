package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type HttpClient struct {
	Url             string
	Method          string
	Headers         []HttpClientHeaders
	Payload         []byte
	Args            *fiber.Args
	ResponseSuccess interface{}
	ResponseError   interface{}
	ErrorPrefix     string
	LogRequest      bool
	LogResponse     bool
}

type HttpClientHeaders struct {
	Key   string
	Value string
}

const USER_AGENT = "irsx"
const JSON_CONTENT_TYPE = "application/json"

func (hc *HttpClient) Send(ctx *fiber.Ctx) (interface{}, error) {
	// setup user agent
	client := fiber.AcquireAgent()
	client.UserAgent(USER_AGENT)

	// setup request
	clientReq := client.Request()
	clientReq.SetRequestURI(hc.Url)
	clientReq.Header.SetMethod(hc.Method)

	// setup custom headers
	if hc.Headers != nil {
		for _, header := range hc.Headers {
			clientReq.Header.Set(header.Key, header.Value)
		}
	} else {
		clientReq.Header.SetContentType(JSON_CONTENT_TYPE)
	}

	if hc.Payload != nil {
		clientReq.SetBody(hc.Payload)
	}

	if hc.Args != nil {
		client.Form(hc.Args)
		fiber.ReleaseArgs(hc.Args)
	}

	// log request
	if hc.LogRequest {
		Logger.Info(clientReq.String())
	}

	if err := client.Parse(); err != nil {
		return nil, err
	}

	// get response raw
	respCode, respBody, respErrs := client.Bytes()
	if respErrs != nil {
		err := errors.New(extractResponseErrors(respErrs))
		return nil, err
	}

	// check response error
	var respJson interface{}
	if respCode >= 300 {
		respJson = hc.ResponseError
		json.Unmarshal(respBody, &respJson)
		Logger.Info(fmt.Sprintf("❌ HTTP ERROR [%d] %s", respCode, string(respBody)))
		err := errors.New(string(respBody))
		return nil, err
	}

	// format response
	respJson = hc.ResponseSuccess
	json.Unmarshal(respBody, &respJson)

	// log response
	if hc.LogResponse {
		Logger.Info(fmt.Sprintf("response : [%d] %s", respCode, respBody))
	}

	return respJson, nil
}

func extractResponseErrors(errors []error) string {
	var message string
	for i, err := range errors {
		if i > 0 {
			message += " | "
		}

		message += err.Error()
	}

	return message
}

func ContentTypeFormHeader() HttpClientHeaders {
	return HttpClientHeaders{
		Key:   "Content-Type",
		Value: "application/x-www-form-urlencoded",
	}
}

func AuthorizationHeader(token string) HttpClientHeaders {
	return HttpClientHeaders{
		Key:   "Authorization",
		Value: "Bearer " + token,
	}
}

func JsonContentTypeHeader() HttpClientHeaders {
	return HttpClientHeaders{
		Key:   "Content-Type",
		Value: JSON_CONTENT_TYPE,
	}
}
