package setlisthandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type setlistGetReqTimeframe struct {
	From time.Time `json:"from" binding:"required"`
	To   time.Time `json:"to" binding:"required"`
}

func (slh setlistHandler) GetAll(ctx *gin.Context) {
	context := ctx.Request.Context()
	setlist, err := slh.sls.FetchAll(context)

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"setlist": setlist,
	})
}

func (slh setlistHandler) GetByID(ctx *gin.Context) {
	idField := ctx.Params.ByName("id")

	setlistID, err := strconv.Atoi(idField)
	if err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	context := ctx.Request.Context()
	setlist, err := slh.sls.FetchByID(context, int64(setlistID))

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	setlistEntries, err := slh.sles.FetchBySetlist(context, setlist)

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"setlist": setlist,
		"entries": setlistEntries,
	})
}

func (slh setlistHandler) GetByTimeframe(ctx *gin.Context) {
	var slReq setlistGetReqTimeframe

	if err := ctx.BindJSON(&slReq); err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	context := ctx.Request.Context()
	setlist, err := slh.sls.FetchByTimeframe(context, slReq.From, slReq.To)

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"setlist": setlist,
	})
}
