package voext

// BaseReq ...
type BaseReq struct {
}

// BaseReqPage ...
type BaseReqPage struct {
	BaseReq
	// 页码，从1开始
	PageNo int `json:"pageNo" binding:"required,gte=1"`
	// 页码，每页大小，必须大于1
	PageSize int `json:"pageSize" binding:"required,gte=1"`
}
