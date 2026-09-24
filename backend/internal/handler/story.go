package handler

import (
	"net/http"
	"strconv"

	"memorylink/internal/model"
	"memorylink/internal/service"
	"memorylink/internal/utils"

	"github.com/gin-gonic/gin"
)

type StoryHandler struct {
	svc *service.StoryService
}

func NewStoryHandler(svc *service.StoryService) *StoryHandler {
	return &StoryHandler{svc: svc}
}

type createStoryReq struct {
	Title      string             `json:"title" binding:"required"`
	Location   string             `json:"location"`
	MemoryText string             `json:"memoryText"`
	OldPhotos  model.StringSlice  `json:"oldPhotos"`
	Reward     int                `json:"reward"`
}

func (h *StoryHandler) Create(c *gin.Context) {
	var req createStoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "标题不能为空")
		return
	}
	story, err := h.svc.CreateStory(utils.CurrentUserID(c), service.CreateStoryInput{
		Title:      req.Title,
		Location:   req.Location,
		MemoryText: req.MemoryText,
		OldPhotos:  req.OldPhotos,
		Reward:     req.Reward,
	})
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.OK(c, story)
}

func (h *StoryHandler) AppendReward(c *gin.Context) {
	storyID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Add int `json:"add"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	story, err := h.svc.AppendReward(utils.CurrentUserID(c), uint(storyID), req.Add)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.OK(c, story)
}

func (h *StoryHandler) Respond(c *gin.Context) {
	storyID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		NewPhoto string `json:"newPhoto" binding:"required"`
		Message  string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "请上传当下现场照片")
		return
	}
	resp, err := h.svc.AddResponse(utils.CurrentUserID(c), uint(storyID), req.NewPhoto, req.Message)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.OK(c, resp)
}

func (h *StoryHandler) Accept(c *gin.Context) {
	storyID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	responseID, _ := strconv.ParseUint(c.Param("rid"), 10, 64)
	story, err := h.svc.AcceptResponse(utils.CurrentUserID(c), uint(storyID), uint(responseID))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.OK(c, story)
}

func (h *StoryHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	status := c.Query("status")
	list, total, err := h.svc.ListStories(page, size, status)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, "加载失败")
		return
	}
	utils.OK(c, gin.H{"list": list, "total": total})
}

func (h *StoryHandler) Detail(c *gin.Context) {
	storyID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	detail, err := h.svc.GetStory(uint(storyID))
	if err != nil {
		utils.Fail(c, http.StatusNotFound, err.Error())
		return
	}
	utils.OK(c, detail)
}
