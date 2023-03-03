package main

import (
	"context"
	"gophermart/internal/api/userapi"
	"gophermart/internal/core/services/config"
	"gophermart/internal/core/services/db"
	"gophermart/internal/core/services/logging"
	"gophermart/internal/core/services/server"
	"gophermart/internal/core/services/userservice"
	"gophermart/internal/core/stores/userstore"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	logService := logging.New()
	mainLogger := logService.ComponentLogger("Main")

	conf, err := config.GetConfig()
	if err != nil {
		mainLogger.Fatal().Err(err).Msg("failed to read configuration")
	}

	db, err := db.NewDB(conf.DatabaseDSN)
	if err != nil {
		mainLogger.Fatal().Err(err).Msg("failed to initiate database")
	}

	engine := createEngine()

	// Stores
	userStore := userstore.New(db)

	// Services
	srv := server.NewServer(":8080", engine, logService)
	userService := userservice.New(conf.Secret, logService, userStore)

	// APIs
	userAPI := userapi.New(logService, userService)
	userAPI.Register(engine)

	// Start
	srv.Start()

	waitSigterm()

	srv.Stop(context.Background())
}

func createEngine() *gin.Engine {
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
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server was interrupted. Shutdown server.")
}
