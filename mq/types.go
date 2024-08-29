package mq

type ItemHandler func(queueName string, payload []byte)

type Connection interface {
	Publish(queueName string, payload []byte)
	Subscribe(queueName string, handler ItemHandler)
	GroupSubscribe(queueName, groupName string, handler ItemHandler)
}

type ConnectionCreator func(configName string) (Connection, error)
