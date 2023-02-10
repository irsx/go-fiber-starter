package workers

import (
	"fmt"
	"go-fiber-starter/utils"
	"strings"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	Channel   *amqp091.Channel
	Consumers []interface{}
}

type Consumer interface {
	GetName() string
	Declare(channel amqp091.Channel)
	Consume()
}

type WorkerProps struct {
	Name    string
	Channel *amqp091.Channel
	Queue   *amqp091.Queue
}

func (c *Worker) Listen() {
	var listWorkers []string
	for _, item := range c.Consumers {
		consumer := item.(Consumer)
		consumer.Declare(*c.Channel)
		listWorkers = append(listWorkers, consumer.GetName())
	}

	utils.Logger.Info(fmt.Sprintf("✅ LISTEN TO %d WORKERS : %s", len(listWorkers), strings.Join(listWorkers, ", ")))
	var waitGroup sync.WaitGroup
	for i, item := range c.Consumers {
		waitGroup.Add(1)
		go func(i int, item interface{}) {
			defer waitGroup.Done()
			consumer := item.(Consumer)
			consumer.Consume()
		}(i, item)
	}

	waitGroup.Wait()
}
