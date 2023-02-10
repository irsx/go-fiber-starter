package jobs

import (
	"fmt"
	"go-fiber-starter/configs"
	"go-fiber-starter/utils"

	"github.com/wagslane/go-rabbitmq"
)

type DefaultJobPayloads struct {
	Pattern string      `json:"pattern"`
	Data    interface{} `json:"data"`
}

func SendJob(routingKey string, payloads []byte) {
	utils.Logger.Info("🔥 RabbitMQ publisher with key " + routingKey)

	publisher, err := configs.RabbitMQPublisher()
	if err != nil {
		utils.Logger.Error("❌ rabbitmq publisher : " + err.Error())
	}

	defer func() {
		err := publisher.Close()
		if err != nil {
			utils.Logger.Error("❌ rabbitmq publisher close : " + err.Error())
		}
	}()

	returns := publisher.NotifyReturn()
	go func() {
		for r := range returns {
			utils.Logger.Error(fmt.Sprintf("❌ rabbitmq message returned : %s", string(r.Body)))
		}
	}()

	confirmations := publisher.NotifyPublish()
	go func() {
		for c := range confirmations {
			utils.Logger.Info(fmt.Sprintf("✅ rabbitmq message confirmed. tag: %v, ack: %v", c.DeliveryTag, c.Ack))
		}
	}()

	err = publisher.Publish(
		payloads,
		[]string{routingKey},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		// rabbitmq.WithPublishOptionsExchange(""),
	)

	if err != nil {
		utils.Logger.Error("❌ JOB ERROR : " + err.Error())
	}
}
