package handler

import (
	"log"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/handler/bundlehandler"
	"github.com/96Asch/mkvstage-server/internal/handler/mehandler"
	"github.com/96Asch/mkvstage-server/internal/handler/rolehandler"
	"github.com/96Asch/mkvstage-server/internal/handler/setlisthandler"
	"github.com/96Asch/mkvstage-server/internal/handler/setlistrolehandler"
	"github.com/96Asch/mkvstage-server/internal/handler/songhandler"
	"github.com/96Asch/mkvstage-server/internal/handler/tokenhandler"
	userhandler "github.com/96Asch/mkvstage-server/internal/handler/userhandler"
	"github.com/96Asch/mkvstage-server/internal/handler/userrolehandler"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Router *gin.Engine
	U      domain.UserService
	T      domain.TokenService
	MH     domain.MiddlewareHandler
	B      domain.BundleService
	S      domain.SongService
	R      domain.RoleService
	UR     domain.UserRoleService
	SL     domain.SetlistService
	SE     domain.SetlistEntryService
	SLR    domain.SetlistRoleService
}

func (cfg *Config) New() *Config {
	return &Config{
		Router: gin.Default(),
	}
}

func Initialize(config *Config) {
	log.Println("Initializing handlers...")

	base := config.Router.Group("api")
	version1 := base.Group("v1")

	ug := userhandler.Initialize(version1, config.U, config.T)
	tokenhandler.Initialize(version1, config.T, config.U)
	mehandler.Initialize(ug, config.U, config.T, config.MH)
	bundlehandler.Initialize(version1, config.B, config.MH)
	songhandler.Initialize(version1, config.S, config.MH)
	rolehandler.Initialize(version1, config.R, config.MH)
	userrolehandler.Initialize(version1, config.UR, config.MH)
	setlisthandler.Initialize(version1, config.SL, config.SE, config.S, config.MH)
	setlistrolehandler.Initialize(version1, config.SLR, config.MH)
}
