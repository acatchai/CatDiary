package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// UserMe 获取用户信息
func UserMe(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}

// UserPatchMe 更新用户信息
func UserPatchMe(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}

// UserPatchPassword 更新用户密码
func UserPatchPassword(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}
