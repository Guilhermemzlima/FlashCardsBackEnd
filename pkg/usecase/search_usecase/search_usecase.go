package search_usecase

import (
	"FlashCardsBackEnd/internal/config/log"
	"FlashCardsBackEnd/internal/errors"
	"FlashCardsBackEnd/pkg/usecase/deck_usecase"
	"FlashCardsBackEnd/pkg/usecase/playlist_usecase"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISearchUseCase interface {
	Search(filter, userId string) (result []map[string]interface{}, err error)
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

func (uc SearchUseCase) Search(filter, userId string) (result []map[string]interface{}, err error) {
	playlistSearched, _, err := uc.playlistUseCase.FindBySearch(filter, userId)
	if err != nil {
		return nil, err
	}
	result = append(result, playlistSearched...)

	decksSearched, _, err := uc.deckUseCase.FindBySearch(filter, userId)
	if err != nil {
		return nil, err
	}
	result = append(result, decksSearched...)

	if err != nil {
		return nil, err
	}
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
