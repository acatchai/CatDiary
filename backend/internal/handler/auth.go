package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// AuthRegister 注册用户
func AuthRegister(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}

// AuthLogin 用户登录
func AuthLogin(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}

// AuthLogout 用户退出登录
func AuthLogout(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}

// AuthMe 获取用户信息
func AuthMe(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}
