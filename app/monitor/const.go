package monitor

import "github.com/prometheus/client_golang/prometheus"

// --- 1. 定义指标 ---
var (
	// 请求总数计数器 (Counter)
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "status"}, // 标签：方法、状态码
	)

	// 请求耗时直方图 (Histogram)
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets, // 默认的分桶: .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10
		},
		[]string{"method"},
	)
)

func init() {
	// 注册指标
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(RequestDuration)
}
