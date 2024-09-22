package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 使用cookie和token进行验证
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取
		tokenString := ctx.GetHeader("Authorization")

		// 验证
		if tokenString == "" || strings.HasPrefix(tokenString, "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}
	}
}

// 从会话获取用户 ID
func getUserIDFromSession(c *gin.Context) uint {
	session := sessions.Default(c)
	// id := c.PostForm("login")
	// fmt.Println("Form id ", id)
	userID, _ := session.Get("login").(uint)
	return userID
}

// 通过Session来验证身份，并保存为userID
func SesssionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserIDFromSession(c)

		fmt.Println("Get session ID :", userID)
		if userID == 0 { // 不存在
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		c.Set("userID", userID)
		c.Next()
		c.JSON(200, gin.H{"code": 200, "msg": "session验证成功"})
	}
}
