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

	"github.com/damndelion/blockchain_justCode/config/user"
	"github.com/damndelion/blockchain_justCode/internal/user/controller/grpc"
	v1 "github.com/damndelion/blockchain_justCode/internal/user/controller/http/v1"
	userEntity "github.com/damndelion/blockchain_justCode/internal/user/entity"
	"github.com/damndelion/blockchain_justCode/internal/user/usecase"
	"github.com/damndelion/blockchain_justCode/internal/user/usecase/repo"
	"github.com/damndelion/blockchain_justCode/pkg/cache"

	"github.com/damndelion/blockchain_justCode/pkg/httpserver"
	"github.com/damndelion/blockchain_justCode/pkg/logger"
	"github.com/damndelion/blockchain_justCode/pkg/postgres"
	"github.com/gin-gonic/gin"
)

func Run(cfg *user.Config) {
	l := logger.New(cfg.Log.Level)

	db, _, err := postgres.New(cfg.PG.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("user - Run - postgres.New: %w", err))
	}
	sqlDB, err := db.DB()
	defer func(sqlDB *sql.DB) {
		err = sqlDB.Close()
		if err != nil {
			l.Fatal(err)
		}
	}(sqlDB)

	tracer, closer, err := jaeger.InitJaeger("user-service", cfg.Jaeger.URL)
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

	redisClient, err := cache.NewRedisClient(cfg.Redis.Host)

	if err != nil {
		l.Fatal(err)
	}

	userCache := cache.NewUserCache(redisClient, 10*time.Minute)

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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("user - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("user - Run - httpServer.Notify: %w", err))

		err = httpServer.Shutdown()
		if err != nil {
			l.Error(fmt.Errorf("user - Run - httpServer.Shutdown: %w", err))
		}
	}
}
