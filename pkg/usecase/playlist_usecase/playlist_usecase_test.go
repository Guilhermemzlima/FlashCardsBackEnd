package playlist_usecase

//
//import (
//	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/validator"
//	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/playlist"
//	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/repository/playlist_repository"
//	utils "github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/test"
//	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/usecase/deck_usecase"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"testing"
//)
//
//func init() {
//	previousPlaylist = utils.BuildPlaylist()
//}
//
//var previousPlaylist *playlist.Playlist
//var resultPlaylist *playlist.Playlist
//var list []*playlist.Playlist
//var count = int64(13)
//var countResult int64
//var err error
//var listResult []*playlist.Playlist
//
//var playlistRepositoryMock *playlist_repository.MockPlaylistRepository
//
//var playlistUseCase *PlaylistUseCase
//
//func TestCreate(t *testing.T) {
//	// given
//	givenNewMock()
//	givenPersist()
//	givenNewUseCase()
//
//	// when
//	whenCreate()
//
//	// then
//	thanAssertPlaylist(t)
//	thanAssertOnePlaylist(t)
//}
//
//func TestFindByUserId(t *testing.T) {
//	// given
//	list = append(list, previousPlaylist)
//	givenNewMock()
//	givenFindByUserId()
//	givenCount()
//	givenNewUseCase()
//	// when
//	whenFindByUserId()
//	// then
//	thenAssertLists(t)
//	thenAssertFindByUserIdAndCountCalls(t)
//}
//
//func givenNewMock() {
//	playlistRepositoryMock = new(playlist_repository.MockPlaylistRepository)
//	playlistRepositoryMock = new(playlist_repository.MockPlaylistRepository)
//}
//
//func givenFindById() {
//	playlistRepositoryMock.On(utils.FindByIdMethodName, mock.Anything, mock.Anything, mock.Anything).Return(previousPlaylist, nil)
//}
//
//func givenFindByUserId() {
//	playlistRepositoryMock.On(utils.FindByUserIdMethodName, mock.Anything).Return(list, nil)
//}
//
//func givenCount() {
//	playlistRepositoryMock.On(utils.CountMethodName, mock.Anything).Return(count, nil)
//}
//
//func givenNewUseCase() {
//	useCase := NewPlaylistUseCase(playlistRepositoryMock, validator.NewValidate())
//	playlistUseCase = &useCase
//}
//
//func givenPersist() {
//	playlistRepositoryMock.On(utils.PersistMethodName, mock.Anything).Return(previousPlaylist, nil)
//}
//
//func whenFindByUserId() {
//	listResult, countResult, err = playlistUseCase.FindByUserId(previousPlaylist.UserId)
//}
//func whenFindById() {
//	resultPlaylist, err = playlistUseCase.FindById(previousPlaylist.UserId, "id")
//}
//
//func thenAssertLists(t *testing.T) {
//	assert.Nil(t, err)
//	assert.Equal(t, count, countResult)
//	assert.Equal(t, list[0], listResult[0])
//}
//
//func thenAssertFindByUserIdAndCountCalls(t *testing.T) {
//	playlistRepositoryMock.AssertNumberOfCalls(t, utils.FindByUserIdMethodName, 1)
//	playlistRepositoryMock.AssertNumberOfCalls(t, utils.CountMethodName, 1)
//}
//
//func thanAssertPlaylist(t *testing.T) {
//	assert.NoError(t, err)
//	assert.NotNil(t, resultPlaylist)
//	assert.Equal(t, resultPlaylist.UserId, previousPlaylist.UserId)
//}
//
//func whenCreate() {
//	resultPlaylist, err = playlistUseCase.Create(previousPlaylist.UserId, previousPlaylist)
//}
//
//func thanAssertOnePlaylist(t *testing.T) {
//	playlistRepositoryMock.AssertNumberOfCalls(t, utils.PersistMethodName, 1)
//}
