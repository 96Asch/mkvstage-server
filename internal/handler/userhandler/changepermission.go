package userhandler

import (
	"log"
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type permissionReq struct {
	RecipientID int64            `json:"recipient_id" binding:"required"`
	Permission  domain.Clearance `json:"permission" binding:"required"`
}

func (uh userHandler) ChangePermissionByID(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		err := domain.NewInternalErr()

		ctx.JSON(domain.Status(err), gin.H{"error": err})
		log.Printf("User does not exist in context")

		return
	}

	principal, ok := val.(*domain.User)
	if !ok {
		err := domain.NewInternalErr()
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	var permReq permissionReq
	if err := ctx.BindJSON(&permReq); err != nil {
		newError := domain.NewBadRequestErr(err.Error())
		log.Println(err)
		ctx.JSON(domain.Status(newError), gin.H{"error": newError})

		return
	}

	context := ctx.Request.Context()

	updatedUser, err := uh.userService.SetPermission(context, permReq.Permission, permReq.RecipientID, principal)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": updatedUser})
}
