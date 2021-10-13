package routers

import (
	"FlashCardsBackEnd/internal/config/log"
	"FlashCardsBackEnd/pkg/api/handler/card_handler"
	"FlashCardsBackEnd/pkg/api/handler/deck_handler"
	"FlashCardsBackEnd/pkg/api/handler/playlist_handler"
	review__handler "FlashCardsBackEnd/pkg/api/handler/review_handler"
	"FlashCardsBackEnd/pkg/api/handler/search_handler"
	"FlashCardsBackEnd/pkg/api/middleware"
	"FlashCardsBackEnd/pkg/model/routers"
	"github.com/gorilla/mux"
	"net/http"
)

type SystemRoutes struct {
	playlistHandler playlist_handler.PlaylistHandler
	deckHandler     deck_handler.DeckHandler
	cardHandler     card_handler.CardHandler
	reviewHandler   review__handler.ReviewHandler
	searchHandler   search_handler.SearchHandler
}

func (sys *SystemRoutes) SetupHandler() http.Handler {
	r := mux.NewRouter().PathPrefix(routers.BasePath).Subrouter()

	r.HandleFunc(routers.PlaylistPath, sys.playlistHandler.Post).Methods(http.MethodPost)
	r.HandleFunc(routers.PlaylistPathAll, sys.playlistHandler.FindByUserIdAndPublic).Methods(http.MethodGet)
	r.HandleFunc(routers.PlaylistPathId, sys.playlistHandler.FindById).Methods(http.MethodGet)
	r.HandleFunc(routers.PlaylistPath, sys.playlistHandler.FindByUserId).Methods(http.MethodGet)
	r.HandleFunc(routers.PlaylistPathId, sys.playlistHandler.Delete).Methods(http.MethodDelete)
	r.HandleFunc(routers.PlaylistPathId, sys.playlistHandler.Patch).Methods(http.MethodPatch)
	r.HandleFunc(routers.PlaylistPathAdd, sys.playlistHandler.PatchDeck).Methods(http.MethodPatch)

	r.HandleFunc(routers.DeckPath, sys.deckHandler.Post).Methods(http.MethodPost)
	r.HandleFunc(routers.DeckPathAll, sys.deckHandler.FindByUserIdAndPublic).Methods(http.MethodGet)
	r.HandleFunc(routers.DeckPathId, sys.deckHandler.FindById).Methods(http.MethodGet)
	r.HandleFunc(routers.DeckPath, sys.deckHandler.FindByUserId).Methods(http.MethodGet)
	r.HandleFunc(routers.DeckPathId, sys.deckHandler.Delete).Methods(http.MethodDelete)
	r.HandleFunc(routers.DeckPathId, sys.deckHandler.Patch).Methods(http.MethodPatch)
	r.HandleFunc(routers.DeckRecentPath, sys.deckHandler.FindRecent).Methods(http.MethodGet)

	r.HandleFunc(routers.CardPathId, sys.cardHandler.Post).Methods(http.MethodPost)
	r.HandleFunc(routers.CardDeckPathId, sys.cardHandler.FindByDeckId).Methods(http.MethodGet)
	r.HandleFunc(routers.CardPathId, sys.cardHandler.Patch).Methods(http.MethodPatch)
	r.HandleFunc(routers.CardPathId, sys.cardHandler.Delete).Methods(http.MethodDelete)
	r.HandleFunc(routers.CardPathId, sys.cardHandler.FindById).Methods(http.MethodGet)

	r.HandleFunc(routers.PlaylistReviewPath, sys.reviewHandler.ReviewPlaylist).Methods(http.MethodGet)
	r.HandleFunc(routers.DeckReviewPath, sys.reviewHandler.ReviewDeck).Methods(http.MethodGet)
	r.HandleFunc(routers.ReviewPathId, sys.reviewHandler.FindById).Methods(http.MethodGet)
	r.HandleFunc(routers.ReviewPathIdWrong, sys.reviewHandler.CardResultWrong).Methods(http.MethodPost)
	r.HandleFunc(routers.ReviewPathIdRight, sys.reviewHandler.CardResultRight).Methods(http.MethodPost)

	r.HandleFunc(routers.SearchPath, sys.searchHandler.FindByFilters).Methods(http.MethodGet)

	r.Use(middleware.Header)
	return r
}
func NewSystemRoutes(playlistHandler playlist_handler.PlaylistHandler, searchHandler search_handler.SearchHandler, reviewHandler review__handler.ReviewHandler, cardHandler card_handler.CardHandler, deckHandler deck_handler.DeckHandler) SystemRoutes {
	log.Logger.Info("Creating System Main Routers")
	return SystemRoutes{
		playlistHandler: playlistHandler,
		deckHandler:     deckHandler,
		cardHandler:     cardHandler,
		reviewHandler:   reviewHandler,
		searchHandler:   searchHandler,
	}
}
