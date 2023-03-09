package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wagslane/go-rabbitmq"
)

type DefaultJobPayloads struct {
	Pattern string      `json:"pattern"`
	Data    interface{} `json:"data"`
}

func rabbitMQPublisher() (*rabbitmq.Publisher, error) {
	var (
		host     string = os.Getenv("RABBITMQ_HOST")
		port     string = os.Getenv("RABBITMQ_PORT")
		user     string = os.Getenv("RABBITMQ_USER")
		password string = os.Getenv("RABBITMQ_PASSWORD")
	)
	urlSpec := fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port)
	publisher, err := rabbitmq.NewPublisher(
		urlSpec,
		rabbitmq.Config{},
		rabbitmq.WithPublisherOptionsLogging,
	)

	if err != nil {
		Logger.Error("❌ rabbitmq publisher error : " + err.Error())
		return nil, err
	}

	return publisher, nil
}

func SendJob(key string, payloads []byte) {
	Logger.Info("🔥 RabbitMQ publisher with key " + key)

	publisher, err := rabbitMQPublisher()
	if err != nil {
		Logger.Error("❌ rabbitmq publisher close : " + err.Error())
	}

	defer func() {
		err := publisher.Close()
		if err != nil {
			Logger.Error("❌ rabbitmq publisher close : " + err.Error())
		}
	}()

	returns := publisher.NotifyReturn()
	go func() {
		for r := range returns {
			Logger.Error(fmt.Sprintf("❌ rabbitmq message returned : %s", string(r.Body)))
		}
	}()

	confirmations := publisher.NotifyPublish()
	go func() {
		for c := range confirmations {
			Logger.Info(fmt.Sprintf("✅ rabbitmq message confirmed. tag: %v, ack: %v", c.DeliveryTag, c.Ack))
		}
	}()

	err = publisher.Publish(
		payloads,
		[]string{key},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		// rabbitmq.WithPublishOptionsExchange(""),
	)

	if err != nil {
		Logger.Error("❌ JOB ERROR : " + err.Error())
	}
}

func SendJobWithDefaultPayloads(key string, data interface{}) {
	payloads := DefaultJobPayloads{
		Pattern: key,
		Data:    data,
	}

	payloadsBytes, _ := json.Marshal(payloads)
	SendJob(key, payloadsBytes)
}
