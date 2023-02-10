package workers

import (
	"go-fiber-starter/app/services"
	"go-fiber-starter/constants"
	"go-fiber-starter/utils"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type UserRegisterWorker struct {
	Name    string
	Channel *amqp091.Channel
	Queue   *amqp091.Queue
}

func (c *UserRegisterWorker) GetName() string {
	return constants.QueueUserRegister
}

func (c *UserRegisterWorker) Declare(channel amqp091.Channel) {
	c.Name = c.GetName()
	c.Channel = &channel

	queue, err := channel.QueueDeclare(c.Name, true, false, false, false, nil)
	utils.FailOnError(err, "Failed to declare a queue `"+c.Name+"`")

	c.Queue = &queue
}

func (c *UserRegisterWorker) Consume() {
	msgs, err := c.Channel.Consume(c.Queue.Name, c.Name, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("error consuming messages: %s", err)
	}

	for msg := range msgs {
		utils.Logger.Info("🔥 CONSUME " + c.Queue.Name + " : " + string(msg.Body))
		userService := new(services.UserService)
		if err := userService.UpdateUserFromConsumer(msg.Body); err != nil {
			utils.Logger.Error("🆘 ========= ERROR =========")
			utils.Logger.Error("CONSUMER : " + c.Queue.Name)
			utils.Logger.Error(err.Error())
			utils.Logger.Error("🆘 ======= END ERROR =======")
		}
		msg.Ack(true)
	}
}
