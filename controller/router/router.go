/* 制定路由分组规则 */

package router

import (
	"fmt"
	ctr "ginEssential/controller/controller"
	mid "ginEssential/controller/middleware"
	"net/http"

	"io"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

/* 给gin制定路由规则 */
func Router() *gin.Engine {
	r := gin.Default()

	// 设置日志格式
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	// gin.DisableConsoleColor()

	// 配置日志文件
	logFile, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	// 配置session
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// 用来测试的路由
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World.")
	})

	//user组下面的路由
	user := r.Group("/user")
	user.GET("", func(c *gin.Context) { // 测试User组路由用
		c.String(http.StatusOK, "This is User group.")
	})

	{
		user.POST("/register", ctr.Register)
		user.POST("/login", ctr.Login)

		user.GET("/cookie", ctr.GetCookie)                  //, ctr.AddCookie
		user.POST("/session", mid.SesssionAuthMiddleware()) //, controller.AddCookie
	}

	music := r.Group("/music")
	{
		music.POST("/vote", ctr.UserVoteMusic)               // 要验证身份 mid.AuthMiddleware,
		music.POST("/cancelVoting", ctr.UserCancelVoteMusic) // 要验证身份 mid.AuthMiddleware,
		music.GET("/rank", ctr.MusicRank)                    // 获取音乐排行榜，不需要验证身份
	}

	return r
}
