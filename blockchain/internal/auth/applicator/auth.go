package applicator

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/evrone/go-clean-template/config/auth"
	consumer "github.com/evrone/go-clean-template/internal/auth/consumer"
	v1 "github.com/evrone/go-clean-template/internal/auth/controller/http/v1"
	authEntity "github.com/evrone/go-clean-template/internal/auth/entity"
	"github.com/evrone/go-clean-template/internal/auth/transport"
	"github.com/evrone/go-clean-template/internal/auth/usecase"
	"github.com/evrone/go-clean-template/internal/auth/usecase/repo"
	natsService "github.com/evrone/go-clean-template/internal/nats"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/jaeger"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go"
)

// Run creates objects via constructors.
func Run(cfg *auth.Config) {
	l := logger.New(cfg.Log.Level)

	tracer, closer, err := jaeger.InitJaeger("auth-service", cfg.Jaeger.URL)
	if err != nil {
		l.Error(fmt.Errorf("auth - Run - jaeger.InitJaeger: %w", err))
	}
	defer func(closer io.Closer) {
		err = closer.Close()
		if err != nil {
			l.Error("Failed to close Jaeger: %v", err)
		}
	}(closer)
	opentracing.SetGlobalTracer(tracer)

	db, _, err := postgres.New(cfg.PG.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("auth - Run - postgres.New: %w", err))
	}
	sqlDB, err := db.DB()
	defer func(sqlDB *sql.DB) {
		err = sqlDB.Close()
		if err != nil {
			l.Error("Failed to close DB: %v", err)
		}
	}(sqlDB)

	nc, err := nats.Connect(cfg.Nats.Server)
	if err != nil {
		l.Error("failed to connect to NATS server: %v", err)

		return
	}
	defer nc.Close()

	userVerificationProducer, err := natsService.NewProducer(cfg)
	if err != nil {
		l.Error("Failed to create NATS producer: %v", err)
	}
	userGrpcTransport := transport.NewUserGrpcTransport(cfg.Transport.UserGrpc)

	userVerificationConsumerCallback := consumer.NewUserVerificationCallback(l, db, nc, userGrpcTransport)
	userVerificationConsumer, err := natsService.NewConsumer(l, cfg, userVerificationConsumerCallback)
	if err != nil {
		l.Fatal("Failed to create NATS consumer: %v", err)
	}
	go userVerificationConsumer.Start()

	err = db.AutoMigrate(authEntity.Token{})
	if err != nil {
		l.Error("Failed to do migrations Token: %v", err)
	}
	err = db.AutoMigrate(authEntity.UserCode{})
	if err != nil {
		l.Error("Failed to do migrations UserCode: %v", err)
	}
	authUseCase := usecase.NewAuth(repo.NewAuthRepo(db, userGrpcTransport), cfg, userVerificationProducer)

	handler := gin.New()
	v1.NewAuthRouter(handler, l, authUseCase)

	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("auth - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("auth - Run - httpServer.Notify: %w", err))

		err = httpServer.Shutdown()
		if err != nil {
			l.Error(fmt.Errorf("auth - Run - httpServer.Shutdown: %w", err))
		}
	}
}
