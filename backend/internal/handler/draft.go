package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// DraftPutDiary 更新草稿日记
func DraftPutDiary(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}

// DraftGetDiary 获取草稿日记
func DraftGetDiary(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}

// DraftDeleteDiary 删除草稿日记
func DraftDeleteDiary(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusNotImplemented, "TODO")
}
