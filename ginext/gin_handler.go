package ginext

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nie312122330/niexq-gotools/stringext"
	"github.com/nie312122330/niexq-gowebapi/voext"
)

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return slogHandlerFunc
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery() gin.HandlerFunc {
	return recoveryHandlerFunc
}

func slogHandlerFunc(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	c.Next()
	cost := time.Since(start).Milliseconds()
	agentStr := stringext.CutString(c.Request.UserAgent(), 32)
	logStr := fmt.Sprintf("%s\t%s\t%d\t%s\t%dms\t%s\t%s", path, c.Request.Method, c.Writer.Status(), c.ClientIP(), cost, query, agentStr)
	slog.Info(logStr)
}

func recoveryHandlerFunc(c *gin.Context) {
	defer recoveryErrorWork(c)
	c.Next()
}

func recoveryErrorWork(c *gin.Context) {
	reqPath := c.Request.URL.Path
	if err := recover(); err != nil {
		//检查连接是否已断开
		var brokenPipe bool
		if ne, ok := err.(*net.OpError); ok {
			if se, ok := ne.Err.(*os.SyscallError); ok {
				if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
					brokenPipe = true
				}
			}
		}
		httpRequest, _ := httputil.DumpRequest(c.Request, false)
		if brokenPipe {
			//发生了异常，但是连接已断开
			slog.Error(fmt.Sprintf("%s\t异常:%v\n%s\n%s", reqPath, err, string(httpRequest), "连接已断开"))
			c.Error(err.(error)) // nolint: errcheck
			c.Abort()
			return
		}

		//这里把错误输出出去
		if vze, ok := err.(ValidZhError); ok {
			result := voext.NewErrBaseResp(fmt.Sprintf("%v", vze.Err))
			result.ExtData = vze.ZhErr
			c.JSON(http.StatusOK, &result)
		} else if vze, ok := err.(RunTimeError); ok {
			result := voext.NewErrBaseResp(fmt.Sprintf("%v", vze.Err))
			c.JSON(http.StatusOK, &result)
		} else {
			slog.Error(fmt.Sprintf("%s\t异常:%v\n%s\n%s", reqPath, err, string(httpRequest), debug.Stack()))
			result := voext.NewErrBaseResp(fmt.Sprintf("%v", err))
			c.JSON(http.StatusOK, &result)
		}
		c.Abort()
	}
}
