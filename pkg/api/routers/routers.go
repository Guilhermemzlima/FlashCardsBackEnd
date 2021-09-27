package routers

import (
	"FlashCardsBackEnd/internal/config/log"
	"FlashCardsBackEnd/pkg/api/handler/playlist_handler"
	"FlashCardsBackEnd/pkg/api/middleware"
	"FlashCardsBackEnd/pkg/model/routers"
	"github.com/gorilla/mux"
	"net/http"
)

type SystemRoutes struct {
	playlistHandler playlist_handler.PlaylistHandler
}

func (sys *SystemRoutes) SetupHandler() http.Handler {
	r := mux.NewRouter().PathPrefix(routers.BasePath).Subrouter()

	r.HandleFunc(routers.PlaylistPath, sys.playlistHandler.Post).Methods(http.MethodPost)
	r.HandleFunc(routers.PlaylistPathAll, sys.playlistHandler.FindByUserIdAndPublic).Methods(http.MethodGet)
	r.HandleFunc(routers.PlaylistPathId, sys.playlistHandler.FindById).Methods(http.MethodGet)
	r.HandleFunc(routers.PlaylistPath, sys.playlistHandler.FindByUserId).Methods(http.MethodGet)
	r.HandleFunc(routers.PlaylistPathId, sys.playlistHandler.Delete).Methods(http.MethodDelete)
	r.HandleFunc(routers.PlaylistPathId, sys.playlistHandler.Patch).Methods(http.MethodPatch)

	r.Use(middleware.Header)
	return r
}
func NewSystemRoutes(playlistHandler playlist_handler.PlaylistHandler) SystemRoutes {
	log.Logger.Info("Creating System Main Routers")
	return SystemRoutes{
		playlistHandler: playlistHandler,
	}
}
