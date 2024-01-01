package infrastructure

import (
	"context"
	"fmt"
	"go-learning/clean-architecture-mongo/database"
	"log"
	"os"
	"time"
)

func NewMongoDatabase() database.Client {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	mongodbURL := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)
	log.Println("mongodbURL: ", mongodbURL)

	if dbUser == "" || dbPass == "" {
		mongodbURL = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	}

	client, err := database.NewClient(ctx, mongodbURL)

	if err != nil {
		log.Fatal("err")
	}

	err = client.Ping(ctx)

	if err != nil {
		log.Fatal(err)
	}

	return client
}

func CloseMongoConnection(client database.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}
