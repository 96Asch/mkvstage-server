package userhandler

import (
	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService domain.UserService
}

func Initialize(rg *gin.RouterGroup, userService domain.UserService) {
	uh := &UserHandler{
		userService: userService,
	}

	us := rg.Group("users")

	us.GET("/", uh.GetAll)
	us.GET("/me", uh.Me)
	us.POST("/create", uh.Create)
	us.PATCH("/update", uh.UpdateByID)
}
