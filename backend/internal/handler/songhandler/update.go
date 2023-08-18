package songhandler

import (
	"net/http"
	"strconv"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type songUpdateReq struct {
	BundleID   int64  `json:"bundle_id" binding:"required"`
	CreatorID  int64  `json:"creator_id" binding:"required"`
	Title      string `json:"title" binding:"required,lte=255"`
	Subtitle   string `json:"subtitle" binding:"required,lte=255"`
	Key        string `json:"key" binding:"required"`
	Bpm        uint   `json:"bpm" binding:"required"`
	ChordSheet string `json:"chord_sheet" binding:"required"`
}

func (sh songHandler) UpdateByID(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	idField := ctx.Params.ByName("id")

	songID, err := strconv.Atoi(idField)
	if err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	var sReq songUpdateReq
	if err := ctx.BindJSON(&sReq); err != nil {
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

	context := ctx.Request.Context()
	song := &domain.Song{
		ID:         int64(songID),
		BundleID:   sReq.BundleID,
		CreatorID:  sReq.CreatorID,
		Title:      sReq.Title,
		Subtitle:   sReq.Subtitle,
		Key:        sReq.Key,
		Bpm:        sReq.Bpm,
		ChordSheet: datatypes.JSON([]byte(sReq.ChordSheet)),
	}

	err = sh.ss.Update(context, song, user)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"song": song})
}
