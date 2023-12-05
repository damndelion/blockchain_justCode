package nats

import (
	"encoding/json"

	"github.com/damndelion/blockchain_justCode/internal/auth/consumer/dto"
	"github.com/damndelion/blockchain_justCode/internal/auth/transport"
	"github.com/damndelion/blockchain_justCode/pkg/logger"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type UserVerificationCallback struct {
	logger logger.Interface
	db     *gorm.DB
	nc     *nats.Conn
	ut     *transport.UserGrpcTransport
}

func NewUserVerificationCallback(logger logger.Interface, db *gorm.DB, nc *nats.Conn, ut *transport.UserGrpcTransport) *UserVerificationCallback {
	return &UserVerificationCallback{logger: logger, db: db, nc: nc, ut: ut}
}

func (c *UserVerificationCallback) Callback(msg *nats.Msg) {
	var userVerification dto.UserCode

	err := json.Unmarshal(msg.Data, &userVerification)
	if err != nil {
		c.logger.Error("failed to unmarshal record value: %v", err)
	} else {
		c.logger.Info("user code: %s", userVerification.Code)
		//send the code to user email
		if err := c.db.Create(&userVerification).Error; err != nil {
			c.logger.Error("failed to save user verification code: %v", err)
		}
	}
}
