package broker

import (
	"fmt"

	"jobsity-code-challenge/config"

	"github.com/streadway/amqp"
)

type Broker struct {
	amqpURL string
	conn    *amqp.Connection
	queues  Queues
}

type Queues struct {
	stockQueueName string
	stockQueue     amqp.Queue
}

func New(conf config.Broker) *Broker {
	conn, err := amqp.Dial(conf.AmqpUrl)
	if err != nil {
		panic(fmt.Sprintf("error connecting to rabbitmq %v\n", err))
	}
	c, err := conn.Channel()
	defer c.Close()
	if err != nil {
		panic(fmt.Sprintf("error ccreating rabbitmq channel %v\n", err))
	}
	q, err := c.QueueDeclare(conf.StockQueueName, true, false, false, false, nil)
	if err != nil {
		panic(fmt.Sprintf("error configuring rabbitmq %v\n", err))
	}
	return &Broker{
		amqpURL: conf.AmqpUrl,
		conn:    conn,
		queues:  Queues{stockQueueName: conf.StockQueueName, stockQueue: q},
	}
}

func (b Broker) Close() error {
	return b.conn.Close()
}

func (b Broker) channel() (*amqp.Channel, error) {
	if b.conn == nil {
		c, err := amqp.Dial(b.amqpURL)
		if err != nil {
			return nil, err
		}
		b.conn = c
	}
	return b.conn.Channel()
}

func (b Broker) GetStockQueueName() string {
	return b.queues.stockQueue.Name
}

func (b Broker) Publish(key string, msg []byte) error {
	channel, err := b.channel()
	if err != nil {
		return err
	}
	return channel.Publish(
		"",
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
}

func (b Broker) Consume(key string) (<-chan amqp.Delivery, error) {
	channel, err := b.channel()
	if err != nil {
		return nil, err
	}
	return channel.Consume(
		key,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}
