package utils

import "github.com/gin-gonic/gin"

// 统一响应结构：{ code, msg, data }。code=0 表示成功。
type apiResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(200, apiResp{Code: 0, Msg: "ok", Data: data})
}

func Fail(c *gin.Context, httpCode int, msg string) {
	c.JSON(httpCode, apiResp{Code: httpCode, Msg: msg})
}
