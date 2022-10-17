package main

import (
	"net/http"
	"zimuzu/admin/instance"
	"zimuzu/admin/middleware"
	"zimuzu/admin/routers"

	"github.com/gin-gonic/gin"
)

func initRouter(router *gin.Engine) {
	router.Use(gin.Logger())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该路由",
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该方法",
		})
	})
	router.GET("/index", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "doc/index.html", gin.H{
			"title": "index",
		})
	})
	// 系列
	routers.RegisterSeriesRouter(router)
	routers.RegisterUserRouter(router)
	routers.RegisterResourceRouter(router)
	routers.RegisterSubGroupRouter(router)

}

func main() {

	// service1.CreateBaseTMDB(438148)
	// service1.CreateBaseSeries("tt1877830")

	r := gin.Default()
	initApp(r)
	r.Static("/img", "./static/img")
	r.Static("/resource", "./static/resource")
	r.Static("/icon", "./static/icon")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// 初始化App
func initApp(r *gin.Engine) {

	// 初始化静态文件夹
	r.Static("/static", "./app/static")
	// // 初始化中间件
	middleware.InitMiddleware(r)
	// 初始化路由
	initRouter(r)
	// 初始化实例
	instance.InitDBInstance()
	instance.InitTMDBInstance()
}
