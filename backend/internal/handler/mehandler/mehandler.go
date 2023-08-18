package mehandler

import (
	"log"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type meHandler struct {
	userService  domain.UserService
	tokenService domain.TokenService
}

func Initialize(
	group *gin.RouterGroup,
	userService domain.UserService,
	tokenService domain.TokenService,
	middleWare domain.MiddlewareHandler,
) {
	log.Println("Setting up me handlers")

	mehandler := &meHandler{
		userService:  userService,
		tokenService: tokenService,
	}

	me := group.Group("me", middleWare.AuthenticateUser())
	me.GET("", mehandler.Me)
	me.PUT("/update", mehandler.Update)
	me.DELETE("/delete", mehandler.Delete)
}
