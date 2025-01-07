//go:build wireinject
// +build wireinject

package app

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	commonconfig "github.com/ognick/word_of_wisdom/internal/common/config"
	"github.com/ognick/word_of_wisdom/internal/server/docs"
	wisdomhttpv1 "github.com/ognick/word_of_wisdom/internal/server/internal/services/wisdom/api/http/v1"
	"github.com/ognick/word_of_wisdom/pkg/http"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Word of Wisdom API
// @version 1.0
// @description This is a simple API for getting wisdoms
// @contact.name Dmitry Aleksandrov
// @contact.email ogneslav.work@gmail.com
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
func registerHTTPHandlers(
	wisdomV1 *wisdomhttpv1.Handler,
) nethttp.Handler {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
	)

	// Init health check
	router.GET("/health", func(c *gin.Context) {
		c.String(nethttp.StatusOK, "ok")
	})

	wisdomV1.Register(router)

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}

func provideHTTPAddr(cfg commonconfig.Config) http.Addr {
	return http.Addr(cfg.HTTPAddress)
}

func initHttpServer() (*http.Server, error) {
	wire.Build(
		Application,
		http.NewServer,
		provideHTTPAddr,
		registerHTTPHandlers,
	)
	return nil, nil
}
