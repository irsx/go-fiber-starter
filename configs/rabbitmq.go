package configs

import (
	"fmt"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func RabbitMQConnection() (connection *amqp091.Connection, channel *amqp091.Channel, err error) {
	var (
		host     string = os.Getenv("RABBITMQ_HOST")
		port     string = os.Getenv("RABBITMQ_PORT")
		user     string = os.Getenv("RABBITMQ_USER")
		password string = os.Getenv("RABBITMQ_PASSWORD")
	)
	url := fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port)
	connection, err = amqp091.Dial(url)
	if err != nil {
		fmt.Println(err)
		fmt.Println("ERROR : Failed to connect to RabbitMQ")
		return nil, nil, err
	}

	channel, err = connection.Channel()
	if err != nil {
		defer connection.Close()
		fmt.Println("ERROR : RabbitMQ Failed to open a channel")
		return nil, nil, err
	}
	return connection, channel, err
}
