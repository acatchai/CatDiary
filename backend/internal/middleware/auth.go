package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func RequireAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := string(c.Request.Header.Peek("Authorization"))
		if token == "" {
			c.Status(consts.StatusUnauthorized)
			return
		}
		c.Next(ctx)
	}
}
