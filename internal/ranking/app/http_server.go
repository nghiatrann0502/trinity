package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/nghiatrann0502/trinity/cmd/ranking/docs"
	"github.com/nghiatrann0502/trinity/internal/ranking/adapters/ginhttp"
	"github.com/nghiatrann0502/trinity/internal/ranking/config"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitGinRouter(handler ginhttp.GinHandler) *gin.Engine {
	router := gin.Default()

	// Define the routes
	router.GET("/ping", handler.Ping())
	router.GET("/health", handler.Health())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	// TODO: Add more routes here
	v1 := router.Group("/v1")
	{
		v1.POST("/videos/:id/score", handler.UpdateVideoScore())
		v1.GET("/videos/ranked", handler.GetTopRanked())
	}

	return router
}

func NewHTTPServer(cfg *config.Config, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
