package userrolehandler

import (
	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type userRoleHandler struct {
	urs domain.UserRoleService
}

func Initialize(group *gin.RouterGroup, urs domain.UserRoleService, mwh domain.MiddlewareHandler) {

	urh := &userRoleHandler{
		urs: urs,
	}

	userroles := group.Group("userroles")
	userroles.PATCH("update", mwh.AuthenticateUser(), urh.UpdateBatch)
	userroles.GET("me", mwh.AuthenticateUser(), urh.Me)
	userroles.GET("", urh.GetAll)
}
