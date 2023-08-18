package userrolehandler

import (
	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type userRoleHandler struct {
	urs domain.UserRoleService
}

func Initialize(group *gin.RouterGroup, urs domain.UserRoleService, mwh domain.MiddlewareHandler) {
	userrolehandler := &userRoleHandler{
		urs: urs,
	}

	userroles := group.Group("userroles")
	userroles.PATCH("update", mwh.AuthenticateUser(), userrolehandler.UpdateBatch)
	userroles.GET("me", mwh.AuthenticateUser(), userrolehandler.Me)
	userroles.GET("", userrolehandler.GetAll)
}
