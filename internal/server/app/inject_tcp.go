//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	commonconfig "github.com/ognick/word_of_wisdom/internal/common/config"
	"github.com/ognick/word_of_wisdom/pkg/tcp"
)

func provideTCPAddr(cfg commonconfig.Config) tcp.Addr {
	return tcp.Addr(cfg.TCPAddress)
}

func initTCPServer() (*tcp.Server, error) {
	wire.Build(
		Application,
		tcp.NewServer,
		provideTCPAddr,
	)
	return nil, nil
}
