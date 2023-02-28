package userhandler

import (
	"log"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService  domain.UserService
	tokenService domain.TokenService
}

func Initialize(group *gin.RouterGroup, us domain.UserService, ts domain.TokenService, mwh domain.MiddlewareHandler) *gin.RouterGroup {
	log.Println("Setting up user handlers")

	userhandler := &userHandler{
		userService:  us,
		tokenService: ts,
	}

	users := group.Group("users")

	users.GET("", userhandler.GetAll)
	users.POST("/create", userhandler.Create)
	users.POST("/login", userhandler.Login)
	users.PUT("/setperm", mwh.AuthenticateUser(), userhandler.ChangePermissionByID)

	return users
}
