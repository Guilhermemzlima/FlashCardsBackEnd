package deck_handler

import (
	"encoding/json"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/errors"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/render"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/deck"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/filter"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/usecase/deck_usecase"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	headerUserId = "userId"
	pathVarID    = "id"
)

type DeckHandler struct {
	deckUseCase deck_usecase.IDeckUseCase
}

func NewDeckHandler(service deck_usecase.IDeckUseCase) DeckHandler {
	return DeckHandler{deckUseCase: service}
}

func (handler *DeckHandler) Post(w http.ResponseWriter, r *http.Request) {
	requestBody, err := handler.extractBody(r)
	if err != nil {
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(errors.ErrInvalidPayload))
		return
	}

	userID := r.Header.Get(headerUserId)
	result, err := handler.deckUseCase.Create(userID, requestBody)
	if err != nil {
		log.Logger.Errorw("Failed to create deck", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Deck has created successfully")
	render.Response(w, result, http.StatusCreated)
}

func (handler *DeckHandler) FindById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[pathVarID]
	userID := r.Header.Get(headerUserId)

	result, err := handler.deckUseCase.FindById(userID, id)
	if err != nil {
		log.Logger.Errorw("Failed to find deck", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Deck has find successfully")
	render.Response(w, result, http.StatusOK)
}

func (handler *DeckHandler) FindByUserIdAndPublic(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(headerUserId)

	pagination := handler.buildPagination(r)

	result, count, err := handler.deckUseCase.FindByUserId(userID, pagination, false)
	if err != nil {
		log.Logger.Errorw("Failed to find deck", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	w.Header().Add("X-Total", strconv.FormatInt(count, 10))
	render.Response(w, result, http.StatusOK)
}

func (handler *DeckHandler) FindRecent(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(headerUserId)

	pagination := handler.buildPagination(r)

	result, count, err := handler.deckUseCase.FindRecent(userID, pagination)
	if err != nil {
		log.Logger.Errorw("Failed to find deck", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	w.Header().Add("X-Total", strconv.FormatInt(count, 10))
	render.Response(w, result, http.StatusOK)
}

func (handler *DeckHandler) FindByUserId(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(headerUserId)
	pagination := handler.buildPagination(r)
	result, count, err := handler.deckUseCase.FindByUserId(userID, pagination, true)
	if err != nil {
		log.Logger.Errorw("Failed to find deck", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	w.Header().Add("X-Total", strconv.FormatInt(count, 10))
	render.Response(w, result, http.StatusOK)
}

func (handler *DeckHandler) Patch(w http.ResponseWriter, r *http.Request) {
	requestBody, err := handler.extractBody(r)
	if err != nil {
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(errors.ErrInvalidPayload))
		return
	}

	id := mux.Vars(r)[pathVarID]
	userID := r.Header.Get(headerUserId)
	result, err := handler.deckUseCase.Update(id, userID, true, requestBody)
	if err != nil {
		log.Logger.Errorw("Failed to update deck", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Deck has updated successfully")
	render.Response(w, result, http.StatusCreated)
}

func (handler *DeckHandler) Delete(w http.ResponseWriter, r *http.Request) {
	customerID := r.Header.Get(headerUserId)

	id := mux.Vars(r)[pathVarID]
	_, err := handler.deckUseCase.Delete(id, customerID)
	if err != nil {
		log.Logger.Errorw("Failed to remove deck", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Deck has removed successfully")
	render.Response(w, nil, http.StatusNoContent)
}

func (handler *DeckHandler) extractBody(r *http.Request) (*deck.Deck, error) {
	var deckJSON deck.Deck

	if err := json.NewDecoder(r.Body).Decode(&deckJSON); err != nil {
		log.Logger.Errorw("Error trying to parse payload", "error", err)
		return nil, err
	}

	return &deckJSON, nil
}

func (handler *DeckHandler) buildPagination(r *http.Request) (pagination *filter.Pagination) {
	reqQuery := r.URL.Query()
	limit, _ := strconv.ParseInt(reqQuery.Get("limit"), 10, 64)
	offset, _ := strconv.ParseInt(reqQuery.Get("offset"), 10, 64)
	return filter.NewPagination(limit, offset)
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
