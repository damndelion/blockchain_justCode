// Package user configures and runs application.
package applicator

import (
	"fmt"
	"github.com/evrone/go-clean-template/config/auth"
	"github.com/evrone/go-clean-template/internal/auth/consumer"
	"github.com/evrone/go-clean-template/internal/auth/controller/http/v1"
	"github.com/evrone/go-clean-template/internal/auth/usecase"
	"github.com/evrone/go-clean-template/internal/auth/usecase/repo"
	"github.com/evrone/go-clean-template/internal/kafka"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/jaeger"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"os"
	"os/signal"
	"syscall"
)

// Run creates objects via constructors.
func Run(cfg *auth.Config) {
	l := logger.New(cfg.Log.Level)

	//tracing
	tracer, closer, _ := jaeger.InitJaeger()
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// Repository
	db, _, err := postgres.New(cfg.PG.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("auth - Run - postgres.New: %w", err))
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()

	userVerificationProducer, err := kafka.NewProducer(cfg)
	if err != nil {
		l.Error("failed NewProducer err: %v", err)
	}

	userVerificationConsumerCallback := consumer.NewUserVerificationCallback(l, db)

	userVerificationConsumer, err := kafka.NewConsumer(l, cfg, userVerificationConsumerCallback)
	if err != nil {
		l.Fatal("failed NewConsumer err: %v", err)
	}

	go userVerificationConsumer.Start()

	// Use case
	authUseCase := usecase.NewAuth(repo.NewAuthRepo(db), cfg, userVerificationProducer)

	// HTTP Server
	handler := gin.New()
	v1.NewAuthRouter(handler, l, authUseCase)

	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("auth - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("auth - Run - httpServer.Notify: %w", err))

		// Shutdown
		err = httpServer.Shutdown()
		if err != nil {
			l.Error(fmt.Errorf("auth - Run - httpServer.Shutdown: %w", err))
		}
	}
}
