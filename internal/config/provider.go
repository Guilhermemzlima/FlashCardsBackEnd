package config

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/internal/config/validator"
	"github.com/google/wire"
)

var AppConfigSet = wire.NewSet(validator.NewValidate)
