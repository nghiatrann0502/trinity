package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nghiatrann0502/trinity/internal/ranking/adapters/ginhttp"
	"github.com/nghiatrann0502/trinity/internal/ranking/adapters/repositories"
	"github.com/nghiatrann0502/trinity/internal/ranking/config"
	"github.com/nghiatrann0502/trinity/internal/ranking/core/services"
	"github.com/nghiatrann0502/trinity/pkg/database"
	"github.com/nghiatrann0502/trinity/pkg/logger"
	"github.com/nghiatrann0502/trinity/pkg/redisc"
	"github.com/nghiatrann0502/trinity/proto/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type app struct {
	cfg        *config.Config
	log        logger.Logger
	db         database.DBEngine
	httpServer *http.Server
}

func NewApp(cfg *config.Config, log logger.Logger) (*app, func(), error) {
	// Connect to MySQL
	tcpUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)
	db, err := database.NewMySQLDatabase(tcpUrl)
	if err != nil {
		log.Fatal("cannot connect to mysql", err, nil)
	}

	ok, state := db.Health()
	if !ok {
		log.Fatal("cannot connect to mysql", errors.New(state["error"]), nil)
	} else {
		log.Info("connected to mysql", nil)
	}

	cur, _ := os.Getwd()

	isDev := cfg.App.Production == false
	var dir string
	dir = filepath.Dir(cur + "/")
	if isDev {
		dir = filepath.Dir(cur + "../../../")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}

	for _, e := range entries {
		fmt.Println(e.Name())
	}

	if err := db.Migrate(dir); err != nil {
		log.Fatal("cannot migrate database", err, nil)
	}

	// Connect to redis
	redisUrl := fmt.Sprintf("redis://:%s@%s:%d", cfg.Redis.Password, cfg.Redis.Host, cfg.Redis.Port)
	redis, err := redisc.NewRedisClient(redisUrl)
	if err != nil {
		log.Fatal("cannot connect to redis", err, nil)
	}

	// redis health check
	ok, err = redis.Health()
	if err != nil {
		log.Fatal("cannot connect to redis", err, nil)
	}
	if ok {
		log.Info("connected to redis", nil)
	}

	// Connect to grpc
	videoGrpcUrl := fmt.Sprintf("%s:%d", cfg.VideoGrpc.Host, cfg.VideoGrpc.Port)
	conn, err := grpc.NewClient(*&videoGrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot connect to grpc", err, nil)
	}
	client := proto.NewVideoServiceClient(conn)

	// Create a repository
	repository := repositories.NewMySQLRepository(log, db)
	cacheRepository := repositories.NewRedisCacheRepository(log, redis)
	grpcRepository := repositories.NewGrpcVideoRepository(log, client)

	service := services.NewRankingService(log, repository, cacheRepository, grpcRepository)
	handler := ginhttp.NewGinHandler(service, log)

	router := InitGinRouter(handler)
	httpServer := NewHTTPServer(cfg, router)

	return &app{
			cfg:        cfg,
			log:        log,
			db:         db,
			httpServer: httpServer,
		}, func() {
			// Clean up function
			// shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			// defer cancel()

			// if err := httpServer.Shutdown(shutdownCtx); err != nil {
			// 	log.Fatal("http server shutdown error", err, nil)
			// }
			conn.Close()
			db.Close()
			redis.Close()
		}, nil
}

func (a *app) Run() error {
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			a.log.Fatal("http server error", err, nil)
		}
	}()

	return nil
}
