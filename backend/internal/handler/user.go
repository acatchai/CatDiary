package handler

import (
	"context"

	"github.com/acatchai/catdiary/backend/internal/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type UserPatchMeReq struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=50"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Avatar   *string `json:"avatar" validate:"omitempty,url"`
}

type UserPatchPasswordReq struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=50"`
}

// getUserIDFromCtx 从上下文获取用户ID
func getUserIDFromCtx(c *app.RequestContext) (uint, bool) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		return 0, false
	}
	return userID, true
}

// UserMe 获取用户信息
func UserMe(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}

	user, err := service.GetMe(userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(consts.StatusNotFound, utils.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "获取用户信息失败",
		})
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"data": user,
	})
}

// UserPatchMe 更新用户信息
func UserPatchMe(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}

	var req UserPatchMeReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}

	user, err := service.UpdateMe(userID, req.Username, req.Email, req.Avatar)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			c.JSON(consts.StatusNotFound, utils.H{
				"error": err.Error(),
			})
		case service.ErrUserExists, service.ErrEmailExists, service.ErrInvalidUsername, service.ErrInvalidEmail:
			status := consts.StatusBadRequest
			if err == service.ErrUserExists || err == service.ErrEmailExists {
				status = consts.StatusConflict
			}
			c.JSON(status, utils.H{
				"error": err.Error(),
			})
		default:
			c.JSON(consts.StatusInternalServerError, utils.H{
				"error": "更新用户信息失败",
			})
		}
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"message": "更新成功",
		"data":    user,
	})
}

// UserPatchPassword 更新用户密码
func UserPatchPassword(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}

	var req UserPatchPasswordReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}

	if err := service.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		switch err {
		case service.ErrUserNotFound:
			c.JSON(consts.StatusNotFound, utils.H{
				"error": err.Error(),
			})
		case service.ErrWrongPassword:
			c.JSON(consts.StatusUnauthorized, utils.H{
				"error": "旧密码错误",
			})
		default:
			c.JSON(consts.StatusInternalServerError, utils.H{
				"error": "修改密码失败",
			})
		}
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"message": "密码修改成功",
	})
}
