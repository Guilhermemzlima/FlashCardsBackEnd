package repository

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/repository/card_repository"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/repository/deck_repository"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/repository/playlist_repository"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/repository/review_repository"
	"github.com/google/wire"
)

var playlistRepositorySet = wire.NewSet(
	playlist_repository.NewPlaylistRepository,
	wire.Bind(new(playlist_repository.IPlaylistRepository), new(playlist_repository.PlaylistRepository)))

var deckRepositorySet = wire.NewSet(
	deck_repository.NewDeckRepository,
	wire.Bind(new(deck_repository.IDeckRepository), new(deck_repository.DeckRepository)))

var cardRepositorySet = wire.NewSet(
	card_repository.NewCardRepository,
	wire.Bind(new(card_repository.ICardRepository), new(card_repository.CardRepository)))

var reviewRepositorySet = wire.NewSet(
	review_repository.NewReviewRepository,
	wire.Bind(new(review_repository.IReviewRepository), new(review_repository.ReviewRepository)))

var Set = wire.NewSet(
	playlistRepositorySet,
	deckRepositorySet,
	cardRepositorySet,
	reviewRepositorySet,
)
