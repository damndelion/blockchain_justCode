package applicator

import (
	"fmt"
	"github.com/evrone/go-clean-template/config/blockchain"
	_ "github.com/evrone/go-clean-template/config/blockchain"
	"github.com/evrone/go-clean-template/internal/blockchain/usecase"
	"github.com/evrone/go-clean-template/internal/blockchain/usecase/repo"
	v1 "github.com/evrone/go-clean-template/internal/controller/http/v1"
	blockchain_logic2 "github.com/evrone/go-clean-template/pkg/blockchain_logic"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
)

// Run creates objects via constructors.
func Run(cfg *blockchain.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	db, _, err := postgres.New(cfg.PG.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("user - Run - postgres.New: %w", err))
	}
	defer db.Close()

	address := blockchain_logic2.CreateWallet()

	// Use case
	chainUseCase := usecase.NewBlockchain(repo.NewBlockchainRepo(db, address))

	blockchain_logic2.ListAddresses()
	//address to create genesis block
	chain := blockchain_logic2.CreateBlockchain(db, address)

	// HTTP Server
	handler := gin.New()
	v1.NewBlockchainRouter(handler, l, chainUseCase, *chain)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("user - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("user - Run - httpServer.Notify: %w", err))

		// Shutdown
		err = httpServer.Shutdown()
		if err != nil {
			l.Error(fmt.Errorf("user - Run - httpServer.Shutdown: %w", err))
		}
	}
}
