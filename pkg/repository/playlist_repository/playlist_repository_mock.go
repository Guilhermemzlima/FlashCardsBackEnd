package playlist_repository

import (
	"FlashCardsBackEnd/pkg/model/playlist"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockPlaylistRepository struct {
	mock.Mock
}

func (uc *MockPlaylistRepository) Persist(playlistToPersist *playlist.Playlist) (*playlist.Playlist, error) {
	args := uc.Called(playlistToPersist)
	if playlistReturn := args.Get(0); playlistReturn != nil {
		return playlistReturn.(*playlist.Playlist), args.Error(1)
	}
	return nil, args.Error(1)
}

func (uc *MockPlaylistRepository) FindByUserId(userId string) (result []*playlist.Playlist, err error) {
	args := uc.Called(userId)
	return args.Get(0).([]*playlist.Playlist), args.Error(1)
}
func (uc *MockPlaylistRepository) FindByUserIdAndPublic(userId string) (result []*playlist.Playlist, err error) {
	args := uc.Called(userId)
	return args.Get(0).([]*playlist.Playlist), args.Error(1)
}

func (uc *MockPlaylistRepository) FindById(userId string, id *primitive.ObjectID, private bool) (*playlist.Playlist, error) {
	args := uc.Called(userId, id, private)
	return args.Get(0).(*playlist.Playlist), args.Error(1)
}

func (uc *MockPlaylistRepository) Update(id *primitive.ObjectID, userId string, playlistToSave *playlist.Playlist) (*playlist.Playlist, error) {
	args := uc.Called(id, userId, playlistToSave)
	if playlistReturn := args.Get(0); playlistReturn != nil {
		return playlistReturn.(*playlist.Playlist), args.Error(1)
	}
	return nil, args.Error(1)
}

func (uc *MockPlaylistRepository) Delete(userId string, id *primitive.ObjectID) (*playlist.Playlist, error) {
	args := uc.Called(id, userId)
	return args.Get(0).(*playlist.Playlist), args.Error(1)
}

func (db *MockPlaylistRepository) Count(userId string) (count int64, err error) {
	args := db.Called(userId)
	return args.Get(0).(int64), args.Error(1)
}
