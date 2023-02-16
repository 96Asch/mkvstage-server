package bundlehandler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type updateReq struct {
	Name     string `json:"name" binding:"required"`
	ParentID int64  `json:"parent_id"`
}

func (bh bundleHandler) UpdateByID(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
		log.Println(newErr)
		return
	}
	principal := val.(*domain.User)

	idField := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idField)
	if err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
		return
	}

	var req updateReq
	if err := ctx.BindJSON(&req); err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
		return
	}

	bundle := &domain.Bundle{
		ID:       int64(id),
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	context := ctx.Request.Context()
	err = bh.bs.Update(context, bundle, principal)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"bundle": bundle})
}
