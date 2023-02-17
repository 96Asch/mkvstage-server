package songhandler

import (
	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type songHandler struct {
	ss domain.SongService
}

func Initialize(group *gin.RouterGroup, ss domain.SongService, mwh domain.MiddlewareHandler) {

	sh := &songHandler{
		ss: ss,
	}

	songs := group.Group("songs")
	songs.POST("create", mwh.AuthenticateUser(), sh.Create)
	songs.GET("", sh.GetAll)
	songs.GET(":id", sh.GetByID)

}
