package userhandler

import (
	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type loginCredentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (us userHandler) Login(ctx *gin.Context) {

	var creds loginCredentials
	if err := ctx.BindJSON(&creds); err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
		return
	}

	context := ctx.Request.Context()
	_, err := us.userService.Authorize(context, creds.Email, creds.Password)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	// TODO Add JWT token and send back Access token + Refresh token

}