package deck_usecase

import (
	"encoding/json"
	"fmt"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/errors"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/Enums"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/deck"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/filter"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/review"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/repository/deck_repository"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/repository/review_repository"
	"github.com/go-playground/validator"
	"github.com/imdario/mergo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IDeckUseCase interface {
	Create(userId string, deck *deck.Deck) (result *deck.Deck, err error)
	FindById(userId, id string) (result *deck.Deck, err error)
	FindByUserId(filter, userId string, pagination *filter.Pagination, private bool, orderBy string, order string) (result []*deck.Deck, count int64, err error)
	Delete(id, userId string) (result *deck.Deck, err error)
	Update(id, userId string, isPartial bool, deck *deck.Deck) (*deck.Deck, error)
	FindBySearch(filter, userId string, pagination *filter.Pagination) (result []map[string]interface{}, count int64, err error)
	FindRecent(userId string, pagination *filter.Pagination) (result []*deck.Deck, count int64, err error)
	AddCounterPlayDeck(id, userId string, savedDeck *deck.Deck) (deckReturn *deck.Deck, err error)
}
type DeckUseCase struct {
	validator  *validator.Validate
	repo       deck_repository.IDeckRepository
	reviewRepo review_repository.IReviewRepository
}

func NewDeckUseCase(deckRepository deck_repository.IDeckRepository, reviewRepo review_repository.IReviewRepository, validator *validator.Validate) DeckUseCase {
	return DeckUseCase{
		validator:  validator,
		repo:       deckRepository,
		reviewRepo: reviewRepo,
	}
}

func (uc DeckUseCase) Create(userId string, deck *deck.Deck) (result *deck.Deck, err error) {
	deck.UserId = userId
	deck.LastUpdate = time.Now()
	deck.CreatedAt = time.Now()
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

func (uc DeckUseCase) FindByIdArray(userId string, ids []string) (result []*deck.Deck, err error) {
	var idsPrimitive []*primitive.ObjectID
	for _, element := range ids {
		objectID, err := uc.parseToObjectID(element)
		if err != nil {
			return nil, err
		}
		idsPrimitive = append(idsPrimitive, &objectID)
	}

	result, err = uc.repo.FindByIdArray(userId, idsPrimitive, false)
	if err != nil {
		log.Logger.Errorw("deck not found", "Error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}
	return result, nil
}

func (uc DeckUseCase) FindByUserId(filter, userId string, pagination *filter.Pagination, private bool, orderBy string, order string) (result []*deck.Deck, count int64, err error) {
	if orderBy == "" {
		orderBy = "createdAt"
	}
	if order == "" {
		order = "DESC"
	}
	result, err = uc.repo.FindByUserId(filter, userId, pagination, private, orderBy, order)
	if err != nil {
		log.Logger.Errorw("deck not found", "Error", err.Error())
		return nil, 0, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}

	count, err = uc.repo.Count(userId, private)
	if err != nil {
		log.Logger.Errorw("Count decks error", "Error", err.Error())
		return nil, 0, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}

	return result, count, nil
}

func (uc DeckUseCase) FindRecent(userId string, pagination *filter.Pagination) (result []*deck.Deck, count int64, err error) {
	reviews, err := uc.reviewRepo.FindRecent(Enums.Deck, pagination, userId)

	var list []*review.Review
	for _, item := range reviews {
		if contains(list, item) == false {
			list = append(list, item)
		}
	}

	for _, s := range list {
		objectID, err := uc.parseToObjectID(s.OriginId)
		if err != nil {
			return nil, 0, err
		}
		deckFound, err := uc.repo.FindById(userId, &objectID, false)
		if err != nil {
			log.Logger.Errorw("deck not found", "Error", err.Error())
			return nil, 0, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
		}
		result = append(result, deckFound)
	}

	count, err = uc.repo.Count(userId, false)
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

func (uc DeckUseCase) FindBySearch(filter, userId string, pagination *filter.Pagination) (result []map[string]interface{}, count int64, err error) {
	result, count, err = uc.repo.FindByFilters(filter, userId, pagination)
	if err != nil {
		log.Logger.Errorw("deck not found", "Error", err.Error())
		return nil, 0, errors.WrapWithMessage(errors.ErrNotFound, err.Error())
	}

	return result, count, nil
}

func (uc DeckUseCase) AddCounterPlayDeck(id, userId string, savedDeck *deck.Deck) (deckReturn *deck.Deck, err error) {
	savedDeck.PlayCount = savedDeck.PlayCount + 1
	deckReturn, err = uc.Update(id, savedDeck.UserId, false, savedDeck)
	if err != nil {
		log.Logger.Errorw("update error", "error", err.Error())
		return nil, errors.WrapWithMessage(errors.ErrInternalServer, err.Error())
	}
	return
}

func contains(s []*review.Review, e *review.Review) bool {
	for _, a := range s {
		if a.OriginId == e.OriginId {
			return true
		}
	}
	return false
}
