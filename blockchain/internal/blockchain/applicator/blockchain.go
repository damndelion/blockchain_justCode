package applicator

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/damndelion/blockchain_justCode/pkg/jaeger"
	"github.com/opentracing/opentracing-go"

	"github.com/damndelion/blockchain_justCode/pkg/cache"

	"github.com/damndelion/blockchain_justCode/config/blockchain"
	_ "github.com/damndelion/blockchain_justCode/config/blockchain"
	v1 "github.com/damndelion/blockchain_justCode/internal/blockchain/controller/http/v1"
	"github.com/damndelion/blockchain_justCode/internal/blockchain/transport"
	"github.com/damndelion/blockchain_justCode/internal/blockchain/usecase"
	"github.com/damndelion/blockchain_justCode/internal/blockchain/usecase/repo"
	blockchainlogic "github.com/damndelion/blockchain_justCode/pkg/blockchain_logic"
	"github.com/damndelion/blockchain_justCode/pkg/httpserver"
	"github.com/damndelion/blockchain_justCode/pkg/logger"
	"github.com/damndelion/blockchain_justCode/pkg/postgres"
	"github.com/gin-gonic/gin"
)

// Run creates objects via constructors.
func Run(cfg *blockchain.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	_, db, err := postgres.New(cfg.PG.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("blockchain - Run - postgres.New: %w", err))
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			l.Error(fmt.Errorf("blockchain - DB: %w", err))
		}
	}(db)

	tracer, closer, err := jaeger.InitJaeger("blockchain-service", cfg.Jaeger.URL)
	if err != nil {
		l.Error(fmt.Errorf("blockchain - Run - jaeger.InitJaeger: %w", err))
	}
	defer func(closer io.Closer) {
		err = closer.Close()
		if err != nil {
			l.Error("Failed to close Jaeger: %v", err)
		}
	}(closer)
	opentracing.SetGlobalTracer(tracer)

	address := blockchainlogic.CreateWallet()

	userGrpcTransport := transport.NewUserGrpcTransport(cfg.Transport.UserGrpc)

	// Use case
	chainUseCase := usecase.NewBlockchain(repo.NewBlockchainRepo(db, address, userGrpcTransport), cfg, userGrpcTransport)
	blockchainlogic.ListAddresses()
	// address to create genesis block
	chain := blockchainlogic.CreateBlockchain(db, address)

	redisClient, err := cache.NewRedisClient(cfg.Redis.Host)
	blockchainCache := cache.NewBlockchainCache(redisClient, 10*time.Minute)

	// HTTP Server
	handler := gin.New()
	v1.NewBlockchainRouter(handler, l, chainUseCase, *chain, cfg, blockchainCache)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("user - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("blockchain - Run - httpServer.Notify: %w", err))

		// Shutdown
		err = httpServer.Shutdown()
		if err != nil {
			l.Error(fmt.Errorf("blockchain - Run - httpServer.Shutdown: %w", err))
		}
	}
}
