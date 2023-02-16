package handler

import (
	"log"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/handler/bundlehandler"
	"github.com/96Asch/mkvstage-server/internal/handler/mehandler"
	"github.com/96Asch/mkvstage-server/internal/handler/tokenhandler"
	userhandler "github.com/96Asch/mkvstage-server/internal/handler/userhandler"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Router *gin.Engine
	U      domain.UserService
	T      domain.TokenService
	MH     domain.MiddlewareHandler
	B      domain.BundleService
}

func (cfg *Config) New() *Config {
	return &Config{
		Router: gin.Default(),
	}
}

func Initialize(config *Config) {

	log.Println("Initializing handlers...")
	base := config.Router.Group("api")
	v1 := base.Group("v1")

	ug := userhandler.Initialize(v1, config.U, config.T)
	tokenhandler.Initialize(v1, config.T, config.U)

	mehandler.Initialize(ug, config.U, config.T, config.MH)

	bundlehandler.Initialize(v1, config.B)
}
