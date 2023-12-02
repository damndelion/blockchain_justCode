package nats

import (
	"github.com/damndelion/blockchain_justCode/config/auth"
	consumer "github.com/damndelion/blockchain_justCode/internal/auth/consumer"
	"github.com/damndelion/blockchain_justCode/pkg/logger"
	"github.com/nats-io/nats.go"
)

type Consumer struct {
	logger   logger.Interface
	topics   []string
	nc       *nats.Conn
	callback *consumer.UserVerificationCallback
}

func NewConsumer(
	logger logger.Interface,
	cfg *auth.Config,
	callback *consumer.UserVerificationCallback,
) (*Consumer, error) {
	nc, err := nats.Connect(cfg.Server)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		logger:   logger,
		topics:   cfg.Consumer.Topics,
		nc:       nc,
		callback: callback,
	}, nil
}

func (c *Consumer) Start() {
	for _, topic := range c.topics {
		// Subscribe to NATS subjects here
		_, err := c.nc.Subscribe(topic, c.callback.Callback)
		if err != nil {
			c.logger.Error("Failed to subscribe to subject %s: %v", topic, err)
		}
		c.logger.Info("Start consuming topic ", topic)
	}
}
