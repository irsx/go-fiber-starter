package dto

type ImportResultDTO struct {
	Index            int         `json:"index"`
	GUID             string      `json:"guid,omitempty"`
	ErrorCode        string      `json:"error_code,omitempty"`
	ErrorDescription string      `json:"error_description,omitempty"`
	Data             interface{} `json:"data,omitempty"`
}

type ImportResultProductDTO struct {
	SKU  string `json:"sku"`
	Name string `json:"name"`
}
