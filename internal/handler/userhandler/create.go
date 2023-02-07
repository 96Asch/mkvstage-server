package userhandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type newUser struct {
	Email        string `json:"email" gorm:"unique" binding:"required,email"`
	Password     string `json:"password" binding:"required"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	ProfileColor string `json:"profile_color" binding:"required"`
}

func (u *UserHandler) Create(ctx *gin.Context) {

	var nUser newUser
	if err := ctx.BindJSON(&nUser); err != nil {
		newError := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newError), gin.H{"error": newError})
		return
	}

	user := domain.User{
		Email:        nUser.Email,
		Password:     nUser.Password,
		FirstName:    nUser.FirstName,
		LastName:     nUser.LastName,
		Permission:   domain.GUEST,
		ProfileColor: nUser.ProfileColor,
	}

	context := ctx.Request.Context()
	if err := u.userService.Store(context, &user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}
