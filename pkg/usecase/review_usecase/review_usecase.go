package review_usecase

import (
	"fmt"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/errors"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/Enums"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/card"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/review"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/repository/review_repository"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/usecase/card_usecase"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/usecase/deck_usecase"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/usecase/playlist_usecase"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IReviewUseCase interface {
	ReviewPlaylists(id, userId string) (map[string]interface{}, error)
	FindById(userId, id string) (result *review.Review, err error)
	AddCardResult(sessionId, userId string, card *card.Card, isRight bool) (*review.Review, error)
	ReviewDecks(id, userId string) (map[string]interface{}, error)
}

type ReviewUseCase struct {
	validator       *validator.Validate
	repo            review_repository.ReviewRepository
	playlistUseCase playlist_usecase.PlaylistUseCase
	deckUseCase     deck_usecase.DeckUseCase
	cardUseCase     card_usecase.CardUseCase
}

func NewReviewUseCase(playlistUseCase playlist_usecase.PlaylistUseCase, repo review_repository.ReviewRepository, cardUseCase card_usecase.CardUseCase, deckUseCase deck_usecase.DeckUseCase, validator *validator.Validate) ReviewUseCase {
	return ReviewUseCase{
		validator:       validator,
		repo:            repo,
		playlistUseCase: playlistUseCase,
		deckUseCase:     deckUseCase,
		cardUseCase:     cardUseCase,
	}
}

func (uc ReviewUseCase) ReviewPlaylists(id, userId string) (map[string]interface{}, error) {
	playlistToReview, err := uc.playlistUseCase.FindById(userId, id)
	if err != nil {
		return nil, err
	}
	var list []*card.Card

	for _, element := range playlistToReview.Decks {
		cards, err := uc.cardUseCase.FindByDeckId(userId, element.Id)
		if err != nil {
			log.Logger.Errorw("error to find cards by deckId", "error", err.Error())
			return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
		}
		list = append(list, cards...)
	}

	review, err := uc.repo.Persist(&review.Review{
		OriginType:    Enums.Playlist,
		OriginId:      playlistToReview.Id.Hex(),
		UserId:        userId,
		Hists:         nil,
		HistsCount:    0,
		Mistakes:      nil,
		MistakesCount: 0,
		LastUpdate:    time.Time{},
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"cards":   list,
		"session": review,
	}, nil
}

func (uc ReviewUseCase) ReviewDecks(id, userId string) (map[string]interface{}, error) {
	deck, err := uc.deckUseCase.FindById(userId, id)
	if err != nil {
		return nil, err
	}
	var list []*card.Card

	cards, err := uc.cardUseCase.FindByDeckId(userId, deck.Id.Hex())
	if err != nil {
		log.Logger.Errorw("error to find cards by deckId", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}
	list = append(list, cards...)

	review, err := uc.repo.Persist(&review.Review{
		OriginType:    Enums.Deck,
		OriginId:      deck.Id.Hex(),
		UserId:        userId,
		Hists:         nil,
		HistsCount:    0,
		Mistakes:      nil,
		MistakesCount: 0,
		LastUpdate:    time.Time{},
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"cards":   list,
		"session": review,
	}, nil
}

func (uc ReviewUseCase) AddCardResult(sessionId, userId string, card *card.Card, isRight bool) (*review.Review, error) {
	savedReview, err := uc.FindById(userId, sessionId)
	if err != nil {
		log.Logger.Errorw("playlist not found", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}
	if isRight {
		savedReview.Hists = append(savedReview.Hists, card)
		savedReview.HistsCount += 1
	}
	if isRight == false {
		savedReview.Mistakes = append(savedReview.Hists, card)
		savedReview.MistakesCount += 1
	}
	id, err := uc.parseToObjectID(sessionId)
	if err != nil {
		return nil, err
	}

	result, err := uc.repo.Update(&id, userId, savedReview)
	if err != nil {
		log.Logger.Errorw("update error", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}

	if result == nil {
		return nil, errors.WrapWithMessage(errors.ErrNotFound, fmt.Sprintf("id %s not found", id))
	}

	return result, nil
}

func (uc ReviewUseCase) FindById(userId, id string) (result *review.Review, err error) {
	objectID, err := uc.parseToObjectID(id)
	if err != nil {
		return nil, err
	}

	result, err = uc.repo.FindById(userId, &objectID, false)
	if err != nil {
		log.Logger.Errorw("deck not found", "Error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}
	return result, nil
}

func (uc ReviewUseCase) parseToObjectID(id string) (objID primitive.ObjectID, err error) {
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
