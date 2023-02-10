package sse

var (
	NotifierController Notifier
	BrokerList         Broker
)

func init() {
	NotifierController = NewNotifier()
	BrokerList = NewBroker(NotifierController)
}
