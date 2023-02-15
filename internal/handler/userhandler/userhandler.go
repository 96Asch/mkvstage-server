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

func Initialize(rg *gin.RouterGroup, us domain.UserService, ts domain.TokenService) *gin.RouterGroup {
	log.Println("Setting up user handlers")
	uh := &userHandler{
		userService:  us,
		tokenService: ts,
	}

	users := rg.Group("users")

	users.GET("", uh.GetAll)
	users.POST("/create", uh.Create)
	users.POST("/login", uh.Login)

	return users
}
