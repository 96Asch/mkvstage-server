package bundlehandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type bundleReq struct {
	Name     string `json:"name" binding:"required"`
	ParentID int64  `json:"parent_id"`
}

func (bh bundleHandler) Create(ctx *gin.Context) {
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
	if err := bh.bs.Store(context, bundle); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"bundle": bundle})

}
