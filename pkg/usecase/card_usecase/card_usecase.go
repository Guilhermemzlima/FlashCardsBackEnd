package card_usecase

import (
	"FlashCardsBackEnd/internal/config/log"
	"FlashCardsBackEnd/internal/errors"
	"FlashCardsBackEnd/pkg/model/card"
	"FlashCardsBackEnd/pkg/repository/card_repository"
	"FlashCardsBackEnd/pkg/usecase/deck_usecase"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/imdario/mergo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ICardUseCase interface {
	Create(userId, deckId string, card *card.Card) (result *card.Card, err error)
	FindByDeckId(userId, id string) (result []*card.Card, err error)
	FindById(userId, id string) (result *card.Card, err error)
	Update(id, userId string, isPartial bool, card *card.Card) (*card.Card, error)
	Delete(id, userId string) (result *card.Card, err error)
}
type CardUseCase struct {
	validator   *validator.Validate
	repo        card_repository.ICardRepository
	deckUseCase deck_usecase.DeckUseCase
}

func NewCardUseCase(cardRepository card_repository.ICardRepository, deckUseCase deck_usecase.DeckUseCase, validator *validator.Validate) CardUseCase {
	return CardUseCase{
		validator:   validator,
		repo:        cardRepository,
		deckUseCase: deckUseCase,
	}
}

func (uc CardUseCase) Create(userId, deckId string, card *card.Card) (result *card.Card, err error) {
	card.UserId = userId
	card.DeckId = deckId
	card.LastUpdate = time.Now()

	err = uc.validator.Struct(card)
	if err != nil {
		cardBytes, _ := json.Marshal(card)
		log.Logger.Errorf("Error to validate input:\n %v;\n error: %v", string(cardBytes), err.Error())
		return nil, &errors.InvalidPayload{Err: err}
	}

	result, err = uc.repo.Persist(card)
	if err != nil {
		log.Logger.Errorw("Card creation error", "Error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}

	//_, err = uc.deckUseCase.AddDeckToPlaylist(deckId, userId, result)
	//if err != nil {
	//	return nil, err
	//}

	return result, nil
}

func (uc CardUseCase) FindByDeckId(userId, id string) (result []*card.Card, err error) {
	result, err = uc.repo.FindByDeckId(userId, id, false)
	if err != nil {
		log.Logger.Errorw("card not found", "Error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}
	return result, nil
}

func (uc CardUseCase) FindById(userId, id string) (result *card.Card, err error) {
	objectID, err := uc.parseToObjectID(id)
	if err != nil {
		return nil, err
	}

	result, err = uc.repo.FindById(userId, &objectID, false)
	if err != nil {
		log.Logger.Errorw("card not found", "Error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}
	return result, nil
}

//func (uc CardUseCase) FindByUserId(userId string) (result []*card.Card, count int64, err error) {
//	result, err = uc.repo.FindByUserId(userId)
//	if err != nil {
//		log.Logger.Errorw("card not found", "Error", err.Error())
//		return nil, 0, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
//	}
//
//	count, err = uc.repo.Count(userId)
//	if err != nil {
//		log.Logger.Errorw("Count cards error", "Error", err.Error())
//		return nil, 0, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
//	}
//
//	return result, count, nil
//}
//
//func (uc CardUseCase) FindByUserIdAndPublic(userId string) (result []*card.Card, count int64, err error) {
//	result, err = uc.repo.FindByUserIdAndPublic(userId)
//	if err != nil {
//		log.Logger.Errorw("card not found", "Error", err.Error())
//		return nil, 0, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
//	}
//
//	count, err = uc.repo.Count(userId)
//	if err != nil {
//		log.Logger.Errorw("Count cards error", "Error", err.Error())
//		return nil, 0, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
//	}
//	return result, count, nil
//}
//
func (uc CardUseCase) parseToObjectID(id string) (objID primitive.ObjectID, err error) {
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

func (uc CardUseCase) Delete(id, userId string) (result *card.Card, err error) {
	objectID, err := uc.parseToObjectID(id)
	if err != nil {
		return nil, err
	}

	result, err = uc.repo.Delete(userId, &objectID)
	if err != nil {
		log.Logger.Errorw("Remove card error", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}
	return
}

func (uc CardUseCase) Update(id, userId string, isPartial bool, card *card.Card) (*card.Card, error) {
	objectID, err := uc.parseToObjectID(id)
	if err != nil {
		return nil, err
	}

	card.UserId = userId
	card.LastUpdate = time.Now()

	savedCard, err := uc.repo.FindById(userId, &objectID, true)
	if err != nil {
		log.Logger.Errorw("card not found", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}

	if isPartial {
		if err := mergo.Merge(card, savedCard); err != nil {
			log.Logger.Errorf("Error to merge entities")
			return nil, &errors.BadRequest{Err: err}
		}
	}

	err = uc.validator.Struct(card)
	if err != nil {
		cardBytes, _ := json.Marshal(card)
		log.Logger.Errorf("Error to validate input:\n %v;\n error: %v", string(cardBytes), err.Error())
		return nil, &errors.InvalidPayload{Err: err}
	}

	result, err := uc.repo.Update(&objectID, userId, card)
	if err != nil {
		log.Logger.Errorw("update error", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}

	if result == nil {
		return nil, errors.WrapWithMessage(errors.ErrNotFound, fmt.Sprintf("id %s not found", id))
	}

	return result, nil
}
