package review_repository

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/infra/mongodb"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/review"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"os"
	"time"
)

type IReviewRepository interface {
	Persist(reviewToPersist *review.Review) (*review.Review, error)
	FindById(userId string, id *primitive.ObjectID, private bool) (reviewReturn *review.Review, err error)
	//FindByDeckId(userId, deckId string, private bool) (reviewReturn []*review.Review, err error)
	//Update(id *primitive.ObjectID, userId string, reviewToSave *review.Review) (*review.Review, error)
	//Count(userId string) (count int64, err error)
	//Delete(userId string, id *primitive.ObjectID) (result *review.Review, err error)
}

type ReviewRepository struct {
	client           *mongo.Client
	database         string
	reviewCollection string
}

func NewReviewRepository(mongoClient *mongo.Client) ReviewRepository {
	return ReviewRepository{
		client:           mongoClient,
		database:         os.Getenv("MONGODB_DATABASE"),
		reviewCollection: os.Getenv("MONGODB_REVIEW_COLLECTION"),
	}
}

const (
	errorString = "error trying to decode result"
)

func (a ReviewRepository) Persist(review *review.Review) (*review.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := a.client.Database(a.database).Collection(a.reviewCollection).InsertOne(ctx, review)
	if err != nil {
		log.Logger.Errorw("Persist has failed", "error", err)
		return nil, errors.Wrap(err, "error trying to persist entity")
	}
	id := result.InsertedID.(primitive.ObjectID)
	review.Id = &id

	return review, nil
}

func (a ReviewRepository) FindById(userId string, id *primitive.ObjectID, private bool) (reviewReturn *review.Review, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	privateResult := bson.M{}
	if !private {
		privateResult = bson.M{"isPrivate": false}
	}

	col := a.client.Database(a.database).Collection(a.reviewCollection)

	query := bson.M{"_id": id, "$or": []interface{}{
		privateResult,
		bson.M{"userId": userId},
	}}
	result := col.FindOne(ctx, query)
	if result.Err() != nil {
		log.Logger.Warn("Find reviewReturn by id has failed", "Error", result.Err())
		return nil, errors.Wrap(result.Err(), "error trying to find reviewReturn by id")
	}

	err = result.Decode(&reviewReturn)
	if err != nil {
		log.Logger.Errorw(errorString, "error", err)
		return nil, errors.Wrap(err, errorString)
	}

	return reviewReturn, nil
}

func (a ReviewRepository) FindRecent(originType, userId string) (reviewResult []*review.Review, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := a.client.Database(a.database).Collection(a.reviewCollection)
	findOptions := options.Find()
	findOptions.SetLimit(10).SetSkip(0)
	findOptions.SetSort(bson.D{{"lastUpdate", mongodb.DESC}})

	query := bson.M{"originType": originType, "$or": []interface{}{
		bson.M{"isPrivate": false},
		bson.M{"userId": userId},
	}}

	result, err := col.Find(ctx, query)
	if result.Err() != nil {
		log.Logger.Warn("Find reviewReturn has failed", "Error", result.Err())
		return nil, errors.Wrap(result.Err(), "error trying to find reviewReturn by id")
	}

	reviewResult = make([]*review.Review, 0)
	for result.Next(ctx) {
		var playlistElement *review.Review
		err := result.Decode(&playlistElement)
		if err != nil {
			log.Logger.Errorw("Parser Playlist has failed", "error", err.Error())
			return nil, errors.Wrap(err, "error trying to parse Playlist")
		}
		reviewResult = append(reviewResult, playlistElement)
	}
	return
}

//func (a ReviewRepository) FindByDeckId(userId, deckId string, private bool) (reviewReturn []*review.Review, err error) {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	privateResult := bson.M{}
//	if !private {
//		privateResult = bson.M{"isPrivate": false}
//	}
//
//	col := a.client.Database(a.database).Collection(a.reviewCollection)
//
//	query := bson.M{"deckId": deckId, "$or": []interface{}{
//		privateResult,
//		bson.M{"userId": userId},
//	}}
//	result, err := col.Find(ctx, query)
//	if err != nil {
//		log.Logger.Errorw("Find has failed", errorString, err.Error())
//		return nil, errors.Wrap(err, "error trying to find review")
//	}
//
//	reviewReturn = make([]*review.Review, 0)
//	for result.Next(ctx) {
//		var reviewElement *review.Review
//		err := result.Decode(&reviewElement)
//		if err != nil {
//			log.Logger.Errorw("Parser Review has failed", "error", err.Error())
//			return nil, errors.Wrap(err, "error trying to parse Review")
//		}
//		reviewReturn = append(reviewReturn, reviewElement)
//	}
//
//	err = result.Close(ctx)
//	if err != nil {
//		log.Logger.Errorw("Error closing context...", "Error", err.Error())
//		return nil, errors.Wrap(err, "error trying to find review")
//	}
//	return reviewReturn, nil
//}
//
//func (a ReviewRepository) FindByUserId(userId string) (reviewResult []*review.Review, err error) {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	col := a.client.Database(a.database).Collection(a.reviewCollection)
//	result, err := col.Find(ctx, bson.M{"userId": userId})
//
//	if err != nil {
//		log.Logger.Errorw("Find has failed", errorString, err.Error())
//		return nil, errors.Wrap(err, "error trying to find review")
//	}
//
//	reviewResult = make([]*review.Review, 0)
//	for result.Next(ctx) {
//		var reviewElement *review.Review
//		err := result.Decode(&reviewElement)
//		if err != nil {
//			log.Logger.Errorw("Parser Review has failed", "error", err.Error())
//			return nil, errors.Wrap(err, "error trying to parse Review")
//		}
//		reviewResult = append(reviewResult, reviewElement)
//	}
//
//	err = result.Close(ctx)
//	if err != nil {
//		log.Logger.Errorw("Error closing context...", "Error", err.Error())
//		return nil, errors.Wrap(err, "error trying to find review")
//	}
//
//	return
//}
//
//func (a ReviewRepository) Count(userId string) (count int64, err error) {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	count, err = a.client.Database(a.database).Collection(a.reviewCollection).CountDocuments(ctx, bson.M{"userId": userId})
//	if err != nil {
//		log.Logger.Errorw("Count Reviews has failed", "error", err.Error())
//		return 0, errors.Wrap(err, "error trying to count Reviews")
//	}
//
//	return count, nil
//}
//
//func (a ReviewRepository) Delete(userId string, id *primitive.ObjectID) (result *review.Review, err error) {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	query := bson.M{"_id": id,
//		"userId": userId,
//	}
//	singleResult := a.client.Database(a.database).
//		Collection(a.reviewCollection).FindOneAndDelete(ctx, query)
//
//	if singleResult.Err() != nil {
//		log.Logger.Errorw("Delete has failed", "Error", err)
//		return nil, errors.Wrap(singleResult.Err(), "error trying to delete entity")
//	}
//
//	err = singleResult.Decode(&result)
//	if err != nil {
//		log.Logger.Errorw("Parser Review has failed", "Error", err)
//		return nil, errors.Wrap(err, "Parser Review has failed")
//	}
//
//	log.Logger.Debug("Review deleted: ", result)
//	return result, err
//}
//

func (a ReviewRepository) Update(id *primitive.ObjectID, userId string, reviewToSave *review.Review) (*review.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filterQuery := bson.M{"_id": bson.M{"$eq": id}, "userId": bson.M{"$eq": userId}}
	updateResult, err := a.client.Database(a.database).Collection(a.reviewCollection).ReplaceOne(ctx, filterQuery, reviewToSave)
	if err != nil {
		log.Logger.Errorw("Update has failed", "error", err)
		return nil, errors.Wrap(err, "error trying to update entity")
	}

	if updateResult != nil && updateResult.MatchedCount == 0 {
		return nil, nil
	}
	reviewToSave.Id = id
	return reviewToSave, nil
}
