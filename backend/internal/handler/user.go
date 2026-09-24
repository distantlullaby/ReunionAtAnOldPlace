package handler

import (
	"net/http"

	"memorylink/internal/service"
	"memorylink/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Profile(c *gin.Context) {
	data, err := h.svc.Profile(utils.CurrentUserID(c))
	if err != nil {
		utils.Fail(c, http.StatusNotFound, "用户不存在")
		return
	}
	utils.OK(c, data)
}

func (h *UserHandler) Update(c *gin.Context) {
	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.UpdateProfile(utils.CurrentUserID(c), req.Nickname, req.Avatar); err != nil {
		utils.Fail(c, http.StatusBadRequest, "更新失败")
		return
	}
	utils.OK(c, gin.H{"ok": true})
}
