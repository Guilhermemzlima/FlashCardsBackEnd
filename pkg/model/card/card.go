package card

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/Enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Card struct {
	Id         *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Front      string              `json:"front" bson:"front" validate:"required,max=400"`
	Back       string              `json:"back" bson:"back" validate:"max=400"`
	UserId     string              `json:"userId" bson:"userId"`
	Color      Enums.Color         `json:"color" bson:"color"`
	DeckId     string              `json:"deckId" bson:"deckId"`
	IsPrivate  bool                `json:"isPrivate" bson:"isPrivate"`
	LastUpdate time.Time           `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
}
