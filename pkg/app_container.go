package pkg

import (
	"FlashCardsBackEnd/internal/config"
	"FlashCardsBackEnd/internal/infra/mongodb"
	handlers "FlashCardsBackEnd/pkg/api/handler"
	"FlashCardsBackEnd/pkg/api/routers"
	"FlashCardsBackEnd/pkg/repository"
	"FlashCardsBackEnd/pkg/usecase"
	"github.com/google/wire"
)

var Container = wire.NewSet(
	config.AppConfigSet,
	ApplicationSet,
	repository.Set,
	usecase.Set,
	mongodb.MongoDatabaseSet,
	handlers.ApplicationHandlersSet,
	routers.RoutesSet,
)
