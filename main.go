package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func main() {
	router.GET("/", func(ctx *gin.Context) {

		ctx.JSON(http.StatusOK, gin.H{
			"code":    "0",
			"message": "success",
		})
	})

	router.GET("/test", Test)

	router.Run(":5501")
}

func Test(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	ctx.Request.URL.Path = "/"
	router.HandleContext(ctx)
}
