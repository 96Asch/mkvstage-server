package userhandler

import (
	"log"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService domain.UserService
}

func Initialize(rg *gin.RouterGroup, userService domain.UserService) {
	log.Println("Setting up user handlers")
	uh := &UserHandler{
		userService: userService,
	}

	us := rg.Group("users")

	us.GET("/", uh.GetAll)
	us.GET("/me", uh.Me)
	us.POST("/create", uh.Create)
	us.PATCH("/me/update", uh.Update)
	us.DELETE("/me/delete", uh.Delete)
}
