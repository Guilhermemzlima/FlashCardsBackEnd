package review_handler

import (
	"FlashCardsBackEnd/internal/config/log"
	"FlashCardsBackEnd/internal/errors"
	"FlashCardsBackEnd/pkg/api/render"
	"FlashCardsBackEnd/pkg/model/card"
	"FlashCardsBackEnd/pkg/usecase/review_usecase"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	headerUserId = "userId"
	pathVarID    = "id"
)

type ReviewHandler struct {
	reviewUseCase review_usecase.IReviewUseCase
}

func NewReviewHandler(service review_usecase.IReviewUseCase) ReviewHandler {
	return ReviewHandler{reviewUseCase: service}
}

func (handler *ReviewHandler) ReviewPlaylist(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[pathVarID]
	userID := r.Header.Get(headerUserId)

	result, err := handler.reviewUseCase.ReviewPlaylists(id, userID)
	if err != nil {
		log.Logger.Errorw("Failed to review playlist", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Playlist has find successfully")
	render.Response(w, result, http.StatusOK)
}

func (handler *ReviewHandler) ReviewDeck(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[pathVarID]
	userID := r.Header.Get(headerUserId)

	result, err := handler.reviewUseCase.ReviewDecks(id, userID)
	if err != nil {
		log.Logger.Errorw("Failed to review playlist", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Deck has find successfully")
	render.Response(w, result, http.StatusOK)
}

func (handler *ReviewHandler) CardResultRight(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[pathVarID]
	requestBody, err := handler.extractBody(r)
	userID := r.Header.Get(headerUserId)

	result, err := handler.reviewUseCase.AddCardResult(id, userID, requestBody, true)
	if err != nil {
		log.Logger.Errorw("Failed to review playlist", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Playlist has find successfully")
	render.Response(w, result, http.StatusOK)
}
func (handler *ReviewHandler) CardResultWrong(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[pathVarID]
	requestBody, err := handler.extractBody(r)
	userID := r.Header.Get(headerUserId)

	result, err := handler.reviewUseCase.AddCardResult(id, userID, requestBody, false)
	if err != nil {
		log.Logger.Errorw("Failed to review playlist", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Playlist has find successfully")
	render.Response(w, result, http.StatusOK)
}

func (handler *ReviewHandler) FindById(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(headerUserId)
	id := mux.Vars(r)[pathVarID]
	result, err := handler.reviewUseCase.FindById(userID, id)
	if err != nil {
		log.Logger.Errorw("Failed to find card", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	render.Response(w, result, http.StatusOK)
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

func (handler *ReviewHandler) extractBody(r *http.Request) (*card.Card, error) {
	var cardJSON card.Card

	if err := json.NewDecoder(r.Body).Decode(&cardJSON); err != nil {
		log.Logger.Errorw("Error trying to parse payload", "error", err)
		return nil, err
	}

	return &cardJSON, nil
}
