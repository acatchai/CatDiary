package handler

import (
	"context"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type DiaryCreateReq struct {
	Title    string `json:"title" validate:"required,min=1,max=100"`
	Content  string `json:"content" validate:"required"`
	Mood     string `json:"mood" validate:"omitempty,max=20"`
	Weather  string `json:"weather" validate:"omitempty,max=20"`
	Location string `json:"location" validate:"omitempty,max=100"`
}

type DiaryPutReq struct {
	Title    string `json:"title" validate:"required,min=1,max=100"`
	Content  string `json:"content" validate:"required"`
	Mood     string `json:"mood" validate:"omitempty,max=20"`
	Weather  string `json:"weather" validate:"omitempty,max=20"`
	Location string `json:"location" validate:"omitempty,max=100"`
}

type DiaryPatchReq struct {
	Title    string `json:"title" validate:"omitempty,min=1,max=100"`
	Content  string `json:"content" validate:"omitempty"`
	Mood     string `json:"mood" validate:"omitempty,max=20"`
	Weather  string `json:"weather" validate:"omitempty,max=20"`
	Location string `json:"location" validate:"omitempty,max=100"`
}

// parseUintParam 解析 uint 类型的参数
func parseUintParam(c *app.RequestContext, key string) (uint, error) {
	v := strings.TrimSpace(c.Param(key))
	u64, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, err
	}
	if u64 == 0 {
		return 0, strconv.ErrSyntax
	}
	return uint(u64), nil
}

// parseIntQueryDefault 解析 int 类型的查询参数，默认值
func parseIntQueryDefault(c *app.RequestContext, key string, def int) int {
	v := strings.TrimSpace(c.Query(key))
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func mapDiaryErr(c *app.RequestContext, err error) {
	switch err {
	case service.ErrDiaryNotFound:
		c.JSON(consts.StatusNotFound, utils.H{
			"error": "日记不存在",
		})
	case service.ErrNoDiaryUpdates:
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "没有可更新的字段",
		})
	default:
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "服务器内部错误",
		})
	}
}

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
