package deck

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Deck struct {
	Id               *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ImageURL         string              `json:"imageURL" bson:"imageURL" validate:"required,max=100"`
	Name             string              `json:"name" bson:"name" validate:"required,max=100"`
	Description      string              `json:"description" bson:"description" validate:"required,max=400"`
	IsPrivate        bool                `json:"isPrivate" bson:"isPrivate"`
	StudySuggestions []string            `json:"studySuggestions" bson:"studySuggestions"`
	CardsCount       int64               `json:"cardsCount" bson:"cardsCount"`
	UserId           string              `json:"userId" bson:"userId"`
	LastUpdate       time.Time           `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
}
