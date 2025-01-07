package app

import (
	"github.com/google/wire"

	commonconfig "github.com/ognick/word_of_wisdom/internal/common/config"
	"github.com/ognick/word_of_wisdom/pkg/tcp"
)

func provideTCPAddr(cfg commonconfig.Config) tcp.Addr {
	return tcp.Addr(cfg.TCPAddress)
}

var initTCPServer = wire.NewSet(
	tcp.NewServer,
	provideTCPAddr,
)
