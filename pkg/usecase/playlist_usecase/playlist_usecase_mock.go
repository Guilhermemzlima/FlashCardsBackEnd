package playlist_usecase

import (
	"FlashCardsBackEnd/pkg/model/playlist"
	"github.com/stretchr/testify/mock"
)

type MockPlaylistUseCase struct {
	mock.Mock
}

func (uc *MockPlaylistUseCase) FindByUserId(customerId string) (playlists []*playlist.Playlist, count int64, err error) {
	args := uc.Called(customerId)
	return args.Get(0).([]*playlist.Playlist), args.Get(1).(int64), args.Error(2)
}
func (uc *MockPlaylistUseCase) FindByUserIdAndPublic(customerId string) (playlists []*playlist.Playlist, count int64, err error) {
	args := uc.Called(customerId)
	return args.Get(0).([]*playlist.Playlist), args.Get(1).(int64), args.Error(2)
}

func (uc *MockPlaylistUseCase) FindByID(userId, id string) (*playlist.Playlist, error) {
	args := uc.Called(userId, id)
	return args.Get(0).(*playlist.Playlist), args.Error(1)
}

func (uc *MockPlaylistUseCase) Create(userId string, playlists *playlist.Playlist) (*playlist.Playlist, error) {
	args := uc.Called(userId, playlists)
	return args.Get(0).(*playlist.Playlist), args.Error(1)
}

func (uc *MockPlaylistUseCase) Update(id, userId string, isPartial bool, playlists *playlist.Playlist) (*playlist.Playlist, error) {
	args := uc.Called(id, userId, isPartial, playlists)
	return args.Get(0).(*playlist.Playlist), args.Error(1)
}

func (uc *MockPlaylistUseCase) Delete(id, userId string) (*playlist.Playlist, error) {
	args := uc.Called(id, userId)
	return args.Get(0).(*playlist.Playlist), args.Error(1)
}
