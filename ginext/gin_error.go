package ginext

import (
	"fmt"

	"github.com/nie312122330/niexq-gotools/jsonext"
)

// ValidZhError 验证异常
type ValidZhError struct {
	Err   string
	ZhErr map[string][]string
}

// Error 实现Error接口
func (e *ValidZhError) Error() string {
	jsonStr, _ := jsonext.ToStr(e.ZhErr)
	return fmt.Sprintf("%s%s", e.Err, jsonStr)
}

// RunTimeError 运行时异常
type RunTimeError struct {
	Err string
}

// Error 实现Error接口
func (e *RunTimeError) Error() string {
	return e.Err
}
