package sse

import (
	"bufio"
	"go-fiber-starter/utils"

	fiberUtils "github.com/gofiber/fiber/v2/utils"
)

type Broker interface {
	Subscribe(writer *bufio.Writer) string
	Unsubscribe(id string)
	CountSubscriber() int64
	GetSubscriber() map[string]*bufio.Writer
	Listen()
}

type broker struct {
	subscriber map[string]*bufio.Writer
	notifier   Notifier
}

func NewBroker(notifier Notifier) Broker {
	return &broker{
		subscriber: make(map[string]*bufio.Writer),
		notifier:   notifier,
	}
}

func (b *broker) Subscribe(writer *bufio.Writer) string {
	id := fiberUtils.UUID()
	b.subscriber[id] = writer
	utils.Logger.Info("--> Subscribed: " + id)
	return id
}
func (b *broker) Unsubscribe(id string) {
	delete(b.subscriber, id)
	utils.Logger.Info("--> Unsubscribed: " + id)
}
func (b *broker) CountSubscriber() int64 {
	return int64(len(b.subscriber))
}
func (b *broker) GetSubscriber() map[string]*bufio.Writer {
	return b.subscriber
}
func (b *broker) Listen() {
	for {
		select {
		case msg := <-b.notifier.Get():
			errList := b.notifier.Send(b.subscriber, msg)
			for id := range errList {
				b.Unsubscribe(id)
			}
		}
	}
}
