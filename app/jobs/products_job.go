package jobs

import (
	"encoding/json"
	"go-fiber-starter/constants"
)

func ProductStreamJob(data interface{}) {
	payloads := DefaultJobPayloads{
		Pattern: constants.QueueNewProduct,
		Data:    data,
	}

	payloadsBytes, _ := json.Marshal(payloads)
	SendJob(constants.QueueNewProduct, payloadsBytes)
}
