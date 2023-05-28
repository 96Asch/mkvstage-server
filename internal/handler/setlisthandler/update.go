package setlisthandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type setlistUpdateReq struct {
	Name      string    `json:"name" binding:"required"`
	CreatorID int64     `json:"creator_id" binding:"required"`
	Global    bool      `json:"is_global" binding:"required"`
	Deadline  time.Time `json:"deadline" binding:"required"`
	// TODO: Add updated entries
}

func (slh setlistHandler) UpdateByID(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	idField := ctx.Params.ByName("id")

	setlistID, err := strconv.Atoi(idField)
	if err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	var slReq setlistUpdateReq
	if err := ctx.BindJSON(&slReq); err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	user, ok := val.(*domain.User)
	if !ok {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	setlist := &domain.Setlist{
		ID:        int64(setlistID),
		Name:      slReq.Name,
		CreatorID: slReq.CreatorID,
		Global:    slReq.Global,
		Deadline:  slReq.Deadline,
	}

	context := ctx.Request.Context()
	updatedSetlist, err := slh.sls.Update(context, setlist, user)
	// TODO: update entries in setlistentryservice

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"setlist": updatedSetlist})
}
