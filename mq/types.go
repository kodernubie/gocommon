package mq

type ItemHandler func(queueName string, payload []byte)

type Connection interface {
	Publish(queueName string, payload []byte) error
	Subscribe(queueName string, handler ItemHandler) error
	GroupSubscribe(queueName, groupName string, handler ItemHandler) error
}

type ConnectionCreator func(configName string) (Connection, error)
