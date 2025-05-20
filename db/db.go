package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"crawler/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func Connect(cfg *config.DBConfig) {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?maxPoolSize=20",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Optional: Ping the DB to make sure it connected
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	MongoClient = client

	log.Println("Database connection established successfully.")
}
