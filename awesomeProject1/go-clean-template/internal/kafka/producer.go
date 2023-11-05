package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/evrone/go-clean-template/config/auth"
)

type Producer struct {
	asyncProducer sarama.AsyncProducer
	topic         string
}

func NewProducer(cfg *auth.Config) (*Producer, error) {
	samaraConfig := sarama.NewConfig()

	asyncProducer, err := sarama.NewAsyncProducer(cfg.Brokers, samaraConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to NewAsyncProducer err: %w", err)
	}

	return &Producer{
		asyncProducer: asyncProducer,
		topic:         cfg.Producer.Topic,
	}, nil
}

func (p *Producer) ProduceMessage(message []byte) {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(message),
	}

	p.asyncProducer.Input() <- msg
}
