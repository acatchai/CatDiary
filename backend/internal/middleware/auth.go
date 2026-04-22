package middleware

import (
	"context"
	"strings"

	"github.com/acatchai/catdiary/backend/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	hutils "github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func RequireAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 期望格式: Authorization: Bearer <token>
		authHeader := string(c.Request.Header.Peek("Authorization"))
		if authHeader == "" {
			c.JSON(consts.StatusUnauthorized, hutils.H{
				"error": "缺少Authorization头部",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(consts.StatusUnauthorized, hutils.H{
				"error": "Authorization 格式错误，应为 Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		userID, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.JSON(consts.StatusUnauthorized, hutils.H{
				"error": "Token 无效或过期",
			})
			c.Abort()
			return
		}

		// 将解析出的 userID 存入 context 供后续 handler 取用
		c.Set("user_id", userID)
		c.Next(ctx)
	}
}
