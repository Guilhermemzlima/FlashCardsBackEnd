package utils

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"os"
)

func InitMongoContainer() {
	playlistsColl := "playlists"
	decksColl := "decks"
	cardsColl := "card"
	database := "flashcards"

	variables := make(map[string]string)
	variables["MONGO_INITDB_DATABASE"] = database

	req := testcontainers.ContainerRequest{
		Image:        "mongo",
		ExposedPorts: []string{"27017/tcp"},
		SkipReaper:   true,
		Env:          variables,
	}

	var err error
	var ctx = context.Background()

	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	port, err := mongoContainer.MappedPort(ctx, "27017")
	if err != nil {
		panic(err)
	}

	host, err := mongoContainer.Host(ctx)

	_ = os.Setenv("MONGODB_ADDRESS", "mongodb://"+host+":"+port.Port()+"/admin")
	_ = os.Setenv("MONGODB_DATABASE", database)
	_ = os.Setenv("MONGODB_DECK_COLLECTION", decksColl)
	_ = os.Setenv("MONGODB_CARD_COLLECTION", cardsColl)
	_ = os.Setenv("MONGODB_PLAYLIST_COLLECTION", playlistsColl)
}
