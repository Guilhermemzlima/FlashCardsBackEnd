package usecase

import (
	"FlashCardsBackEnd/pkg/usecase/playlist_usecase"
	"github.com/google/wire"
)

var playlistSet = wire.NewSet(
	playlist_usecase.NewPlaylistUseCase,
	wire.Bind(new(playlist_usecase.IPlaylistUseCase), new(playlist_usecase.PlaylistUseCase)))

var Set = wire.NewSet(
	playlistSet,
)
