package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const maxUploadSizeBytes int64 = 5 << 20 // 5MB

// UploadCreate 创建上传
func UploadCreate(ctx context.Context, c *app.RequestContext) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"error": "未登录",
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "上传内容不合法",
		})
		return
	}
	files := form.File["file"]
	if len(files) == 0 || files[0] == nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "缺少文件",
		})
		return
	}
	fh := files[0]
	if fh.Size <= 0 {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "文件为空",
		})
		return
	}
	if fh.Size > maxUploadSizeBytes {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "文件过大",
		})
		return
	}

	in, err := fh.Open()
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "读取文件失败",
		})
		return
	}
	defer in.Close()

	head := make([]byte, 512)
	n, _ := io.ReadFull(in, head)
	head = head[:n]
	mimeType := http.DetectContentType(head)
	ext := ""
	switch mimeType {
	case "image/png":
		ext = ".png"
	case "image/jpeg":
		ext = ".jpg"
	case "image/webp":
		ext = ".webp"
	case "image/gif":
		ext = ".gif"
	default:
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "不支持的文件类型",
		})
		return
	}

	rnd := make([]byte, 16)
	if _, err := rand.Read(rnd); err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "生成文件名失败",
		})
		return
	}
	name := hex.EncodeToString(rnd) + ext
	dateDir := time.Now().Format("20060101")

	root := filepath.Join("data", "uploads")
	relDir := filepath.Join(strconv.FormatUint(uint64(userID), 10), dateDir)
	absDir := filepath.Join(root, relDir)
	if err := os.MkdirAll(absDir, 0o755); err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "创建目录失败",
		})
		return
	}

	absPath := filepath.Join(absDir, name)
	out, err := os.OpenFile(absPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "保存文件失败",
		})
		return
	}
	defer out.Close()

	if len(head) > 0 {
		if _, err := out.Write(head); err != nil {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"error": "保存文件失败",
			})
			return
		}
	}
	if _, err := io.Copy(out, in); err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "保存文件失败",
		})
		return
	}
	urlPath := "/uploads/" + strings.ReplaceAll(filepath.ToSlash(filepath.Join(relDir, name)), "\\", "/")
	c.JSON(consts.StatusCreated, utils.H{
		"data": utils.H{
			"url": urlPath,
		},
	})
}
