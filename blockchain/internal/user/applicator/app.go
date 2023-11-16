package applicator

import (
	"database/sql"
	"fmt"
	"github.com/evrone/go-clean-template/config/user"
	"github.com/evrone/go-clean-template/internal/user/controller/grpc"
	v1 "github.com/evrone/go-clean-template/internal/user/controller/http/v1"
	userEntity "github.com/evrone/go-clean-template/internal/user/entity"
	"github.com/evrone/go-clean-template/internal/user/usecase"
	"github.com/evrone/go-clean-template/internal/user/usecase/repo"
	"github.com/evrone/go-clean-template/pkg/cache"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/gin-gonic/gin"
)

func Runn(cfg *user.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository

	db, _, err := postgres.New(cfg.PG.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("user - Run - postgres.New: %w", err))
	}
	sqlDB, err := db.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			l.Fatal(err)
		}
	}(sqlDB)

	//Redis client
	redisClient, err := cache.NewRedisClient()
	if err != nil {
		l.Error(err)
	}

	userCache := cache.NewUserCache(redisClient, 10*time.Minute)

	// Use case
	userRepo := repo.NewUserRepo(db)
	userUseCase := usecase.NewUser(userRepo, cfg)

	err = db.AutoMigrate(&userEntity.User{})
	if err != nil {
		l.Error(err)
	}
	err = db.AutoMigrate(&userEntity.UserInfo{})
	if err != nil {
		l.Error(err)

	}
	err = db.AutoMigrate(&userEntity.UserCredentials{})
	if err != nil {
		l.Error(err)

	}
	// HTTP Server

	handler := gin.New()
	v1.NewUserRouter(handler, l, userUseCase, userCache, cfg)

	grpcService := grpc.NewService(l, userRepo)
	grpcServer := grpc.NewServer(cfg.GrpcServer.Port, grpcService)
	err = grpcServer.Start()
	if err != nil {
		l.Fatal("failed to start grpc-server err: %v", err)
	}

	defer grpcServer.Close()

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
