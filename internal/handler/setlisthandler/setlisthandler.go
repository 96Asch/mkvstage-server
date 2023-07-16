package setlisthandler

import (
	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type setlistHandler struct {
	sls  domain.SetlistService
	sles domain.SetlistEntryService
	ss   domain.SongService
}

func Initialize(group *gin.RouterGroup, sls domain.SetlistService, sles domain.SetlistEntryService, ss domain.SongService, mwh domain.MiddlewareHandler) {
	setlisthandler := &setlistHandler{
		sls:  sls,
		sles: sles,
		ss:   ss,
	}

	setlists := group.Group("setlists")
	setlists.POST("create", mwh.AuthenticateUser(), setlisthandler.Create)
	// setlists.GET("", setlisthandler.GetAll)
	// setlists.GET(":id", setlisthandler.GetByID)
	setlists.DELETE(":id/delete", mwh.AuthenticateUser(), setlisthandler.DeleteByID)
	setlists.PUT(":id/update", mwh.AuthenticateUser(), setlisthandler.UpdateByID)
}
