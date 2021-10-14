package mongodb

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"os"
	"time"
)

const (
	DESC = -1
	ASC  = 1
)

func NewMongoDbClient() (*mongo.Client, error) {
	// Set client options
	address := os.Getenv("MONGODB_ADDRESS")
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")

	clientOptions := options.Client().ApplyURI(address)

	if username != "" && password != "" {
		clientOptions.
			SetAuth(options.Credential{
				Username: username,
				Password: password,
			})
	}

	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Logger.Errorw("Error on create MongoDB Client", "error", err.Error())
		return &mongo.Client{}, errors.New("Error on create MongoDB Client: " + err.Error())
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Logger.Errorw("Connect failed: connection was not established", "error", err.Error())
		return &mongo.Client{}, errors.New("Error to create MongoDB connection: " + err.Error())
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Logger.Errorw("Ping failed: connection was not established", "error", err.Error())
		return &mongo.Client{}, errors.New("Error on check MongoDB connection: " + err.Error())
	}

	return client, nil
}
