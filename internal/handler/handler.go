package handler

import (
	"os"

	"github.com/96Asch/mkvstage-server/internal/domain"
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

	base := config.Router.Group(os.Getenv("API_BASE"))
	v1 := base.Group("v1")

	NewUserHandler(v1, &config.U)

}
