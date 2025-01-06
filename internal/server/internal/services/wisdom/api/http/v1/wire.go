package v1

import "github.com/google/wire"

var Init = wire.NewSet(
	NewHandler,
)
