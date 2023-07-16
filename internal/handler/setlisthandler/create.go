package setlisthandler

import (
	"net/http"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type setlistCreateReq struct {
	Name           string                `json:"name" binding:"required"`
	CreatorID      int64                 `json:"creator_id" binding:"required"`
	Deadline       time.Time             `json:"deadline" binding:"required"`
	CreatedEntries []domain.SetlistEntry `json:"created_entries" binding:"required"`
}

func (slh setlistHandler) Create(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	var slReq setlistCreateReq
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
		Name:      slReq.Name,
		CreatorID: slReq.CreatorID,
		Deadline:  slReq.Deadline.Local(),
	}

	setlistEntries := make([]domain.SetlistEntry, len(slReq.CreatedEntries))
	copy(setlistEntries, slReq.CreatedEntries)

	context := ctx.Request.Context()
	if err := slh.sls.Store(context, setlist, user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	if err := slh.sles.StoreBatch(context, &setlistEntries, user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"setlist": setlist,
		"entries": setlistEntries,
	})
}
