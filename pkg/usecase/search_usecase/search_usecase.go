package search_usecase

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/errors"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/filter"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/usecase/deck_usecase"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/usecase/playlist_usecase"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISearchUseCase interface {
	Search(filter, userId string, pagination *filter.Pagination) (result []map[string]interface{}, count int64, err error)
}

type SearchUseCase struct {
	validator       *validator.Validate
	playlistUseCase playlist_usecase.PlaylistUseCase
	deckUseCase     deck_usecase.DeckUseCase
}

func NewSearchUseCase(playlistUseCase playlist_usecase.PlaylistUseCase, deckUseCase deck_usecase.DeckUseCase, validator *validator.Validate) SearchUseCase {
	return SearchUseCase{
		validator:       validator,
		playlistUseCase: playlistUseCase,
		deckUseCase:     deckUseCase,
	}
}

func (uc SearchUseCase) Search(filter, userId string, pagination *filter.Pagination) (result []map[string]interface{}, count int64, err error) {
	playlistSearched, countPlaylist, err := uc.playlistUseCase.FindBySearch(filter, userId, pagination)
	if err != nil {
		return nil, 0, err
	}
	result = append(result, playlistSearched...)

	decksSearched, countDeck, err := uc.deckUseCase.FindBySearch(filter, userId, pagination)
	if err != nil {
		return nil, 0, err
	}
	result = append(result, decksSearched...)

	if err != nil {
		return nil, 0, err
	}
	count = countPlaylist + countDeck
	return
}

func (uc SearchUseCase) parseToObjectID(id string) (objID primitive.ObjectID, err error) {
	if id == "" {
		err := errors.WrapWithMessage(errors.ErrInvalidPayload, "id is required")
		log.Logger.Errorw("id is required")
		return objID, err
	}

	objID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Logger.Errorw("invalid Id", "Error", err.Error())
		return objID, errors.WrapWithMessage(errors.ErrInvalidPayload, err.Error())
	}

	return
}
