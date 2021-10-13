package playlist_repository
//
//import (
//	"FlashCardsBackEnd/internal/config/log"
//	"FlashCardsBackEnd/internal/infra/mongodb"
//	"FlashCardsBackEnd/pkg/model/playlist"
//	utils "FlashCardsBackEnd/pkg/test"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//var playlistRepository PlaylistRepository
//
//var (
//	previousPlaylist *playlist.Playlist
//)
//
//func init() {
//	log.SetupLogger()
//	utils.InitMongoContainer()
//
//	mongodbClient, err := mongodb.NewMongoDbClient()
//	if err != nil {
//		panic(err)
//	}
//	playlistRepository = NewPlaylistRepository(mongodbClient)
//	previousPlaylist = utils.BuildPlaylist()
//}
//
//func TestCreate(t *testing.T) {
//	// given
//	i, err := playlistRepository.Count(previousPlaylist.UserId)
//	assert.NoError(t, err)
//	// when
//	playlistSaved, err := playlistRepository.Persist(previousPlaylist)
//	// then
//	assert.NoError(t, err, "Persist with error")
//	assert.NotNil(t, playlistSaved)
//
//	c, err := playlistRepository.Count(previousPlaylist.UserId)
//	assert.NoError(t, err)
//	assert.NotEqual(t, i, c)
//}
//
//func TestFindById(t *testing.T) {
//	//given
//	clearMongo()
//	playlistSaved, err := playlistRepository.Persist(previousPlaylist)
//	assert.NoError(t, err)
//
//	// when
//	playlistFound, err := playlistRepository.FindById(playlistSaved.UserId, playlistSaved.Id, false)
//
//	// then
//	assert.NoError(t, err)
//	assert.NotNil(t, playlistFound)
//}
//
//func TestFindByUserId(t *testing.T) {
//	//given
//	clearMongo()
//	playlistSaved, err := playlistRepository.Persist(previousPlaylist)
//	assert.NoError(t, err)
//
//	// when
//	playlistsFound, err := playlistRepository.FindByUserId(playlistSaved.UserId)
//
//	// then
//	assert.NoError(t, err)
//	assert.NotNil(t, playlistsFound)
//}
//
//func TestFindByUserIdAndPublic(t *testing.T) {
//	//given
//	clearMongo()
//	playlistSaved, err := playlistRepository.Persist(previousPlaylist)
//	assert.NoError(t, err)
//
//	// when
//	playlistsFound, err := playlistRepository.FindByUserIdAndPublic(playlistSaved.UserId)
//
//	// then
//	assert.NoError(t, err)
//	assert.NotNil(t, playlistsFound)
//}
//
//func TestDelete(t *testing.T) {
//	//given
//	clearMongo()
//	savedPlaylist, _ := playlistRepository.Persist(previousPlaylist)
//
//	// when
//	_, err := playlistRepository.Delete(previousPlaylist.UserId, previousPlaylist.Id)
//
//	playlistFound, _ := playlistRepository.FindById(savedPlaylist.UserId, savedPlaylist.Id, true)
//
//	// then
//	assert.Nil(t, err)
//	assert.Nil(t, playlistFound)
//}
//
//func TestCount(t *testing.T) {
//	//given
//	clearMongo()
//	_, _ = playlistRepository.Persist(previousPlaylist)
//	_, _ = playlistRepository.Persist(previousPlaylist)
//
//	// when
//	count, err := playlistRepository.Count(
//		previousPlaylist.UserId,
//	)
//
//	// then
//	assert.Nil(t, err)
//	assert.NotNil(t, count)
//}
//
//func TestUpdate(t *testing.T) {
//	//given
//	clearMongo()
//	savedPlaylist, _ := playlistRepository.Persist(previousPlaylist)
//	previousPlaylist.Name = "New Name"
//	// when
//	_, err := playlistRepository.Update(savedPlaylist.Id, previousPlaylist.UserId, previousPlaylist)
//
//	// then
//	playlistFound, _ := playlistRepository.FindById(previousPlaylist.UserId, savedPlaylist.Id,true)
//
//	assert.Nil(t, err)
//	assert.Equal(t, "New Name", playlistFound.Name)
//}
//
//
//func clearMongo() {
//	log.Logger.Errorw("cleaning mongo test")
//	result, _ := playlistRepository.FindByUserIdAndPublic(previousPlaylist.UserId)
//
//	for _, playlistFound := range result {
//		id := playlistFound.Id
//		_, err := playlistRepository.Delete(previousPlaylist.UserId, id)
//
//		if err != nil {
//			panic("Fail to clear mongo")
//		}
//	}
//}
