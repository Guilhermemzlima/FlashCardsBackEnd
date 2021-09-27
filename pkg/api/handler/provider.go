package handlers

import (
	"FlashCardsBackEnd/pkg/api/handler/playlist_handler"
	"github.com/google/wire"
)

var ApplicationHandlersSet = wire.NewSet(
	playlist_handler.NewPlaylistHandler,
)
