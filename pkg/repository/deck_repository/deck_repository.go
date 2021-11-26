package deck_repository

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/infra/mongodb"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/deck"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/filter"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"os"
	"time"
)

type IDeckRepository interface {
	Persist(deckToPersist *deck.Deck) (*deck.Deck, error)
	FindById(userId string, id *primitive.ObjectID, private bool) (result *deck.Deck, err error)
	FindByUserId(filter,userId string, pagination *filter.Pagination, private bool) (deckResult []*deck.Deck, err error)
	Count(userId string, private bool) (count int64, err error)
	Delete(userId string, id *primitive.ObjectID) (result *deck.Deck, err error)
	Update(id *primitive.ObjectID, userId string, deckToSave *deck.Deck) (*deck.Deck, error)
	FindByFilters(filter, userId string, pagination *filter.Pagination) (deckResult []map[string]interface{}, count int64, err error)
	FindByIdArray(userId string, ids []*primitive.ObjectID, private bool) (deckResult []*deck.Deck, err error)
}

type DeckRepository struct {
	client         *mongo.Client
	database       string
	deckCollection string
}

func NewDeckRepository(mongoClient *mongo.Client) DeckRepository {
	return DeckRepository{
		client:         mongoClient,
		database:       os.Getenv("MONGODB_DATABASE"),
		deckCollection: os.Getenv("MONGODB_DECK_COLLECTION"),
	}
}

const (
	errorString = "error trying to decode result"
)

func (a DeckRepository) Persist(deck *deck.Deck) (*deck.Deck, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := a.client.Database(a.database).Collection(a.deckCollection).InsertOne(ctx, deck)
	if err != nil {
		log.Logger.Errorw("Persist has failed", "error", err)
		return nil, errors.Wrap(err, "error trying to persist entity")
	}
	id := result.InsertedID.(primitive.ObjectID)
	deck.Id = &id

	return deck, nil
}

func (a DeckRepository) FindById(userId string, id *primitive.ObjectID, private bool) (deckReturn *deck.Deck, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	privateResult := bson.M{}
	if !private {
		privateResult = bson.M{"isPrivate": false}
	}

	col := a.client.Database(a.database).Collection(a.deckCollection)

	query := bson.M{"_id": id, "$or": []interface{}{
		privateResult,
		bson.M{"userId": userId},
	}}
	result := col.FindOne(ctx, query)
	if result.Err() != nil {
		log.Logger.Warn("Find deckReturn by id has failed", "Error", result.Err())
		return nil, errors.Wrap(result.Err(), "error trying to find deckReturn by id")
	}

	err = result.Decode(&deckReturn)
	if err != nil {
		log.Logger.Errorw(errorString, "error", err)
		return nil, errors.Wrap(err, errorString)
	}

	return deckReturn, nil
}

func (a DeckRepository) FindByIdArray(userId string, ids []*primitive.ObjectID, private bool) (deckResult []*deck.Deck, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	privateResult := bson.M{}
	if !private {
		privateResult = bson.M{"isPrivate": false}
	}

	col := a.client.Database(a.database).Collection(a.deckCollection)

	query := bson.M{"_id": bson.M{"$in": ids}, "$or": []interface{}{
		privateResult,
		bson.M{"userId": userId},
	}}
	result, err := col.Find(ctx, query)
	if result.Err() != nil {
		log.Logger.Warn("Find deckReturn by id has failed", "Error", result.Err())
		return nil, errors.Wrap(result.Err(), "error trying to find deckReturn by id")
	}
	deckResult = make([]*deck.Deck, 0)
	for result.Next(ctx) {
		var deckElement *deck.Deck
		err := result.Decode(&deckElement)
		if err != nil {
			log.Logger.Errorw("Parser Deck has failed", "error", err.Error())
			return nil, errors.Wrap(err, "error trying to parse Deck")
		}
		deckResult = append(deckResult, deckElement)
	}

	err = result.Close(ctx)
	if err != nil {
		log.Logger.Errorw("Error closing context...", "Error", err.Error())
		return nil, errors.Wrap(err, "error trying to find decks")
	}

	return
}

func (a DeckRepository) FindByUserId(filter, userId string, pagination *filter.Pagination, private bool) (deckResult []*deck.Deck, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := a.client.Database(a.database).Collection(a.deckCollection)
	findOptions := options.Find()
	findOptions.SetLimit(pagination.Limit).SetSkip(pagination.Skip)
	findOptions.SetSort(bson.D{{"createdAt", mongodb.ASC}})

	privateResult := bson.M{}
	if !private {
		privateResult = bson.M{"isPrivate": false}
	}

	query := bson.M{"name": bson.M{"$regex": primitive.Regex{
		Pattern: ".*" + filter + ".*",
		Options: "i",
	}}, "$or": []interface{}{
		privateResult,
		bson.M{"userId": userId},
	}}

	result, err := col.Find(ctx, query, findOptions)
	if err != nil {
		log.Logger.Errorw("Find has failed", errorString, err.Error())
		return nil, errors.Wrap(err, "error trying to find deck")
	}

	deckResult = make([]*deck.Deck, 0)
	for result.Next(ctx) {
		var deckElement *deck.Deck
		err := result.Decode(&deckElement)
		if err != nil {
			log.Logger.Errorw("Parser Deck has failed", "error", err.Error())
			return nil, errors.Wrap(err, "error trying to parse Deck")
		}
		deckResult = append(deckResult, deckElement)
	}

	err = result.Close(ctx)
	if err != nil {
		log.Logger.Errorw("Error closing context...", "Error", err.Error())
		return nil, errors.Wrap(err, "error trying to find decks")
	}

	return
}

func (a DeckRepository) Count(userId string, private bool) (count int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	privateResult := bson.M{}
	if !private {
		privateResult = bson.M{"isPrivate": false}
	}

	query := bson.M{"$or": []interface{}{
		privateResult,
		bson.M{"userId": userId},
	}}

	count, err = a.client.Database(a.database).Collection(a.deckCollection).CountDocuments(ctx, query)
	if err != nil {
		log.Logger.Errorw("Count Decks has failed", "error", err.Error())
		return 0, errors.Wrap(err, "error trying to count Decks")
	}

	return count, nil
}

func (a DeckRepository) Delete(userId string, id *primitive.ObjectID) (result *deck.Deck, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := bson.M{"_id": id,
		"userId": userId,
	}
	singleResult := a.client.Database(a.database).
		Collection(a.deckCollection).FindOneAndDelete(ctx, query)

	if singleResult.Err() != nil {
		log.Logger.Errorw("Delete has failed", "Error", err)
		return nil, errors.Wrap(singleResult.Err(), "error trying to delete entity")
	}

	err = singleResult.Decode(&result)
	if err != nil {
		log.Logger.Errorw("Parser Deck has failed", "Error", err)
		return nil, errors.Wrap(err, "Parser Deck has failed")
	}

	log.Logger.Debug("Deck deleted: ", result)
	return result, err
}

func (a DeckRepository) Update(id *primitive.ObjectID, userId string, deckToSave *deck.Deck) (*deck.Deck, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filterQuery := bson.M{"_id": bson.M{"$eq": id}, "userId": bson.M{"$eq": userId}}
	updateResult, err := a.client.Database(a.database).Collection(a.deckCollection).ReplaceOne(ctx, filterQuery, deckToSave)
	if err != nil {
		log.Logger.Errorw("Update has failed", "error", err)
		return nil, errors.Wrap(err, "error trying to update entity")
	}

	if updateResult != nil && updateResult.MatchedCount == 0 {
		return nil, nil
	}
	deckToSave.Id = id
	return deckToSave, nil
}

func (a DeckRepository) FindByFilters(filter, userId string, pagination *filter.Pagination) (deckResult []map[string]interface{}, count int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetLimit(pagination.Limit / 2).SetSkip(pagination.Skip)
	findOptions.SetSort(bson.D{{"lastUpdate", mongodb.DESC}})

	col := a.client.Database(a.database).Collection(a.deckCollection)
	query := bson.M{"name": bson.M{"$regex": primitive.Regex{
		Pattern: ".*" + filter + ".*",
		Options: "i",
	}}, "$or": []interface{}{
		bson.M{"isPrivate": false},
		bson.M{"userId": userId},
	}}

	result, err := col.Find(ctx, query, findOptions)
	if err != nil {
		log.Logger.Errorw("Find has failed", errorString, err.Error())
		return nil, 0, errors.Wrap(err, "error trying to find decks")
	}

	count, err = col.CountDocuments(ctx, query)
	if err != nil {
		log.Logger.Errorw("Count Decks has failed", "error", err.Error())
		return nil, 0, errors.Wrap(err, "error trying to count Decks")
	}

	deckResult = make([]map[string]interface{}, 0)
	for result.Next(ctx) {
		var deckElement map[string]interface{}
		err := result.Decode(&deckElement)
		if err != nil {
			log.Logger.Errorw("Parser Deck has failed", "error", err.Error())
			return nil, 0, errors.Wrap(err, "error trying to parse Deck")
		}
		deckResult = append(deckResult, deckElement)
	}
	return
}
