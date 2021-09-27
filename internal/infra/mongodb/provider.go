package mongodb

import "github.com/google/wire"

var MongoDatabaseSet = wire.NewSet(NewMongoDbClient)
