package main

import (
	_ "net/http"
	_ "zimuzu/admin/routers"

	_ "github.com/gin-gonic/gin"
)

// func initRouter(router *gin.Engine) {
// 	router.Use(gin.Logger())

// 	router.NoRoute(func(c *gin.Context) {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"code": 404,
// 			"msg":  "找不到该路由",
// 		})
// 	})

// 	router.NoMethod(func(c *gin.Context) {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"code": 404,
// 			"msg":  "找不到该方法",
// 		})
// 	})

// 	// 系列
// 	routers.RegisterSeriesRouter(router)
// 	routers.RegisterUserRouter(router)
// 	routers.RegisterResourceRouter(router)
// }
