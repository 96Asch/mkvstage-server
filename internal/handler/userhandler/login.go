package userhandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type loginCredentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type tokenPair struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

func (uh userHandler) Login(ctx *gin.Context) {
	var creds loginCredentials
	if err := ctx.BindJSON(&creds); err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	context := ctx.Request.Context()

	user, err := uh.userService.Authorize(context, creds.Email, creds.Password)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	refresh, err := uh.tokenService.CreateRefresh(context, user.ID, "")
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	access, err := uh.tokenService.CreateAccess(context, refresh.Refresh)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	tokens := &tokenPair{
		Access:  access.Access,
		Refresh: refresh.Refresh,
	}

	ctx.JSON(http.StatusOK, gin.H{"tokens": tokens})
}
