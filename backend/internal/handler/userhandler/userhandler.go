package userhandler

import (
	"log"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService domain.UserService
}

func Initialize(group *gin.RouterGroup, us domain.UserService, mwh domain.MiddlewareHandler) *gin.RouterGroup {
	log.Println("Setting up user handlers")

	userhandler := &userHandler{
		userService: us,
	}

	users := group.Group("users")

	users.GET("", userhandler.GetAll)
	users.POST("/create", mwh.JWTExtractEmail(), userhandler.Create)

	return users
}
