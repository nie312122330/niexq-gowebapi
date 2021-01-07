package tools

// LimitParams ...
func LimitParams(pageNo int, pageSize int) (limitStart int, limitSize int) {
	if pageNo == 0 {
		pageNo = 1
	}
	return (pageNo - 1) * pageSize, pageSize
}
