package repository

import (
	"FlashCardsBackEnd/pkg/repository/playlist_repository"
	"github.com/google/wire"
)

var playlistRepositorySet = wire.NewSet(
	playlist_repository.NewPlaylistRepository,
	wire.Bind(new(playlist_repository.IPlaylistRepository), new(playlist_repository.PlaylistRepository)))

var Set = wire.NewSet(playlistRepositorySet)
