package main

import (
	"context"
	"gophermart/internal/core/services/config"
	"gophermart/internal/core/services/db"
	"gophermart/internal/core/services/logging"
	"gophermart/internal/core/services/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	logService := logging.New()
	logger := logService.ComponentLogger("Main")

	conf, err := config.GetConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to read configuration")
	}

	_, err = db.NewDB(conf.DatabaseDSN)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initiate database")
	}

	handler := createHandler()

	// Services
	srv := server.NewServer(":8080", handler, logService)

	srv.Start()

	waitSigterm()

	srv.Stop(context.Background())
}

func createHandler() http.Handler {
	r := gin.New()
	r.Use(gin.Recovery(), logger.SetLogger())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return r
}

func waitSigterm() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server was interrupted. Shutdown server.")
}
