package handler

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type DraftCreateReq struct {
	Title    string `json:"title" validate:"required,min=1,max=100"`
	Content  string `json:"content" validate:"required"`
	Mood     string `json:"mood" validate:"omitempty,max=20"`
	Weather  string `json:"weather" validate:"omitempty,max=20"`
	Location string `json:"location" validate:"omitempty,max=100"`
}

type DraftPutReq struct {
	ExpectedVersion *uint64 `json:"expected_version" validate:"omitempty"`
	Title           string  `json:"title" validate:"required,min=1,max=100"`
	Content         string  `json:"content" validate:"required"`
	Mood            string  `json:"mood" validate:"omitempty,max=20"`
	Weather         string  `json:"weather" validate:"omitempty,max=20"`
	Location        string  `json:"location" validate:"omitempty,max=100"`
}

type DraftPatchReq struct {
	ExpectedVersion *uint64 `json:"expected_version" validate:"omitempty"`
	Title           *string `json:"title" validate:"omitempty,min=1,max=100"`
	Content         *string `json:"content" validate:"omitempty,min=1"`
	Mood            *string `json:"mood" validate:"omitempty,max=20"`
	Weather         *string `json:"weather" validate:"omitempty,max=20"`
	Location        *string `json:"location" validate:"omitempty,max=100"`
}

func parseUint64Param(c *app.RequestContext, key string) (uint64, error) {
	v := strings.TrimSpace(c.Param(key))
	u64, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, strconv.ErrSyntax
	}
	return u64, nil
}

// DraftCreate 创建草稿日记
func DraftCreate(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}
	var req DraftCreateReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}
	draft, err := service.CreateDraftDiary(userID, req.Title, req.Content, req.Mood, req.Weather, req.Location)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "创建草稿失败",
		})
		return
	}
	c.JSON(consts.StatusCreated, utils.H{
		"data": draft,
	})
}

// DraftList 获取草稿日记列表
func DraftList(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}

	page := clampInt(parseIntQueryDefault(c, "page", 1), 1, 1000000)
	pageSize := clampInt(parseIntQueryDefault(c, "page_size", 20), 1, 100)
	items, total, err := service.ListDraftDiaries(userID, page, pageSize)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "获取草稿列表失败",
		})
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"data": utils.H{
			"items":     items,
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// DraftGet 获取草稿日记详情
func DraftGet(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}
	id, err := parseUint64Param(c, "id")
	if err != nil || id == 0 {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "id 参数不合法",
		})
		return
	}
	draft, err := service.GetDraftDiary(userID, id)
	if err != nil {
		if err == service.ErrDraftNotFound {
			c.JSON(consts.StatusNotFound, utils.H{
				"error": "草稿不存在",
			})
			return
		}
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "获取草稿失败",
		})
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"data": draft,
	})

}

// DraftPut 更新草稿日记
func DraftPut(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}
	id, err := parseUint64Param(c, "id")
	if err != nil || id == 0 {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "id 参数不合法",
		})
		return
	}
	var req DraftPutReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}

	draft, err := service.PutDraftDiary(userID, id, req.ExpectedVersion, req.Title, req.Content, req.Mood, req.Weather, req.Location)
	if err != nil {
		if err == service.ErrDraftNotFound {
			c.JSON(consts.StatusNotFound, utils.H{
				"error": "草稿不存在",
			})
			return
		}
		var ce *service.DraftConflictError
		if errors.As(err, &ce) {
			c.JSON(consts.StatusConflict, utils.H{
				"error":           "草稿已被更新",
				"current_version": ce.CurrentVersion,
			})
			return
		}
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "更新草稿失败",
		})
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"data": draft,
	})
}

// DraftPatch 更新草稿日记
func DraftPatch(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}
	id, err := parseUint64Param(c, "id")
	if err != nil || id == 0 {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "id 参数不合法",
		})
		return
	}

	var req DraftPatchReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}

	draft, err := service.PatchDraftDiary(userID, id, req.ExpectedVersion, req.Title, req.Content, req.Mood, req.Weather, req.Location)
	if err != nil {
		if err == service.ErrDraftNotFound {
			c.JSON(consts.StatusNotFound, utils.H{
				"error": "草稿不存在",
			})
			return
		}
		if err == service.ErrNoDraftUpdates {
			c.JSON(consts.StatusBadRequest, utils.H{
				"error": "没有可更新的字段",
			})
			return
		}
		var ce *service.DraftConflictError
		if errors.As(err, &ce) {
			c.JSON(consts.StatusConflict, utils.H{
				"error":           "草稿已被更新",
				"current_version": ce.CurrentVersion,
			})
			return
		}
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "更新草稿失败",
		})
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"data": draft,
	})
}

// DraftDelete 删除草稿日记
func DraftDelete(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}
	id, err := parseUint64Param(c, "id")
	if err != nil || id == 0 {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "id 参数不合法",
		})
		return
	}
	if err := service.DeleteDraftDiary(userID, id); err != nil {
		if err == service.ErrDraftNotFound {
			c.JSON(consts.StatusNotFound, utils.H{
				"error": "草稿不存在",
			})
			return
		}
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "删除草稿失败",
		})
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"message": "删除成功",
	})
}

func DraftFlush(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}
	id, err := parseUint64Param(c, "id")
	if err != nil || id == 0 {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "id 参数不合法",
		})
		return
	}
	if err := service.FlushDraftDiary(userID, id); err != nil {
		if err == service.ErrDraftNotFound {
			c.JSON(consts.StatusNotFound, utils.H{
				"error": "草稿不存在",
			})
			return
		}
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "触发落库失败",
		})
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"message": "触发落库成功",
	})
}
