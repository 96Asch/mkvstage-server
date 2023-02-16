package bundlehandler

import (
	"log"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type bundleHandler struct {
	bs domain.BundleService
}

func Initialize(rg *gin.RouterGroup, bs domain.BundleService) {
	log.Println("Setting up bundle handlers")
	bh := &bundleHandler{
		bs: bs,
	}

	bundle := rg.Group("bundle")
	bundle.POST("create", bh.Create)

}
