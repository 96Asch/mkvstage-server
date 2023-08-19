package songhandler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/96Asch/mkvstage-server/backend/internal/util"
	"github.com/gin-gonic/gin"
)

func (sh songHandler) GetByID(ctx *gin.Context) {
	idField := ctx.Params.ByName("id")

	songID, err := strconv.Atoi(idField)
	if err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	context := ctx.Request.Context()

	song, err := sh.ss.FetchByID(context, int64(songID))
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"song": song})
}

func (sh songHandler) Get(ctx *gin.Context) {
	idsQuery := ctx.Query("ids")
	bundlesQuery := ctx.Query("bids")
	creatorsQuery := ctx.Query("cids")
	titleSearchQuery := ctx.Query("title")
	keysQuery := ctx.Query("keys")
	bpmsQuery := ctx.Query("bpms")

	ids, err := util.StringToInt64Slice(idsQuery)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	bids, err := util.StringToInt64Slice(bundlesQuery)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	cids, err := util.StringToInt64Slice(creatorsQuery)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	keys := strings.Split(keysQuery, ",")
	if keys[0] == "" {
		keys = nil
	}

	bpms, err := util.StringToUintSlice(bpmsQuery)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	filterOptions := &domain.SongFilterOptions{
		IDs:   ids,
		BIDs:  bids,
		CIDs:  cids,
		Keys:  keys,
		Title: titleSearchQuery,
		Bpms:  bpms,
	}

	context := ctx.Request.Context()

	songs, err := sh.ss.Fetch(context, filterOptions)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"songs": songs})
}
