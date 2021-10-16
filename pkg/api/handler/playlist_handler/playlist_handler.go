package playlist_handler

import (
	"encoding/json"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/errors"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/render"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/filter"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/playlist"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/usecase/playlist_usecase"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	headerUserId = "userId"
	pathVarID    = "id"
)

type PlaylistHandler struct {
	playlistUseCase playlist_usecase.IPlaylistUseCase
}

func NewPlaylistHandler(service playlist_usecase.IPlaylistUseCase) PlaylistHandler {
	return PlaylistHandler{playlistUseCase: service}
}

func (handler *PlaylistHandler) Post(w http.ResponseWriter, r *http.Request) {
	requestBody, err := handler.extractBody(r)
	if err != nil {
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(errors.ErrInvalidPayload))
		return
	}

	userID := r.Header.Get(headerUserId)
	result, err := handler.playlistUseCase.Create(userID, requestBody)
	if err != nil {
		log.Logger.Errorw("Failed to create playlist", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Playlist has created successfully")
	render.Response(w, result, http.StatusCreated)
}

func (handler *PlaylistHandler) FindById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[pathVarID]
	userID := r.Header.Get(headerUserId)

	result, err := handler.playlistUseCase.FindById(userID, id)
	if err != nil {
		log.Logger.Errorw("Failed to find playlist", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Playlist has find successfully")
	render.Response(w, result, http.StatusOK)
}

func (handler *PlaylistHandler) FindByUserIdAndPublic(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(headerUserId)

	pagination := handler.buildPagination(r)
	result, count, err := handler.playlistUseCase.FindByUserId(userID, pagination, false)
	if err != nil {
		log.Logger.Errorw("Failed to find playlist", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	w.Header().Add("X-Total", strconv.FormatInt(count, 10))
	render.Response(w, result, http.StatusOK)
}

func (handler *PlaylistHandler) FindDecksByPlaylistId(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(headerUserId)
	id := mux.Vars(r)[pathVarID]
	result, err := handler.playlistUseCase.FindDecksOnPlaylist(userID, id)
	if err != nil {
		log.Logger.Errorw("Failed to find playlist", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	render.Response(w, result, http.StatusOK)
}

func (handler *PlaylistHandler) FindByUserId(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(headerUserId)

	pagination := handler.buildPagination(r)
	result, count, err := handler.playlistUseCase.FindByUserId(userID, pagination, true)
	if err != nil {
		log.Logger.Errorw("Failed to find playlist", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	w.Header().Add("X-Total", strconv.FormatInt(count, 10))
	render.Response(w, result, http.StatusOK)
}

func (handler *PlaylistHandler) Patch(w http.ResponseWriter, r *http.Request) {
	requestBody, err := handler.extractBody(r)
	if err != nil {
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(errors.ErrInvalidPayload))
		return
	}

	id := mux.Vars(r)[pathVarID]
	userID := r.Header.Get(headerUserId)
	result, err := handler.playlistUseCase.Update(id, userID, true, requestBody)
	if err != nil {
		log.Logger.Errorw("Failed to update playlist", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Playlist has updated successfully")
	render.Response(w, result, http.StatusCreated)
}

func (handler *PlaylistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	customerID := r.Header.Get(headerUserId)

	id := mux.Vars(r)[pathVarID]
	_, err := handler.playlistUseCase.Delete(id, customerID)
	if err != nil {
		log.Logger.Errorw("Failed to remove playlist", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Playlist has removed successfully")
	render.Response(w, nil, http.StatusNoContent)
}

func (handler *PlaylistHandler) extractBody(r *http.Request) (*playlist.Playlist, error) {
	var playlistJSON playlist.Playlist

	if err := json.NewDecoder(r.Body).Decode(&playlistJSON); err != nil {
		log.Logger.Errorw("Error trying to parse payload", "error", err)
		return nil, err
	}

	return &playlistJSON, nil
}

func (handler *PlaylistHandler) PatchDeck(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[pathVarID]
	userID := r.Header.Get(headerUserId)
	requestBody, err := handler.extractBody(r)
	result, err := handler.playlistUseCase.AddDeckToPlaylist(id, userID, requestBody.Id.Hex())
	if err != nil {
		log.Logger.Errorw("Failed to update playlist", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Playlist has updated successfully")
	render.Response(w, result, http.StatusCreated)
}

func (handler *PlaylistHandler) RemoveDeckFromPlaylist(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[pathVarID]
	userID := r.Header.Get(headerUserId)
	requestBody, err := handler.extractBody(r)
	result, err := handler.playlistUseCase.RemoveDeckFromPlaylist(id, userID, requestBody.Id.Hex())
	if err != nil {
		log.Logger.Errorw("Failed to update playlist", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Playlist has updated successfully")
	render.Response(w, result, http.StatusCreated)
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

func (handler *PlaylistHandler) buildPagination(r *http.Request) (pagination *filter.Pagination) {
	reqQuery := r.URL.Query()
	limit, _ := strconv.ParseInt(reqQuery.Get("limit"), 10, 64)
	offset, _ := strconv.ParseInt(reqQuery.Get("offset"), 10, 64)
	return filter.NewPagination(limit, offset)
}
