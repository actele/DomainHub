package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware 检查当前用户是否具有管理员角色
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("user_role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}
