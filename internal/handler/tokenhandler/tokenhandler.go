package tokenhandler

import (
	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type tokenHandler struct {
	tokenService domain.TokenService
	userService  domain.UserService
}

func Initialize(rg *gin.RouterGroup, ts domain.TokenService, us domain.UserService) {

	th := &tokenHandler{
		tokenService: ts,
		userService:  us,
	}

	t := rg.Group("tokens")

	t.POST("/create", th.CreateAccess)

}
