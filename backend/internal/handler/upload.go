package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"memorylink/internal/utils"

	"github.com/gin-gonic/gin"
)

// Upload 接收单张图片，存入 uploads/yyyy-mm/ 并返回可访问 URL。
func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "请选择要上传的图片")
		return
	}

	// 限制 10MB。
	if file.Size > 10<<20 {
		utils.Fail(c, http.StatusBadRequest, "图片不能超过 10MB")
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	name := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), strconv.FormatInt(time.Now().Unix(), 36), ext)
	rel := filepath.Join("uploads", time.Now().Format("2006-01"), name)
	// 统一为 URL 用的正斜杠。
	savePath := filepath.FromSlash(rel)
	urlPath := "/" + filepath.ToSlash(rel)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.Fail(c, http.StatusInternalServerError, "图片保存失败")
		return
	}
	utils.OK(c, gin.H{"url": urlPath})
}
