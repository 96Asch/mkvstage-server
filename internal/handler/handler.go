package handler

import (
	"github.com/96Asch/mkvstage-server/internal/domain"
	userhandler "github.com/96Asch/mkvstage-server/internal/handler/userhandler"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Router *gin.Engine
	U      domain.UserService
}

func (cfg *Config) New() *Config {
	return &Config{
		Router: gin.Default(),
	}
}

func Initialize(config *Config) {

	base := config.Router.Group("api")
	v1 := base.Group("v1")

	userhandler.Initialize(v1, config.U)
}
