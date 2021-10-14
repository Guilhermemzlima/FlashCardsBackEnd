package pkg

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/log"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/routers"
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
