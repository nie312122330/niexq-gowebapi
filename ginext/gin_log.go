package ginext

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nie312122330/niexq-gotools/logext"
	"github.com/nie312122330/niexq-gotools/stringext"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger = logext.DefaultLogger("gin_log")
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start).Milliseconds()
		agentStr := stringext.CutString(c.Request.UserAgent(), 32)
		logStr := fmt.Sprintf("%s\t%s\t%d\t%s\t%dms\t%s\t%s", path, c.Request.Method, c.Writer.Status(), c.ClientIP(), cost, query, agentStr)
		logger.Info(logStr)
	}
}
