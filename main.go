package main

import (
	"github.com/gin-gonic/gin"
	"go-learning/clean-architecture-mongo/infrastructure"
	"go-learning/clean-architecture-mongo/infrastructure/routes"
	"os"
	"time"
)

func main() {
	config := infrastructure.LoadConfig()
	logger := infrastructure.NewLogger()
	client := infrastructure.NewMongoDatabase()
	defer infrastructure.CloseMongoConnection(client)
	db := client.Database("news")

	timeout := time.Duration(config.ContextTimeout) * time.Second
	app := gin.Default()

	routes.Dispatch(app, config, db, logger, timeout)

	app.Run(os.Getenv("SERVER_ADDRESS"))
}
