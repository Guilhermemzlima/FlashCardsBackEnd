package pkg

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/infra/mongodb"
	handlers "github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/handler"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/routers"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/repository"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/usecase"
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
