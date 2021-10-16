package review

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/card"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Review struct {
	Id            *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OriginType    string              `json:"originType" bson:"originType"`
	OriginId      string              `json:"originId" bson:"originId"`
	UserId        string              `json:"userId" bson:"userId"`
	Hists         []*card.Card        `json:"hists" bson:"hists"`
	HistsCount    int64               `json:"histsCount" bson:"histsCount"`
	Mistakes      []*card.Card        `json:"mistakes" bson:"mistakes"`
	MistakesCount int64               `json:"mistakesCount" bson:"mistakesCount"`
	LastUpdate    time.Time           `json:"lastUpdate" bson:"lastUpdate"`
}
