package userhandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/gin-gonic/gin"
)

type userCreateReq struct {
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	ProfileColor string `json:"profile_color" binding:"required"`
}

func (uh userHandler) Create(ctx *gin.Context) {
	val, exists := ctx.Get("email")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(newErr.Status(), gin.H{"error": newErr.Error()})

		return
	}

	email, ok := val.(string)
	if !ok {
		newErr := domain.NewInternalErr()
		ctx.JSON(newErr.Status(), gin.H{"error": newErr.Error()})

		return
	}

	var nUser userCreateReq
	if err := util.BindModel(ctx, &nUser); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	user := domain.User{
		Email:        email,
		FirstName:    nUser.FirstName,
		LastName:     nUser.LastName,
		Permission:   domain.GUEST,
		ProfileColor: nUser.ProfileColor,
	}

	context := ctx.Request.Context()
	if err := uh.userService.Store(context, &user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}
