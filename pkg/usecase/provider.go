package usecase

import (
	"FlashCardsBackEnd/pkg/usecase/card_usecase"
	"FlashCardsBackEnd/pkg/usecase/deck_usecase"
	"FlashCardsBackEnd/pkg/usecase/playlist_usecase"
	"FlashCardsBackEnd/pkg/usecase/review_usecase"
	"github.com/google/wire"
)

var playlistSet = wire.NewSet(
	playlist_usecase.NewPlaylistUseCase,
	wire.Bind(new(playlist_usecase.IPlaylistUseCase), new(playlist_usecase.PlaylistUseCase)))

var deckSet = wire.NewSet(
	deck_usecase.NewDeckUseCase,
	wire.Bind(new(deck_usecase.IDeckUseCase), new(deck_usecase.DeckUseCase)))

var cardSet = wire.NewSet(
	card_usecase.NewCardUseCase,
	wire.Bind(new(card_usecase.ICardUseCase), new(card_usecase.CardUseCase)))

var reviewSet = wire.NewSet(
	review_usecase.NewReviewUseCase,
	wire.Bind(new(review_usecase.IReviewUseCase), new(review_usecase.ReviewUseCase)))

var Set = wire.NewSet(
	playlistSet,
	deckSet,
	cardSet,
	reviewSet,
)
