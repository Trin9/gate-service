package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- 配置 ---
const (
	vllmURL   = "http://localhost:8000/v1/chat/completions" // vLLM 地址
	jwtSecret = "your-secret-key"
)

func ProxyHandler(c *gin.Context) {
	// A. 读取客户端请求体
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	// (可选) 这里可以解析 bodyBytes 里的 JSON，看看 model 参数对不对，或者统计 token 数

	// B. 构建发往 vLLM 的请求
	// 重点：使用 c.Request.Context()，这样客户端断开时，vLLM 请求也会被 Cancel
	proxyReq, err := http.NewRequestWithContext(c.Request.Context(), "POST", vllmURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create request"})
		return
	}
	proxyReq.Header.Set("Content-Type", "application/json")
	proxyReq.Header.Set("Authorization", "Bearer your-vllm-api-key") // 如果 vLLM 设置了 key

	// C. 发送请求
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		// 这里处理如果是 Context Cancelled 导致的错误
		c.JSON(500, gin.H{"error": "Upstream error"})
		return
	}
	defer resp.Body.Close()

	// D. 处理响应
	// 设置流式响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// E. 核心循环：读取 vLLM 的流，实时写回 Client
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			// Log error
			break
		}

		// 这里可以直接转发，也可以做一些处理（比如记录日志）
		// line 格式通常是 "data: {...}\n\n"

		fmt.Fprintf(c.Writer, "%s", line)
		c.Writer.Flush() // 关键！必须立即刷新缓冲区，否则前端看不到打字机效果
	}
}

// 💡 修改 HealthCheckHandler 以接受 *gin.Context
func HealthCheckHandler(c *gin.Context) { // 注意：参数现在是 c *gin.Context
    // Gin 框架中，我们不再直接使用 w http.ResponseWriter 和 r *http.Request
    // 而是通过 c.Writer 和 c.Request 来访问它们，但通常不需要直接操作它们。

    // 使用 Gin 推荐的 c.String() 或 c.JSON() 方法来返回响应
    // 这样它会自动设置状态码和响应头
    c.String(http.StatusOK, "Status: OK") 
    
    // 如果想要返回 JSON:
    // c.JSON(http.StatusOK, gin.H{"status": "ok"})
    
    // log.Println("Health check accessed.")
    // 注意：Gin 默认集成了 Logger 中间件，日志记录会更自动化
}