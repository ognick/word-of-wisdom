//go:build wireinject

package app

import (
	"github.com/google/wire"
)

func InitializeApp() (*App, error) {
	wire.Build(Application)
	return nil, nil
}
