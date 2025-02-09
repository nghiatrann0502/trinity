package app

import (
	"fmt"
	"net"

	"github.com/nghiatrann0502/trinity/internal/video/adapters/repositories"
	"github.com/nghiatrann0502/trinity/internal/video/config"
	"github.com/nghiatrann0502/trinity/internal/video/core/services"
	"github.com/nghiatrann0502/trinity/pkg/logger"
	"github.com/nghiatrann0502/trinity/pkg/redisc"
	"google.golang.org/grpc"
)

type app struct {
	cfg *config.Config
	log logger.Logger

	grpcServer *grpc.Server
}

func NewApp(cfg *config.Config, log logger.Logger) (*app, func(), error) {
	// Connect to MySQL
	// tcpUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)
	// db, err := database.NewMySQLDatabase(tcpUrl)
	// if err != nil {
	// 	log.Fatal("cannot connect to mysql", err, nil)
	// }
	//
	// ok, state := db.Health()
	// if !ok {
	// 	log.Fatal("cannot connect to mysql", errors.New(state["error"]), nil)
	// } else {
	// 	log.Info("connected to mysql", nil)
	// }

	// Connect to redis
	redisUrl := fmt.Sprintf("redis://:%s@%s:%d", cfg.Redis.Password, cfg.Redis.Host, cfg.Redis.Port)
	redis, err := redisc.NewRedisClient(redisUrl)
	if err != nil {
		log.Fatal("cannot connect to redis", err, nil)
	}

	// redis health check
	ok, err := redis.Health()
	if err != nil {
		log.Fatal("cannot connect to redis", err, nil)
	}
	if ok {
		log.Info("connected to redis", nil)
	}

	// repository := repositories.NewMysqlRepo(log, db)
	repository := repositories.NewMemoryRepo()
	service := services.NewService(log, repository)

	grpcServer := NewGRPCServer(log, service)

	return &app{
			cfg:        cfg,
			log:        log,
			grpcServer: grpcServer,
		}, func() {
			grpcServer.Stop()
			redis.Close()
		}, nil
}

func (a *app) Run() error {
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.GRPC.Port))
		if err != nil {
			a.log.Fatal("Failed to listen", err, nil)
		}
		a.log.Info(fmt.Sprintf("Starting gRPC server on :%d", a.cfg.GRPC.Port), nil)
		if err := a.grpcServer.Serve(lis); err != nil {
			a.log.Fatal("Failed to serve", err, nil)
		}
	}()

	return nil
}
