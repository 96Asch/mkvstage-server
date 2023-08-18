package songhandler

import (
	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type songHandler struct {
	ss domain.SongService
}

func Initialize(group *gin.RouterGroup, ss domain.SongService, mwh domain.MiddlewareHandler) {
	songhandler := &songHandler{
		ss: ss,
	}

	songs := group.Group("songs")
	songs.POST("/", mwh.AuthenticateUser(), songhandler.Create)
	songs.GET("", songhandler.GetAll)
	songs.GET(":id", songhandler.GetByID)
	songs.DELETE(":id", mwh.AuthenticateUser(), songhandler.DeleteByID)
	songs.PUT(":id", mwh.AuthenticateUser(), songhandler.UpdateByID)
}
