package main

import (
	// "fmt"
	// "ginEssential/controller/router"
	// "ginEssential/ocean_learn/controller"

	// "ginEssential/tutu"

	"ginEssential/controller/router"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

/* 基本的使用gin的模版 */
func gin_demo() {
	r := gin.Default()

	r.Use(
		cors.Default(),
	)

	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, Geektutu")
		// c.JSON(200)
	})
	r.Run(":15480") // listen and serve on 0.0.0.0:8080
}

func main() {
	r := router.Router()
	r.Run("127.0.0.1:8080")

}
