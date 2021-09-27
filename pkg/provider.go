package pkg

import "github.com/google/wire"

var ApplicationSet = wire.NewSet(NewApplication)
