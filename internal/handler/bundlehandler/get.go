package bundlehandler

import (
	"net/http"
	"strconv"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

func (bh bundleHandler) Get(ctx *gin.Context) {

	idField := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idField)
	if err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
		return
	}

	context := ctx.Request.Context()
	bundles, err := bh.bs.FetchByID(context, int64(id))
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"bundle": bundles})
}

func (bh bundleHandler) GetAll(ctx *gin.Context) {

	context := ctx.Request.Context()
	bundles, err := bh.bs.FetchAll(context)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"bundles": bundles})
}
