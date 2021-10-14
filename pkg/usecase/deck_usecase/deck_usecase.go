package deck_usecase

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/errors"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/deck"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/repository/deck_repository"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/imdario/mergo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IDeckUseCase interface {
	Create(userId string, deck *deck.Deck) (result *deck.Deck, err error)
	FindByUserIdAndPublic(userId string) (result []*deck.Deck, count int64, err error)
	FindById(userId, id string) (result *deck.Deck, err error)
	FindByUserId(userId string) (result []*deck.Deck, count int64, err error)
	Delete(id, userId string) (result *deck.Deck, err error)
	Update(id, userId string, isPartial bool, deck *deck.Deck) (*deck.Deck, error)
}
type DeckUseCase struct {
	validator *validator.Validate
	repo      deck_repository.IDeckRepository
}

func NewDeckUseCase(deckRepository deck_repository.IDeckRepository, validator *validator.Validate) DeckUseCase {
	return DeckUseCase{
		validator: validator,
		repo:      deckRepository,
	}
}

func (uc DeckUseCase) Create(userId string, deck *deck.Deck) (result *deck.Deck, err error) {
	deck.UserId = userId
	deck.LastUpdate = time.Now()
	deck.CardsCount = 0

	err = uc.validator.Struct(deck)
	if err != nil {
		deckBytes, _ := json.Marshal(deck)
		log.Logger.Errorf("Error to validate input:\n %v;\n error: %v", string(deckBytes), err.Error())
		return nil, &errors.InvalidPayload{Err: err}
	}

	result, err = uc.repo.Persist(deck)
	if err != nil {
		log.Logger.Errorw("Deck creation error", "Error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}

	return result, nil
}

func (uc DeckUseCase) FindById(userId, id string) (result *deck.Deck, err error) {
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

func (uc DeckUseCase) FindByUserId(userId string) (result []*deck.Deck, count int64, err error) {
	result, err = uc.repo.FindByUserId(userId)
	if err != nil {
		log.Logger.Errorw("deck not found", "Error", err.Error())
		return nil, 0, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}

	count, err = uc.repo.Count(userId)
	if err != nil {
		log.Logger.Errorw("Count decks error", "Error", err.Error())
		return nil, 0, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}

	return result, count, nil
}

func (uc DeckUseCase) FindByUserIdAndPublic(userId string) (result []*deck.Deck, count int64, err error) {
	result, err = uc.repo.FindByUserIdAndPublic(userId)
	if err != nil {
		log.Logger.Errorw("deck not found", "Error", err.Error())
		return nil, 0, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}

	count, err = uc.repo.Count(userId)
	if err != nil {
		log.Logger.Errorw("Count decks error", "Error", err.Error())
		return nil, 0, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}
	return result, count, nil
}

func (uc DeckUseCase) parseToObjectID(id string) (objID primitive.ObjectID, err error) {
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

func (uc DeckUseCase) Delete(id, userId string) (result *deck.Deck, err error) {
	objectID, err := uc.parseToObjectID(id)
	if err != nil {
		return nil, err
	}

	result, err = uc.repo.Delete(userId, &objectID)
	if err != nil {
		log.Logger.Errorw("Remove deck error", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}
	return
}

func (uc DeckUseCase) Update(id, userId string, isPartial bool, deck *deck.Deck) (*deck.Deck, error) {
	objectID, err := uc.parseToObjectID(id)
	if err != nil {
		return nil, err
	}

	deck.UserId = userId
	deck.LastUpdate = time.Now()

	savedDeck, err := uc.repo.FindById(userId, &objectID, true)
	if err != nil {
		log.Logger.Errorw("deck not found", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}

	if isPartial {
		if err := mergo.Merge(deck, savedDeck); err != nil {
			log.Logger.Errorf("Error to merge entities")
			return nil, &errors.BadRequest{Err: err}
		}
	}

	err = uc.validator.Struct(deck)
	if err != nil {
		deckBytes, _ := json.Marshal(deck)
		log.Logger.Errorf("Error to validate input:\n %v;\n error: %v", string(deckBytes), err.Error())
		return nil, &errors.InvalidPayload{Err: err}
	}

	result, err := uc.repo.Update(&objectID, userId, deck)
	if err != nil {
		log.Logger.Errorw("update error", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}
	if deck.IsPrivate == !savedDeck.IsPrivate {
		//todo cardsUsecase.updatePrivacy(deckId,boleano)
	}

	if result == nil {
		return nil, errors.WrapWithMessage(errors.ErrNotFound, fmt.Sprintf("id %s not found", id))
	}

	return result, nil
}

//func (uc DeckUseCase) AddDeckToPlaylist(id, userId string, card *card.Card) (*deck.Deck, error) {
//
//	savedDeck, err := uc.FindById(userId, id)
//	if err != nil {
//		log.Logger.Errorw("playlist not found", "error", err.Error())
//		return nil, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
//	}
//
//	savedDeck.Cards = append(savedDeck.Cards, card)
//	savedDeck.CardsCount += 1
//
//	result, err := uc.Update(id, userId, true, savedDeck)
//	if err != nil {
//		log.Logger.Errorw("update error", "error", err.Error())
//		return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
//	}
//
//	if result == nil {
//		return nil, errors.WrapWithMessage(errors.ErrNotFound, fmt.Sprintf("id %s not found", id))
//	}
//
//	return result, nil
//}
