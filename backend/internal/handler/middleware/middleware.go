package middleware

import "github.com/96Asch/mkvstage-server/internal/domain"

type ginMiddlewareHandler struct {
	US domain.UserService
	TS domain.TokenService
}

//revive:disable:unexported-return
func NewGinMiddlewareHandler(us domain.UserService, ts domain.TokenService) *ginMiddlewareHandler {
	return &ginMiddlewareHandler{
		US: us,
		TS: ts,
	}
}
