package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Healthz 健康检查
func Healthz(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusOK, "ok")
}

// Readyz 就绪检查
func Readyz(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusOK, "ok")
}
