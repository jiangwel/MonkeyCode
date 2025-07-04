package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/chaitin/MonkeyCode/backend/consts"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type ActiveMiddleware struct {
	redis  *redis.Client
	logger *slog.Logger
}

func NewActiveMiddleware(redis *redis.Client, logger *slog.Logger) *ActiveMiddleware {
	return &ActiveMiddleware{
		redis:  redis,
		logger: logger,
	}
}

func (a *ActiveMiddleware) Active() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if apikey := GetApiKey(c); apikey != nil {
				if err := a.redis.Set(context.Background(), fmt.Sprintf(consts.UserActiveKeyFmt, apikey.UserID), time.Now().Unix(), 0).Err(); err != nil {
					a.logger.With("error", err).ErrorContext(c.Request().Context(), "failed to set user active status in Redis")
				}
			}
			return next(c)
		}
	}
}
