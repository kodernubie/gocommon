package mq

type RabbitConnection struct {
}

func (o *RabbitConnection) Publish(queueName string, payload []byte) {

}

func (o *RabbitConnection) Subscribe(queueName string, handler ItemHandler) {

}

func (o *RabbitConnection) GroupSubscribe(queueName, groupName string, handler ItemHandler) {

}

func init() {

	RegConnectionCreator("rabbit", func(configName string) (Connection, error) {

		ret := &RabbitConnection{}

		return ret, nil
	})
}
