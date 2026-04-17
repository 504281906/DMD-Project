package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 中间件 - 请求日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		
		c.Next()
		
		latency := time.Since(start)
		status := c.Writer.Status()
		
		if status >= 400 {
			println("[WARN]", path, status, latency)
		} else {
			println("[INFO]", path, status, latency)
		}
	}
}

// Recovery 中间件 - 错误恢复
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}

// CORS 中间件 - 跨域支持
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}