package search_handler

import (
	"FlashCardsBackEnd/internal/config/log"
	"FlashCardsBackEnd/internal/errors"
	"FlashCardsBackEnd/pkg/api/render"
	"FlashCardsBackEnd/pkg/usecase/search_usecase"
	"net/http"
)

const (
	headerUserId = "userId"
	pathVarID    = "id"
)

type SearchHandler struct {
	uc search_usecase.ISearchUseCase
}

func NewSearchHandler(service search_usecase.ISearchUseCase) SearchHandler {
	return SearchHandler{uc: service}
}
func (handler *SearchHandler) FindByFilters(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(headerUserId)
	filter := r.URL.Query().Get("filter")
	search, err := handler.uc.Search(filter, userID)
	if err != nil {
		log.Logger.Errorw("Failed to find by filters", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Deck has find successfully")
	render.Response(w, search, http.StatusOK)
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
