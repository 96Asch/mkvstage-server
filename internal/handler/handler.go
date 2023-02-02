package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func Initialize(config *Config) {

	h := Handler{}

	base := config.Router.Group(os.Getenv("API_BASE"))
	v1 := base.Group("v1")

	user := v1.Group("users")
	user.GET("/me", h.Me)

}

func (h *Handler) Me(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Hello me!")
}
