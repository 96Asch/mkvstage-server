package tokenhandler

import (
	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type tokenHandler struct {
	tokenService domain.TokenService
	userService  domain.UserService
}

func Initialize(group *gin.RouterGroup, ts domain.TokenService, us domain.UserService) {
	tokenhandler := &tokenHandler{
		tokenService: ts,
		userService:  us,
	}

	t := group.Group("tokens")

	t.POST("/renew", tokenhandler.CreateAccess)
}
