package handlers

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/handler/card_handler"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/handler/deck_handler"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/handler/playlist_handler"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/handler/review_handler"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/handler/search_handler"
	"github.com/google/wire"
)

var ApplicationHandlersSet = wire.NewSet(
	playlist_handler.NewPlaylistHandler,
	deck_handler.NewDeckHandler,
	card_handler.NewCardHandler,
	review_handler.NewReviewHandler,
	search_handler.NewSearchHandler,
)
