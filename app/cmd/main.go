package main

import (
	"gate-service/app/handler"
	"gate-service/app/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := gin.Default() // 自带 Logger 和 Recovery 中间件

	// API 路由组
	api := r.Group("/v1")
	api.Use(middleware.AuthMiddleware())       // 挂载鉴权
	api.Use(middleware.RateLimitMiddleware())  // 挂载限流
	api.Use(middleware.PrometheusMiddleware()) // 挂载限流

	// --- 3. 暴露 /metrics 接口 ---
	// Prometheus 会访问这个接口来“刮取”数据
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api.POST("/chat/completions", handler.ProxyHandler)
	api.POST("/health", handler.HealthCheckHandler)

	r.Run(":8080")
}
