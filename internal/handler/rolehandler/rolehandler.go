package rolehandler

import (
	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type roleHandler struct {
	rs domain.RoleService
}

func Initialize(group *gin.RouterGroup, rs domain.RoleService, mwh domain.MiddlewareHandler) {
	rh := roleHandler{
		rs: rs,
	}

	roles := group.Group("roles")
	roles.POST("create", mwh.AuthenticateUser(), rh.Create)
	roles.PUT(":id/update", mwh.AuthenticateUser(), rh.UpdateByID)
}
