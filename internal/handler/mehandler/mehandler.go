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

func Initialize(rg *gin.RouterGroup, us domain.UserService, ts domain.TokenService, mwh domain.MiddlewareHandler) {
	log.Println("Setting up me handlers")
	mh := &meHandler{
		userService:  us,
		tokenService: ts,
	}

	me := rg.Group("me", mwh.AuthenticateUser())
	me.GET("", mh.Me)
	me.PATCH("/update", mh.Update)
	me.DELETE("/delete", mh.Delete)
	me.DELETE("/logout", mh.Logout)
}
