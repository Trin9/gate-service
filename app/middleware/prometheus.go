package middleware

import (
	"gate-service/app/monitor"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// --- 2. 监控中间件 ---
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // 执行后续逻辑

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// 记录指标
		monitor.RequestCount.WithLabelValues(c.Request.Method, status).Inc()
		monitor.RequestDuration.WithLabelValues(c.Request.Method).Observe(duration)
	}
}
