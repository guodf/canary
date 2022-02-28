package dtos

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode int

type ResultDto struct {
	Code  ErrorCode   `json:"code"`
	Msg   string      `json:"msg,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Total int64       `json:"total,omitempty"`
}

const (
	Success   ErrorCode = 0             //success
	Unknow    ErrorCode = 99999         // unknow
	ArgsError ErrorCode = iota + 100000 // 参数错误
)

var errorsMap = map[ErrorCode]string{
	Success:   "成功",
	Unknow:    "未知错误",
	ArgsError: "参数错误",
}

// RegisterCode 扩展、覆盖错误信息
func RegisterCode(errMap map[ErrorCode]string) {
	for k, v := range errMap {
		errorsMap[k] = v
	}
}

// FailedResult 返回值
func FailedResult(c *gin.Context, code ErrorCode) {
	FailedResultWithMsg(c, code, errorsMap[code])
}

// FailedResultWithMsg
func FailedResultWithMsg(c *gin.Context, code ErrorCode, msg string) {
	c.AbortWithStatusJSON(http.StatusOK, &ResultDto{
		Code: code,
		Msg:  msg,
	})
}

// OkResult 返回成功数据
func Ok(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, ResultDto{
		Code: Success,
		Msg:  errorsMap[Success],
		Data: data,
	})
}

// OkPagedResult 返回翻页数据
func OkPaged(c *gin.Context, data interface{}, total int64) {
	c.AbortWithStatusJSON(http.StatusOK, ResultDto{
		Code:  Success,
		Msg:   errorsMap[Success],
		Data:  data,
		Total: total,
	})
}
