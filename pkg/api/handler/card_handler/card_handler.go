package card_handler

import (
	"encoding/json"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/errors"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/render"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/card"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/usecase/card_usecase"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	headerUserId = "userId"
	pathVarID    = "id"
)

type CardHandler struct {
	cardUseCase card_usecase.ICardUseCase
}

func NewCardHandler(service card_usecase.ICardUseCase) CardHandler {
	return CardHandler{cardUseCase: service}
}

func (handler *CardHandler) Post(w http.ResponseWriter, r *http.Request) {
	requestBody, err := handler.extractBody(r)
	if err != nil {
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(errors.ErrInvalidPayload))
		return
	}
	id := mux.Vars(r)[pathVarID]

	userID := r.Header.Get(headerUserId)
	result, err := handler.cardUseCase.Create(userID, id, requestBody)
	if err != nil {
		log.Logger.Errorw("Failed to create card", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Card has created successfully")
	render.Response(w, result, http.StatusCreated)
}

func (handler *CardHandler) FindByDeckId(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[pathVarID]
	userID := r.Header.Get(headerUserId)

	result, err := handler.cardUseCase.FindByDeckId(userID, id)
	if err != nil {
		log.Logger.Errorw("Failed to find card", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Card has find successfully")
	render.Response(w, result, http.StatusOK)
}

func (handler *CardHandler) FindById(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(headerUserId)
	id := mux.Vars(r)[pathVarID]
	result, err := handler.cardUseCase.FindById(userID, id)
	if err != nil {
		log.Logger.Errorw("Failed to find card", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	render.Response(w, result, http.StatusOK)
}

func (handler *CardHandler) Patch(w http.ResponseWriter, r *http.Request) {
	requestBody, err := handler.extractBody(r)
	if err != nil {
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(errors.ErrInvalidPayload))
		return
	}

	id := mux.Vars(r)[pathVarID]
	userID := r.Header.Get(headerUserId)
	result, err := handler.cardUseCase.Update(id, userID, true, requestBody)
	if err != nil {
		log.Logger.Errorw("Failed to update card", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Card has updated successfully")
	render.Response(w, result, http.StatusOK)
}

func (handler *CardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	customerID := r.Header.Get(headerUserId)

	id := mux.Vars(r)[pathVarID]
	_, err := handler.cardUseCase.Delete(id, customerID)
	if err != nil {
		log.Logger.Errorw("Failed to remove card", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Card has removed successfully")
	render.Response(w, nil, http.StatusNoContent)
}

func (handler *CardHandler) extractBody(r *http.Request) (*card.Card, error) {
	var cardJSON card.Card

	if err := json.NewDecoder(r.Body).Decode(&cardJSON); err != nil {
		log.Logger.Errorw("Error trying to parse payload", "error", err)
		return nil, err
	}

	return &cardJSON, nil
}

func GenerateHTTPErrorStatusCode(err error) int {
	switch errors.Cause(err).(type) {
	case *errors.NotFound:
		return http.StatusNotFound
	case *errors.InvalidPayload:
		return http.StatusPreconditionFailed
	case *errors.BadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
