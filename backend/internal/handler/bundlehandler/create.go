package bundlehandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type bundleReq struct {
	Name     string `json:"name" binding:"required,lte=255"`
	ParentID int64  `json:"parent_id"`
}

func (bh bundleHandler) Create(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	var bReq bundleReq
	if err := ctx.BindJSON(&bReq); err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	bundle := &domain.Bundle{
		Name:     bReq.Name,
		ParentID: bReq.ParentID,
	}

	context := ctx.Request.Context()

	principal, ok := val.(*domain.User)
	if !ok {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	if err := bh.bs.Store(context, bundle, principal); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"bundle": bundle})
}
