package main

import "github.com/gin-gonic/gin"

// failed to run http server: conflict path
func main() {
	r := gin.Default()
	r.GET("/ping/:id", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/ping/hoge", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong",
		})
	})
	if err := r.Run(); err != nil {
		panic(err)
	}
}
