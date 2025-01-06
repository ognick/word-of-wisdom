package v1

import (
	"net"

	"github.com/google/wire"
)

func ProvideTCPHandle(handler *Handler) func(net.Conn) {
	return handler.Handle
}

var Init = wire.NewSet(
	NewHandler,
	ProvideTCPHandle,
)
