package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/chaitin/MonkeyCode/backend/config"
	"github.com/chaitin/MonkeyCode/backend/errcode"
)

type ReadOnlyMiddleware struct {
	cfg *config.Config
}

func NewReadOnlyMiddleware(cfg *config.Config) *ReadOnlyMiddleware {
	return &ReadOnlyMiddleware{cfg: cfg}
}

func (m *ReadOnlyMiddleware) Guard() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if m.cfg.ReadOnly && c.Request().Method != http.MethodGet {
				return errcode.ErrReadOnly
			}
			return next(c)
		}
	}
}
