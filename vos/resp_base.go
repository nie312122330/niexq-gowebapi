package vos

import (
	"github.com/nie312122330/niexq-gotools/dateext"
)

// EmptyObj ...
type EmptyObj struct {
}

// BaseResp 基础响应
type BaseResp struct {
	// 响应码
	Code int `json:"code"`
	// 响应消息
	Msg string `json:"msg"`
	// 列表时的数据总量
	Count int64 `json:"count"`
	// 列表时的数据总页
	PageCount int64 `json:"pageCount"`
	// 返回成功时可能包含警告信息
	Warn string `json:"warn"`
	// 服务器当前时间
	ServerTime string `json:"serverTime"`
	// 每个接口返回的数据类型不同，在接口文档中具体指定
	Data interface{} `json:"data" swaggerignore:"true" `
	// 每个接口返回的数据类型不同，在接口文档中具体指定
	ExtData interface{} `json:"extData" swaggerignore:"true" `
}

func emptyObj() *BaseResp {
	instance := new(BaseResp)
	serverTimeStr, _ := dateext.Now().Format("yyyy-MM-dd HH:mm:ss")
	instance.ServerTime = serverTimeStr
	instance.Data = new(EmptyObj)
	instance.ExtData = new(EmptyObj)
	return instance
}

// NewErrBaseResp ...
func NewErrBaseResp(msg string) *BaseResp {
	instance := emptyObj()
	instance.Code = 9999
	instance.Msg = msg
	return instance
}

// NewOkBaseResp ...
func NewOkBaseResp(data interface{}) *BaseResp {
	instance := emptyObj()
	instance.Data = data
	instance.Code = 0
	return instance
}

// NewNoBaseResp ...
func NewNoBaseResp(code int, msg string) *BaseResp {
	instance := emptyObj()
	instance.Code = code
	instance.Msg = msg
	return instance
}
