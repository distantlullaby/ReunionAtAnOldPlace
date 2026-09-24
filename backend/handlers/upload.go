package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// POST /api/upload  上传照片，返回可访问的 URL
func Upload(uploadDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
			return
		}
		if file.Size > 10<<20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件不能超过 10MB"})
			return
		}
		ext := filepath.Ext(file.Filename)
		name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		dst := filepath.Join(uploadDir, name)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"url": "/uploads/" + name})
	}
}
