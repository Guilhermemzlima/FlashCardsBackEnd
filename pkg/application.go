package pkg

import (
	"FlashCardsBackEnd/internal/config/log"
	"FlashCardsBackEnd/pkg/api/routers"
)

type Application struct {
	SystemRoutes routers.SystemRoutes
}

func NewApplication(systemRoutes routers.SystemRoutes) Application {
	log.Logger.Info("Creating System Main Routers")
	return Application{
		SystemRoutes: systemRoutes,
	}
}
