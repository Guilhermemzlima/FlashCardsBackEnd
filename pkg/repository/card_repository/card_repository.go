package card_repository

import (
	"FlashCardsBackEnd/internal/config/log"
	"FlashCardsBackEnd/pkg/model/card"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"os"
	"time"
)

type ICardRepository interface {
	Persist(cardToPersist *card.Card) (*card.Card, error)
	FindById(userId string, id *primitive.ObjectID, private bool) (result *card.Card, err error)
	FindByDeckId(userId, deckId string, private bool) (cardReturn []*card.Card, err error)
	Update(id *primitive.ObjectID, userId string, cardToSave *card.Card) (*card.Card, error)
	Count(userId string) (count int64, err error)
	Delete(userId string, id *primitive.ObjectID) (result *card.Card, err error)
}

type CardRepository struct {
	client         *mongo.Client
	database       string
	cardCollection string
}

func NewCardRepository(mongoClient *mongo.Client) CardRepository {
	return CardRepository{
		client:         mongoClient,
		database:       os.Getenv("MONGODB_DATABASE"),
		cardCollection: os.Getenv("MONGODB_CARD_COLLECTION"),
	}
}

const (
	errorString = "error trying to decode result"
)

func (a CardRepository) Persist(card *card.Card) (*card.Card, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := a.client.Database(a.database).Collection(a.cardCollection).InsertOne(ctx, card)
	if err != nil {
		log.Logger.Errorw("Persist has failed", "error", err)
		return nil, errors.Wrap(err, "error trying to persist entity")
	}
	id := result.InsertedID.(primitive.ObjectID)
	card.Id = &id

	return card, nil
}

func (a CardRepository) FindById(userId string, id *primitive.ObjectID, private bool) (cardReturn *card.Card, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	privateResult := bson.M{}
	if !private {
		privateResult = bson.M{"isPrivate": false}
	}

	col := a.client.Database(a.database).Collection(a.cardCollection)

	query := bson.M{"_id": id, "$or": []interface{}{
		privateResult,
		bson.M{"userId": userId},
	}}
	result := col.FindOne(ctx, query)
	if result.Err() != nil {
		log.Logger.Warn("Find cardReturn by id has failed", "Error", result.Err())
		return nil, errors.Wrap(result.Err(), "error trying to find cardReturn by id")
	}

	err = result.Decode(&cardReturn)
	if err != nil {
		log.Logger.Errorw(errorString, "error", err)
		return nil, errors.Wrap(err, errorString)
	}

	return cardReturn, nil
}
func (a CardRepository) FindByDeckId(userId, deckId string, private bool) (cardReturn []*card.Card, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	privateResult := bson.M{}
	if !private {
		privateResult = bson.M{"isPrivate": false}
	}

	col := a.client.Database(a.database).Collection(a.cardCollection)

	query := bson.M{"deckId": deckId, "$or": []interface{}{
		privateResult,
		bson.M{"userId": userId},
	}}
	result, err := col.Find(ctx, query)
	if err != nil {
		log.Logger.Errorw("Find has failed", errorString, err.Error())
		return nil, errors.Wrap(err, "error trying to find card")
	}

	cardReturn = make([]*card.Card, 0)
	for result.Next(ctx) {
		var cardElement *card.Card
		err := result.Decode(&cardElement)
		if err != nil {
			log.Logger.Errorw("Parser Card has failed", "error", err.Error())
			return nil, errors.Wrap(err, "error trying to parse Card")
		}
		cardReturn = append(cardReturn, cardElement)
	}

	err = result.Close(ctx)
	if err != nil {
		log.Logger.Errorw("Error closing context...", "Error", err.Error())
		return nil, errors.Wrap(err, "error trying to find card")
	}
	return cardReturn, nil
}

func (a CardRepository) FindByUserId(userId string) (cardResult []*card.Card, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := a.client.Database(a.database).Collection(a.cardCollection)
	result, err := col.Find(ctx, bson.M{"userId": userId})

	if err != nil {
		log.Logger.Errorw("Find has failed", errorString, err.Error())
		return nil, errors.Wrap(err, "error trying to find card")
	}

	cardResult = make([]*card.Card, 0)
	for result.Next(ctx) {
		var cardElement *card.Card
		err := result.Decode(&cardElement)
		if err != nil {
			log.Logger.Errorw("Parser Card has failed", "error", err.Error())
			return nil, errors.Wrap(err, "error trying to parse Card")
		}
		cardResult = append(cardResult, cardElement)
	}

	err = result.Close(ctx)
	if err != nil {
		log.Logger.Errorw("Error closing context...", "Error", err.Error())
		return nil, errors.Wrap(err, "error trying to find card")
	}

	return
}

func (a CardRepository) Count(userId string) (count int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err = a.client.Database(a.database).Collection(a.cardCollection).CountDocuments(ctx, bson.M{"userId": userId})
	if err != nil {
		log.Logger.Errorw("Count Cards has failed", "error", err.Error())
		return 0, errors.Wrap(err, "error trying to count Cards")
	}

	return count, nil
}

func (a CardRepository) Delete(userId string, id *primitive.ObjectID) (result *card.Card, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := bson.M{"_id": id,
		"userId": userId,
	}
	singleResult := a.client.Database(a.database).
		Collection(a.cardCollection).FindOneAndDelete(ctx, query)

	if singleResult.Err() != nil {
		log.Logger.Errorw("Delete has failed", "Error", err)
		return nil, errors.Wrap(singleResult.Err(), "error trying to delete entity")
	}

	err = singleResult.Decode(&result)
	if err != nil {
		log.Logger.Errorw("Parser Card has failed", "Error", err)
		return nil, errors.Wrap(err, "Parser Card has failed")
	}

	log.Logger.Debug("Card deleted: ", result)
	return result, err
}

func (a CardRepository) Update(id *primitive.ObjectID, userId string, cardToSave *card.Card) (*card.Card, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filterQuery := bson.M{"_id": bson.M{"$eq": id}, "userId": bson.M{"$eq": userId}}
	updateResult, err := a.client.Database(a.database).Collection(a.cardCollection).ReplaceOne(ctx, filterQuery, cardToSave)
	if err != nil {
		log.Logger.Errorw("Update has failed", "error", err)
		return nil, errors.Wrap(err, "error trying to update entity")
	}

	if updateResult != nil && updateResult.MatchedCount == 0 {
		return nil, nil
	}
	cardToSave.Id = id
	return cardToSave, nil
}
