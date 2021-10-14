package playlist

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/deck"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Playlist struct {
	Id               *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ImageURL         string              `json:"imageURL" bson:"imageURL" validate:"required,max=100"`
	Name             string              `json:"name" bson:"name" validate:"required,max=100"`
	Description      string              `json:"description" bson:"description" validate:"required,max=400"`
	IsPrivate        bool                `json:"isPrivate" bson:"isPrivate"`
	StudySuggestions []string            `json:"studySuggestions" bson:"studySuggestions"`
	Decks            []deck.DeckPreview  `json:"decks" bson:"decks"`
	UserId           string              `json:"userId" bson:"userId"`
	LastUpdate       time.Time           `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
}
