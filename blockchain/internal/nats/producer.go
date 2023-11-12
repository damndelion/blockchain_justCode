package nats

import (
	"github.com/evrone/go-clean-template/config/auth"
	"github.com/nats-io/nats.go"
)

type Producer struct {
	conn  *nats.Conn
	topic string
}

func NewProducer(cfg *auth.Config) (*Producer, error) {
	nc, err := nats.Connect(cfg.Server)
	if err != nil {
		return nil, err
	}

	return &Producer{
		conn:  nc,
		topic: cfg.Producer.Topic,
	}, nil
}

func (p *Producer) ProduceMessage(message []byte) error {
	err := p.conn.Publish(p.topic, message)

	if err != nil {
		return err
	}
	return nil
}
