package search_handler

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/errors"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/render"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/filter"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/usecase/search_usecase"
	"net/http"
	"strconv"
)

const (
	headerUserId = "userId"
)

type SearchHandler struct {
	uc search_usecase.ISearchUseCase
}

func NewSearchHandler(service search_usecase.ISearchUseCase) SearchHandler {
	return SearchHandler{uc: service}
}
func (handler *SearchHandler) FindByFilters(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(headerUserId)
	filterString := r.URL.Query().Get("filter")
	pagination := handler.buildPagination(r)

	search, count, err := handler.uc.Search(filterString, userID, pagination)
	if err != nil {
		log.Logger.Errorw("Failed to find by filters", "error", err)
		render.ResponseError(w, err, GenerateHTTPErrorStatusCode(err))
		return
	}

	log.Logger.Debug("Deck has find successfully")
	w.Header().Add("X-Total", strconv.FormatInt(count, 10))
	render.Response(w, search, http.StatusOK)
}

func (handler *SearchHandler) buildPagination(r *http.Request) (pagination *filter.Pagination) {
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
