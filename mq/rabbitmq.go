package mq

import (
	"sync"

	"github.com/kodernubie/gocommon/conf"
	"github.com/kodernubie/gocommon/log"
	"github.com/oklog/ulid/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

var exchanges map[string]bool = map[string]bool{}
var queues map[string]bool = map[string]bool{}
var rabbitLock sync.RWMutex

type RabbitConnection struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func (o *RabbitConnection) makeSureExchange(name string) error {

	rabbitLock.RLock()
	declared := exchanges[name]
	rabbitLock.RUnlock()

	if !declared {
		rabbitLock.Lock()
		defer rabbitLock.Unlock()

		declared = exchanges[name]

		if !declared {

			err := o.ch.ExchangeDeclare(name,
				"fanout", // kind
				true,     // durable
				false,    // autodelete
				false,    // internal
				false,    // nowait
				nil)

			if err != nil {

				log.Fatal("error declaring exchange :", err)
				return err
			}

			exchanges[name] = true
		}
	}

	return nil
}

func (o *RabbitConnection) makeSureQueue(queueName, groupName string, exclusive bool) error {

	err := o.makeSureExchange(queueName)

	if err != nil {
		return err
	}

	targetQueueName := queueName + "_" + groupName

	rabbitLock.RLock()
	declared := queues[targetQueueName]
	rabbitLock.RUnlock()

	if !declared {

		rabbitLock.Lock()
		defer rabbitLock.Unlock()

		declared = queues[targetQueueName]

		if !declared {

			_, err = o.ch.QueueDeclare(queueName+"_"+groupName,
				true,      // durable
				false,     // autodelete
				exclusive, // exclusive
				false,     // no wait
				nil)

			if err != nil {
				return err
			}

			err = o.ch.QueueBind(targetQueueName, // destination
				"",        // key
				queueName, // source
				false,     // no wait
				nil)

			if err != nil {
				return err
			}

			queues[targetQueueName] = true
		}
	}

	return nil
}

func (o *RabbitConnection) Publish(queueName string, payload []byte) error {

	err := o.makeSureExchange(queueName)

	if err != nil {

		return err
	}

	return o.ch.Publish(queueName, // exchange
		"",    // key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        payload,
		})
}

func (o *RabbitConnection) Subscribe(queueName string, handler ItemHandler) error {

	groupName := ulid.Make().String()

	err := o.makeSureQueue(queueName, groupName, true)

	if err != nil {
		return err
	}

	msgs, err := o.ch.Consume(
		queueName+"_"+groupName,
		"",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		return err
	}

	go func() {

		for msg := range msgs {
			callHandler(queueName, handler, msg.Body)
		}
	}()

	return nil
}

func (o *RabbitConnection) GroupSubscribe(queueName, groupName string, handler ItemHandler) error {

	err := o.makeSureQueue(queueName, groupName, false)

	if err != nil {
		return err
	}

	msgs, err := o.ch.Consume(
		queueName+"_"+groupName,
		"",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		return err
	}

	go func() {

		for msg := range msgs {
			callHandler(queueName, handler, msg.Body)
		}
	}()

	return nil
}

func callHandler(queueName string, handler ItemHandler, payload []byte) {

	defer func() {

		err := recover()

		if err != nil {
			log.Error("RabbitMQ call handler recover because :", err)
		}
	}()

	handler(queueName, payload)
}

func init() {

	RegConnectionCreator("rabbit", func(configName string) (Connection, error) {

		ret := &RabbitConnection{}
		var err error

		ret.conn, err = amqp.Dial(conf.Str("MQ_"+configName+"_URL", "default"))

		if err != nil {
			return nil, err
		}

		ret.ch, err = ret.conn.Channel()

		return ret, err
	})
}
