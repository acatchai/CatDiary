package handler

import (
	"context"

	"github.com/acatchai/catdiary/backend/internal/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type RegisterReq struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6,max=50"`
	Email    string `json:"email" validate:"omitempty,email"`
}

// AuthRegister 注册用户
func AuthRegister(ctx context.Context, c *app.RequestContext) {
	var req RegisterReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}
	user, err := service.Register(req.Username, req.Password, req.Email)
	if err != nil {
		if err == service.ErrUserExists {
			c.JSON(consts.StatusConflict, utils.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"error": "注册失败",
			})
		}
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"message": "注册成功",
		"data": utils.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthLogin 用户登录
func AuthLogin(ctx context.Context, c *app.RequestContext) {
	var req LoginReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}

	token, user, err := service.Login(req.Username, req.Password)
	if err != nil {
		if err == service.ErrUserNotFound || err == service.ErrWrongPassword {
			c.JSON(consts.StatusUnauthorized, utils.H{
				"error": "用户名或密码错误",
			})
		} else {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"error": "登录失败",
			})
		}
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"message": "登录成功",
		"token":   token,
		"user": utils.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// AuthLogout 用户退出登录
func AuthLogout(ctx context.Context, c *app.RequestContext) {
	// JWT 无状态， 前端删除token即可。如果需要强退出，可以引入Redis黑名单，这里目前返回成功。
	c.JSON(consts.StatusOK, utils.H{
		"message": "退出成功",
	})
}

// AuthMe 获取用户信息
func AuthMe(ctx context.Context, c *app.RequestContext) {
	// 从 context 中获取中间件存入的 user_id
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}
	userID := userIDAny.(uint)
	user, err := service.GetMe(userID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "获取用户信息失败",
		})
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"data": user,
	})
}
