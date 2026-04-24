package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// DiaryCreate 创建日记
func DiaryCreate(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}

// DiaryList 获取日记列表
func DiaryList(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}

// DiaryGet 获取日记详情
func DiaryGet(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}

// DiaryPut 更新日记
func DiaryPut(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}

// DiaryPatch 更新日记内容
func DiaryPatch(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}

// DiaryDelete 删除日记
func DiaryDelete(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}
