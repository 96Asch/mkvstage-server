package setlistrolehandler

import (
	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type setlistRoleHandler struct {
	slrs domain.SetlistRoleService
}

func Initialize(group *gin.RouterGroup, slrs domain.SetlistRoleService, mwh domain.MiddlewareHandler) {
	setlistRolehandler := &setlistRoleHandler{
		slrs: slrs,
	}

	setlistRole := group.Group("setlistroles")
	setlistRole.POST("", mwh.AuthenticateUser(), setlistRolehandler.Create)
	setlistRole.GET("", setlistRolehandler.GetAll)
	setlistRole.DELETE("", mwh.AuthenticateUser(), setlistRolehandler.Delete)
}
