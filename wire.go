//+build wireinject

package main

import (
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg"
	"github.com/google/wire"
)

func SetupApplication() (pkg.Application, error) {
	wire.Build(pkg.Container)
	return pkg.Application{}, nil
}
