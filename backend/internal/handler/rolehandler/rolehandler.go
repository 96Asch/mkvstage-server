package rolehandler

import (
	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type roleHandler struct {
	rs domain.RoleService
}

func Initialize(group *gin.RouterGroup, roleservice domain.RoleService, middleware domain.MiddlewareHandler) {
	rolehandler := roleHandler{
		rs: roleservice,
	}

	roles := group.Group("roles")
	roles.POST("create", middleware.AuthenticateUser(), rolehandler.Create)
	roles.GET("", rolehandler.GetAll)
	roles.PUT(":id/update", middleware.AuthenticateUser(), rolehandler.UpdateByID)
	roles.DELETE(":id/delete", middleware.AuthenticateUser(), rolehandler.DeleteByID)
}
