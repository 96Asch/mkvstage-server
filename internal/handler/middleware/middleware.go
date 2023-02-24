package middleware

import "github.com/96Asch/mkvstage-server/internal/domain"

type ginMiddlewareHandler struct {
	TS domain.TokenService
}

//revive:disable:unexported-return
func NewGinMiddlewareHandler(ts domain.TokenService) *ginMiddlewareHandler {
	return &ginMiddlewareHandler{
		TS: ts,
	}
}
