package Rabbit

import (
	"context"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type Mq struct {
	ID         uuid.UUID
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func New(url string) (*Mq, error) {
	connectRabbitMQ, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		return nil, err
	}
	return &Mq{
		ID:         uuid.New(),
		Connection: connectRabbitMQ,
		Channel:    channelRabbitMQ,
	}, nil
}
func (mq *Mq) DeclareQueue(name string) error {

	_, err := mq.Channel.QueueDeclare(
		name,  // queue name
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)
	return err
}
func (mq *Mq) Produce(queue, msg string) error {
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	}
	err := mq.Channel.Publish(
		"",
		queue,
		false,
		false,
		message,
	)

	return err
}

func (mq *Mq) Consumer(ctx context.Context, queue string, body *chan string, chanError chan error) {
	messages, err := mq.Channel.Consume(
		queue, // queue name
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		chanError <- err
		return
	}
	for {
		select {
		case msg := <-messages:
			*body <- string(msg.Body)
		case <-ctx.Done():
			return
		}
	}

}
