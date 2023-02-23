package main

import (
	"context"
	"gophermart/internal/core/services/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	handler := createHandler()

	srv := server.NewServer(":8080", handler)

	srv.Start()

	waitSigterm()

	srv.Stop(context.Background())
}

func createHandler() http.Handler {
	r := gin.Default()
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
