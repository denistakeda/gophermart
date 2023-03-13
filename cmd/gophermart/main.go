package main

import (
	"context"
	"gophermart/internal/api/userapi"
	"gophermart/internal/core/services/config"
	"gophermart/internal/core/services/db"
	"gophermart/internal/core/services/logging"
	"gophermart/internal/core/services/orderservice"
	"gophermart/internal/core/services/server"
	"gophermart/internal/core/services/userservice"
	"gophermart/internal/core/stores/orderstore"
	"gophermart/internal/core/stores/userstore"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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
	orderStore := orderstore.New(db)

	// Services
	srv := server.NewServer(":8080", engine, logService)
	userService := userservice.New(conf.Secret, logService, userStore)
	orderService := orderservice.New(logService, orderStore)

	// APIs
	userAPI := userapi.New(logService, userService, orderService)
	userAPI.Register(engine)

	// Start
	srv.Start()

	waitSigterm(mainLogger)

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

func waitSigterm(logger zerolog.Logger) {
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Warn().Msg("Server was interrupted. Shutdown server.")
}
