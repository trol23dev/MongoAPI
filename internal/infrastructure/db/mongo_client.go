package db

import (
	"context"
	"errors"
	"time"

	"mongodb-api/internal/infrastructure/logger"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoClient struct {
	Client *mongo.Client
}

func NewMongoClient(uri string) (*MongoClient, error) {
	if uri == "" {
		logger.New(logger.Error).Print("🍃 MongoDB URI is empty. Please set MONGO_URI in .env")
		return nil, errors.New("MONGO_URI is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		logger.New(logger.Error).Print("🍃 MongoDB connection has been failed: " + err.Error())
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		logger.New(logger.Error).Print("🍃 MongoDB ping failed: " + err.Error())
		return nil, err
	}

	logger.New(logger.Success).Print("🍃 MongoDB connection has been established.")
	return &MongoClient{Client: client}, nil
}

func (m *MongoClient) Close(ctx context.Context) error {
	if m.Client == nil {
		return nil
	}
	err := m.Client.Disconnect(ctx)
	if err != nil {
		logger.New(logger.Error).Print("🍃 Error closing MongoDB connection: " + err.Error())
		return err
	}
	logger.New(logger.System).Print("🍃 MongoDB connection has been closed.")
	return nil
}

func (m *MongoClient) GetDatabase(databaseName string) *mongo.Database {
	return m.Client.Database(databaseName)
}

func (m *MongoClient) GetCollection(databaseName, collectionName string) *mongo.Collection {
	return m.Client.Database(databaseName).Collection(collectionName)
}
