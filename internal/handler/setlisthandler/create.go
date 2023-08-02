package setlisthandler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type setlistCreateReq struct {
	Name           string                 `json:"name" binding:"required"`
	CreatorID      int64                  `json:"creator_id" binding:"required"`
	Deadline       time.Time              `json:"deadline" binding:"required"`
	CreatedEntries []setlistRoleCreateReq `json:"created_entries" binding:"required,dive"`
}

type setlistRoleCreateReq struct {
	SongID      int64    `json:"song_id" binding:"required"`
	Transpose   int16    `json:"transpose"`
	Notes       string   `json:"notes"`
	Arrangement []string `json:"arrangement"`
}

type setlistResponse struct {
	*domain.Setlist
	Entries *[]domain.SetlistEntry `json:"entries"`
}

func (slh setlistHandler) Create(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	var slReq setlistCreateReq
	if err := util.BindModel(ctx, &slReq); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

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
		Deadline:  slReq.Deadline.Local().Truncate(time.Minute),
	}

	setlistEntries := make([]domain.SetlistEntry, len(slReq.CreatedEntries))

	for idx, entry := range slReq.CreatedEntries {
		jsonArray, err := json.Marshal(entry.Arrangement)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		setlistEntries[idx] = domain.SetlistEntry{
			SongID:      entry.SongID,
			Transpose:   entry.Transpose,
			Notes:       entry.Notes,
			Arrangement: datatypes.JSON(jsonArray),
		}
	}

	context := ctx.Request.Context()
	if err := slh.sls.Store(context, setlist, user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	if err := slh.sles.StoreBatch(context, &setlistEntries, user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	response := setlistResponse{
		setlist,
		&setlistEntries,
	}

	ctx.JSON(http.StatusCreated, gin.H{"setlist": response})
}
