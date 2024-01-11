package rabbit

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Sender struct {
	channel *amqp091.Channel
}

func NewSender(channel *amqp091.Channel) Sender {
	return Sender{channel: channel}
}

func (s *Sender) SendMessage(ctx context.Context, message protoreflect.ProtoMessage, exchange, rk string) error {
	bytes, err := proto.Marshal(message)

	if err != nil {
		return err
	}

	err = s.channel.PublishWithContext(ctx, exchange, rk, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        bytes,
	})

	if err != nil {
		return err
	}
	return nil
}
