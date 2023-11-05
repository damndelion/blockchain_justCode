package consumer

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/evrone/go-clean-template/internal/auth/consumer/dto"
	"github.com/evrone/go-clean-template/pkg/logger"
	"gorm.io/gorm"
)

type UserVerificationCallback struct {
	logger logger.Interface
	db     *gorm.DB
}

func NewUserVerificationCallback(logger logger.Interface, db *gorm.DB) *UserVerificationCallback {
	return &UserVerificationCallback{logger: logger, db: db}
}

func (c *UserVerificationCallback) Callback(message <-chan *sarama.ConsumerMessage, error <-chan *sarama.ConsumerError) {
	for {
		select {
		case msg := <-message:
			var userVerification dto.UserCode

			err := json.Unmarshal(msg.Value, &userVerification)
			if err != nil {
				c.logger.Error("failed to unmarshall record value err: %v", err)
			} else {
				c.logger.Info("user code: %s", userVerification.Code)
				//TODO grpc
				if err := c.db.Create(&userVerification).Error; err != nil {
					c.logger.Error("failed to save user verification code: %v", err)
				}
			}
		case err := <-error:
			c.logger.Error("failed consume err: %v", err)
		}
	}
}
