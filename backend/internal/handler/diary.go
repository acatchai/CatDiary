package handler

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/acatchai/catdiary/backend/internal/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type DiaryCreateReq struct {
	Title      string  `json:"title" validate:"required,min=1,max=100"`
	Content    string  `json:"content" validate:"required"`
	Mood       string  `json:"mood" validate:"omitempty,max=20"`
	Weather    string  `json:"weather" validate:"omitempty,max=20"`
	Location   string  `json:"location" validate:"omitempty,max=100"`
	OccurredAt *string `json:"occurred_at" validate:"omitempty"`
}

type DiaryPutReq struct {
	Title      string  `json:"title" validate:"required,min=1,max=100"`
	Content    string  `json:"content" validate:"required"`
	Mood       string  `json:"mood" validate:"omitempty,max=20"`
	Weather    string  `json:"weather" validate:"omitempty,max=20"`
	Location   string  `json:"location" validate:"omitempty,max=100"`
	OccurredAt *string `json:"occurred_at" validate:"omitempty"`
}

type DiaryPatchReq struct {
	Title      *string `json:"title" validate:"omitempty,min=1,max=100"`
	Content    *string `json:"content" validate:"omitempty"`
	Mood       *string `json:"mood" validate:"omitempty,max=20"`
	Weather    *string `json:"weather" validate:"omitempty,max=20"`
	Location   *string `json:"location" validate:"omitempty,max=100"`
	OccurredAt *string `json:"occurred_at" validate:"omitempty"`
}

func parseTimeRFC3339(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, strings.TrimSpace(s))
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
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}

	var req DiaryCreateReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}

	var occurredAt *time.Time
	if req.OccurredAt != nil {
		t, err := parseTimeRFC3339(*req.OccurredAt)
		if err != nil {
			c.JSON(consts.StatusBadRequest, utils.H{
				"error": "occurred_at 参数不合法",
			})
			return
		}
		occurredAt = &t
	}

	diary, err := service.CreateDiary(userID, occurredAt, req.Title, req.Content, req.Mood, req.Weather, req.Location)
	if err != nil {
		mapDiaryErr(c, err)
		return
	}
	c.JSON(consts.StatusCreated, utils.H{
		"data": diary,
	})
}

// DiaryList 获取日记列表
func DiaryList(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}

	page := clampInt(parseIntQueryDefault(c, "page", 1), 1, 1000000)
	pageSize := clampInt(parseIntQueryDefault(c, "page_size", 20), 1, 100)

	items, total, err := service.ListDiaries(userID, page, pageSize)
	if err != nil {
		mapDiaryErr(c, err)
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"items":     items,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// DiaryGet 获取日记详情
func DiaryGet(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}

	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "id 参数不合法",
		})
		return
	}

	diary, err := service.GetDiary(userID, id)
	if err != nil {
		mapDiaryErr(c, err)
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"data": diary,
	})
}

// DiaryPut 更新日记
func DiaryPut(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}

	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "id 参数不合法",
		})
		return
	}

	var req DiaryPutReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}

	var occurredAt *time.Time
	if req.OccurredAt != nil {
		t, err := parseTimeRFC3339(*req.OccurredAt)
		if err != nil {
			c.JSON(consts.StatusBadRequest, utils.H{
				"error": "occurred_at 参数不合法",
			})
			return
		}
		occurredAt = &t
	}

	diary, err := service.PutDiary(userID, id, occurredAt, req.Title, req.Content, req.Mood, req.Weather, req.Location)
	if err != nil {
		mapDiaryErr(c, err)
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"data": diary,
	})
}

// DiaryPatch 更新日记内容
func DiaryPatch(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "id 参数不合法",
		})
		return
	}
	var req DiaryPatchReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}

	var occurredAt *time.Time
	if req.OccurredAt != nil {
		t, err := parseTimeRFC3339(*req.OccurredAt)
		if err != nil {
			c.JSON(consts.StatusBadRequest, utils.H{
				"error": "occurred_at 参数不合法",
			})
			return
		}
		occurredAt = &t
	}

	diary, err := service.PatchDiary(userID, id, occurredAt, req.Title, req.Content, req.Mood, req.Weather, req.Location)
	if err != nil {
		mapDiaryErr(c, err)
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"data": diary,
	})
}

// DiaryDelete 删除日记
func DiaryDelete(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "id 参数不合法",
		})
		return
	}
	if err := service.DeleteDiary(userID, id); err != nil {
		mapDiaryErr(c, err)
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"message": "删除成功",
	})
}
