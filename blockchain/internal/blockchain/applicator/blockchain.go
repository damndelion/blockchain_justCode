package applicator

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evrone/go-clean-template/pkg/cache"

	"github.com/evrone/go-clean-template/config/blockchain"
	_ "github.com/evrone/go-clean-template/config/blockchain"
	v1 "github.com/evrone/go-clean-template/internal/blockchain/controller/http/v1"
	"github.com/evrone/go-clean-template/internal/blockchain/transport"
	"github.com/evrone/go-clean-template/internal/blockchain/usecase"
	"github.com/evrone/go-clean-template/internal/blockchain/usecase/repo"
	blockchainlogic "github.com/evrone/go-clean-template/pkg/blockchain_logic"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
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
