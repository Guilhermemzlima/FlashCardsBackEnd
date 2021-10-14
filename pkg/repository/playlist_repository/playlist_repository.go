package playlist_repository

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/playlist"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"time"
)

type IPlaylistRepository interface {
	Persist(playlistToPersist *playlist.Playlist) (*playlist.Playlist, error)
	FindById(userId string, id *primitive.ObjectID, private bool) (result *playlist.Playlist, err error)
	FindByUserIdAndPublic(userId string) (result []*playlist.Playlist, err error)
	FindByUserId(userId string) (result []*playlist.Playlist, err error)
	Count(userId string) (count int64, err error)
	Delete(userId string, id *primitive.ObjectID) (result *playlist.Playlist, err error)
	Update(id *primitive.ObjectID, userId string, playlistToSave *playlist.Playlist) (*playlist.Playlist, error)
	FindFilter(filter, userId string) (playlistResult []map[string]interface{}, err error)
}

type PlaylistRepository struct {
	client             *mongo.Client
	database           string
	playlistCollection string
}

func NewPlaylistRepository(mongoClient *mongo.Client) PlaylistRepository {
	return PlaylistRepository{
		client:             mongoClient,
		database:           os.Getenv("MONGODB_DATABASE"),
		playlistCollection: os.Getenv("MONGODB_PLAYLIST_COLLECTION"),
	}
}

const (
	errorString = "error trying to decode result"
)

func (a PlaylistRepository) Persist(playlist *playlist.Playlist) (*playlist.Playlist, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := a.client.Database(a.database).Collection(a.playlistCollection).InsertOne(ctx, playlist)
	if err != nil {
		log.Logger.Errorw("Persist has failed", "error", err)
		return nil, errors.Wrap(err, "error trying to persist entity")
	}
	id := result.InsertedID.(primitive.ObjectID)
	playlist.Id = &id

	return playlist, nil
}

func (a PlaylistRepository) FindById(userId string, id *primitive.ObjectID, private bool) (playlist *playlist.Playlist, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	privateResult := bson.M{}
	if !private {
		privateResult = bson.M{"isPrivate": false}
	}

	col := a.client.Database(a.database).Collection(a.playlistCollection)

	query := bson.M{"_id": id, "$or": []interface{}{
		privateResult,
		bson.M{"userId": userId},
	}}
	result := col.FindOne(ctx, query)
	if result.Err() != nil {
		log.Logger.Warn("Find playlist by id has failed", "Error", result.Err())
		return nil, errors.Wrap(result.Err(), "error trying to find playlist by id")
	}

	err = result.Decode(&playlist)
	if err != nil {
		log.Logger.Errorw(errorString, "error", err)
		return nil, errors.Wrap(err, errorString)
	}

	return playlist, nil
}

func (a PlaylistRepository) FindByUserId(userId string) (playlistResult []*playlist.Playlist, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := a.client.Database(a.database).Collection(a.playlistCollection)
	result, err := col.Find(ctx, bson.M{"userId": userId})

	if err != nil {
		log.Logger.Errorw("Find has failed", errorString, err.Error())
		return nil, errors.Wrap(err, "error trying to find playlists")
	}

	playlistResult = make([]*playlist.Playlist, 0)
	for result.Next(ctx) {
		var playlistElement *playlist.Playlist
		err := result.Decode(&playlistElement)
		if err != nil {
			log.Logger.Errorw("Parser Playlist has failed", "error", err.Error())
			return nil, errors.Wrap(err, "error trying to parse Playlist")
		}
		playlistResult = append(playlistResult, playlistElement)
	}

	err = result.Close(ctx)
	if err != nil {
		log.Logger.Errorw("Error closing context...", "Error", err.Error())
		return nil, errors.Wrap(err, "error trying to find playlists")
	}

	return
}

func (a PlaylistRepository) FindByUserIdAndPublic(userId string) (playlistResult []*playlist.Playlist, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := a.client.Database(a.database).Collection(a.playlistCollection)
	query := bson.M{"$or": []interface{}{
		bson.M{"isPrivate": false},
		bson.M{"userId": userId},
	}}
	result, err := col.Find(ctx, query)

	if err != nil {
		log.Logger.Errorw("Find has failed", errorString, err.Error())
		return nil, errors.Wrap(err, "error trying to find playlists")
	}

	playlistResult = make([]*playlist.Playlist, 0)
	for result.Next(ctx) {
		var playlistElement *playlist.Playlist
		err := result.Decode(&playlistElement)
		if err != nil {
			log.Logger.Errorw("Parser Playlist has failed", "error", err.Error())
			return nil, errors.Wrap(err, "error trying to parse Playlist")
		}
		playlistResult = append(playlistResult, playlistElement)
	}
	return
}
func (a PlaylistRepository) Count(userId string) (count int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err = a.client.Database(a.database).Collection(a.playlistCollection).CountDocuments(ctx, bson.M{"userId": userId})
	if err != nil {
		log.Logger.Errorw("Count playlists has failed", "error", err.Error())
		return 0, errors.Wrap(err, "error trying to count playlists")
	}

	return count, nil
}

func (a PlaylistRepository) Delete(userId string, id *primitive.ObjectID) (result *playlist.Playlist, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := bson.M{"_id": id,
		"userId": userId,
	}
	singleResult := a.client.Database(a.database).
		Collection(a.playlistCollection).FindOneAndDelete(ctx, query)

	if singleResult.Err() != nil {
		log.Logger.Errorw("Delete has failed", "Error", err)
		return nil, errors.Wrap(singleResult.Err(), "error trying to delete entity")
	}

	err = singleResult.Decode(&result)
	if err != nil {
		log.Logger.Errorw("Parser Playlist has failed", "Error", err)
		return nil, errors.Wrap(err, "Parser Playlist has failed")
	}

	log.Logger.Debug("playlist deleted: ", result)
	return result, err
}

func (a PlaylistRepository) Update(id *primitive.ObjectID, userId string, playlist *playlist.Playlist) (*playlist.Playlist, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filterQuery := bson.M{"_id": bson.M{"$eq": id}, "userId": bson.M{"$eq": userId}}
	updateResult, err := a.client.Database(a.database).Collection(a.playlistCollection).ReplaceOne(ctx, filterQuery, playlist)
	if err != nil {
		log.Logger.Errorw("Update has failed", "error", err)
		return nil, errors.Wrap(err, "error trying to update entity")
	}

	if updateResult != nil && updateResult.MatchedCount == 0 {
		return nil, nil
	}
	playlist.Id = id
	return playlist, nil
}

//func (a PlaylistRepository) FindFilter(filter, userId string) (playlistResult []*playlist.Playlist, err error) {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	result, err := a.client.Database(a.database).Collection(a.playlistCollection).Find(ctx, bson.M{"name": bson.M{"$regex": primitive.Regex{
//		Pattern: "/.*" + filter + ".*/",
//		Options: "i",
//	}},
//		"$or": []interface{}{
//			bson.M{"isPrivate": false},
//			bson.M{"userId": userId},
//		},
//	})
//
//	if err != nil {
//		log.Logger.Errorw("Find has failed", errorString, err.Error())
//		return nil, errors.Wrap(err, "error trying to find playlists")
//	}
//
//	playlistResult = make([]*playlist.Playlist, 0)
//	for result.Next(ctx) {
//		var playlistElement *playlist.Playlist
//		err := result.Decode(&playlistElement)
//		if err != nil {
//			log.Logger.Errorw("Parser Playlist has failed", "error", err.Error())
//			return nil, errors.Wrap(err, "error trying to parse Playlist")
//		}
//		playlistResult = append(playlistResult, playlistElement)
//	}
//	return
//}

func (a PlaylistRepository) FindFilter(filter, userId string) (playlistResult []map[string]interface{}, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := a.client.Database(a.database).Collection(a.playlistCollection)

	query := bson.M{"name": bson.M{"$regex": primitive.Regex{
		Pattern: ".*" + filter + ".*",
		Options: "i",
	}},
		"$or": []interface{}{
			bson.M{"isPrivate": false},
			bson.M{"userId": userId},
		}}

	result, err := col.Find(ctx, query)

	if err != nil {
		log.Logger.Errorw("Find has failed", errorString, err.Error())
		return nil, errors.Wrap(err, "error trying to find playlists")
	}

	playlistResult = make([]map[string]interface{}, 0)
	for result.Next(ctx) {
		var playlistElement map[string]interface{}
		err := result.Decode(&playlistElement)
		if err != nil {
			log.Logger.Errorw("Parser Playlist has failed", "error", err.Error())
			return nil, errors.Wrap(err, "error trying to parse Playlist")
		}
		playlistResult = append(playlistResult, playlistElement)
	}
	return
}
