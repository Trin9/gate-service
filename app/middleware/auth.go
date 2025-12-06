package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- 2. JWT 中间件 ---
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token"})
			return
		}
		// ... 这里省略具体的 JWT 解析代码 ...
		// 假设解析成功，拿到 UserID
		c.Set("userID", "user_123")
		c.Next()
	}
}
