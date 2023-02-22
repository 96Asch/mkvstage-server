package rolehandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type roleCreateReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (rh roleHandler) Create(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
		return
	}

	var rReq roleCreateReq
	if err := ctx.BindJSON(&rReq); err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
		return
	}

	role := &domain.Role{
		Name:        rReq.Name,
		Description: rReq.Description,
	}

	user := val.(*domain.User)
	context := ctx.Request.Context()
	if err := rh.rs.Store(context, role, user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"role": role})
}
