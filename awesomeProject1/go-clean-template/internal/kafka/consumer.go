package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/evrone/go-clean-template/config/auth"
	"github.com/evrone/go-clean-template/pkg/logger"
	"strings"
)

type ConsumerCallback interface {
	Callback(message <-chan *sarama.ConsumerMessage, error <-chan *sarama.ConsumerError)
}

type Consumer struct {
	logger   logger.Interface
	topics   []string
	master   sarama.Consumer
	callback ConsumerCallback
}

func NewConsumer(
	logger logger.Interface,
	cfg *auth.Config,
	callback ConsumerCallback,
) (*Consumer, error) {
	samaraCfg := sarama.NewConfig()
	samaraCfg.ClientID = "go-kafka-consumer"
	samaraCfg.Consumer.Return.Errors = true

	master, err := sarama.NewConsumer(cfg.Brokers, samaraCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create NewConsumer err: %w", err)
	}

	return &Consumer{
		logger:   logger,
		topics:   cfg.Consumer.Topics,
		master:   master,
		callback: callback,
	}, nil
}

func (c *Consumer) Start() {
	consumers := make(chan *sarama.ConsumerMessage, 1)
	errors := make(chan *sarama.ConsumerError)

	for _, topic := range c.topics {
		if strings.Contains(topic, "__consumer_offsets") {
			continue
		}

		partitions, _ := c.master.Partitions(topic)

		consumer, err := c.master.ConsumePartition(topic, partitions[0], sarama.OffsetNewest)
		if nil != err {
			c.logger.Error("Topic %v Partitions: %v, err: %w", topic, partitions, err)
			continue
		}

		c.logger.Info(" Start consuming topic ", topic)

		go func(topic string, consumer sarama.PartitionConsumer) {
			for {
				select {
				case consumerError := <-consumer.Errors():
					errors <- consumerError

				case msg := <-consumer.Messages():
					consumers <- msg
				}
			}
		}(topic, consumer)
	}

	c.callback.Callback(consumers, errors)
}
