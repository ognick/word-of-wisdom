package v1

import (
	"net/http"

	"github.com/google/wire"
)

func ProvideHTTPEngine(handler *Handler) http.Handler {
	return handler.Init()
}

var Set = wire.NewSet(ProofOfWorkLimiter, NewHandler, ProvideHTTPEngine)
