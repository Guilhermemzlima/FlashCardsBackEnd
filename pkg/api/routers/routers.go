package routers

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/handler/card_handler"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/handler/deck_handler"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/handler/playlist_handler"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/handler/review_handler"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/handler/search_handler"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/middleware"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/model/routers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

type SystemRoutes struct {
	playlistHandler playlist_handler.PlaylistHandler
	deckHandler     deck_handler.DeckHandler
	cardHandler     card_handler.CardHandler
	reviewHandler   review_handler.ReviewHandler
	searchHandler   search_handler.SearchHandler
}

func (sys *SystemRoutes) SetupHandler() http.Handler {
	r := mux.NewRouter().PathPrefix(routers.BasePath).Subrouter()

	r.HandleFunc(routers.PlaylistPath, sys.playlistHandler.Post).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc(routers.PlaylistPathAll, sys.playlistHandler.FindByUserIdAndPublic).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(routers.PlaylistPathId, sys.playlistHandler.FindById).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(routers.PlaylistPath, sys.playlistHandler.FindByUserId).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(routers.PlaylistPathId, sys.playlistHandler.Delete).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc(routers.PlaylistPathId, sys.playlistHandler.Patch).Methods(http.MethodPatch, http.MethodOptions)
	r.HandleFunc(routers.PlaylistPathAdd, sys.playlistHandler.PatchDeck).Methods(http.MethodPatch, http.MethodOptions)
	r.HandleFunc(routers.PlaylistPathDelete, sys.playlistHandler.RemoveDeckFromPlaylist).Methods(http.MethodPatch, http.MethodOptions)
	r.HandleFunc(routers.PlaylistFindDecks, sys.playlistHandler.FindDecksByPlaylistId).Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc(routers.DeckPath, sys.deckHandler.Post).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc(routers.DeckPathAll, sys.deckHandler.FindByUserIdAndPublic).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(routers.DeckPathId, sys.deckHandler.FindById).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(routers.DeckPath, sys.deckHandler.FindByUserId).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(routers.DeckPathId, sys.deckHandler.Delete).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc(routers.DeckPathId, sys.deckHandler.Patch).Methods(http.MethodPatch, http.MethodOptions)
	r.HandleFunc(routers.DeckRecentPath, sys.deckHandler.FindRecent).Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc(routers.CardPathId, sys.cardHandler.Post).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc(routers.CardDeckPathId, sys.cardHandler.FindByDeckId).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(routers.CardPathId, sys.cardHandler.Patch).Methods(http.MethodPatch, http.MethodOptions)
	r.HandleFunc(routers.CardPathId, sys.cardHandler.Delete).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc(routers.CardPathId, sys.cardHandler.FindById).Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc(routers.PlaylistReviewPath, sys.reviewHandler.ReviewPlaylist).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(routers.DeckReviewPath, sys.reviewHandler.ReviewDeck).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(routers.ReviewPathId, sys.reviewHandler.FindById).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(routers.ReviewPathIdWrong, sys.reviewHandler.CardResultWrong).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc(routers.ReviewPathIdRight, sys.reviewHandler.CardResultRight).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc(routers.SearchPath, sys.searchHandler.FindByFilters).Methods(http.MethodGet, http.MethodOptions)

	r.Use(middleware.Header)
	return configCORS().Handler(r)
}
func NewSystemRoutes(playlistHandler playlist_handler.PlaylistHandler, deckHandler deck_handler.DeckHandler,
	cardHandler card_handler.CardHandler, reviewHandler review_handler.ReviewHandler, searchHandler search_handler.SearchHandler) SystemRoutes {
	log.Logger.Info("Creating System Main Routers")
	return SystemRoutes{
		playlistHandler: playlistHandler,
		deckHandler:     deckHandler,
		cardHandler:     cardHandler,
		reviewHandler:   reviewHandler,
		searchHandler:   searchHandler,
	}
}

func configCORS() *cors.Cors {
	var c = cors.New(cors.Options{
		AllowedOrigins:  []string{"*"},
		AllowOriginFunc: AllowOriginFunc,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:     []string{"*"},
		ExposedHeaders:     []string{"*"},
		MaxAge:             0,
		AllowCredentials:   false,
		OptionsPassthrough: false,
		Debug:              false,
	})

	log.Logger.Debug("Cors handler with routes created successfully att aqui ")
	return c
}

func AllowOriginFunc(_ string) bool {
	return true
}
