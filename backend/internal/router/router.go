package router

import (
	"github.com/acatchai/catdiary/backend/internal/handler"
	"github.com/acatchai/catdiary/backend/internal/middleware"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func Register(h *server.Hertz) {
	h.GET("/healthz", handler.Healthz)
	h.GET("/readyz", handler.Readyz)

	v1 := h.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.POST("/register", handler.AuthRegister)
	auth.POST("/login", handler.AuthLogin)
	auth.POST("/logout", middleware.RequireAuth(), handler.AuthLogout)
	auth.GET("/me", middleware.RequireAuth(), handler.AuthMe)

	users := v1.Group("/users", middleware.RequireAuth())
	users.GET("/me", handler.UserMe)
	users.PATCH("/me", handler.UserPatchMe)
	users.PATCH("/me/password", handler.UserPatchPassword)

	diaries := v1.Group("/diaries", middleware.RequireAuth())
	diaries.POST("", handler.DiaryCreate)
	diaries.GET("", handler.DiaryList)
	diaries.GET("/:id", handler.DiaryGet)
	diaries.PUT("/:id", handler.DiaryPut)
	diaries.PATCH("/:id", handler.DiaryPatch)
	diaries.DELETE("/:id", handler.DiaryDelete)

	drafts := v1.Group("/drafts", middleware.RequireAuth())
	drafts.POST("", handler.DraftCreate)
	drafts.GET("", handler.DraftList)
	drafts.GET("/:id", handler.DraftGet)
	drafts.PUT("/:id", handler.DraftPut)
	drafts.PATCH("/:id", handler.DraftPatch)
	drafts.DELETE("/:id", handler.DraftDelete)
	drafts.POST("/:id/flush", handler.DraftFlush)

	uploads := v1.Group("/uploads", middleware.RequireAuth())
	uploads.POST("", handler.UploadCreate)
}
