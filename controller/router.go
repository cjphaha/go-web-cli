package controller

import "github.com/gin-gonic/gin"

func HelloGroup(r *gin.RouterGroup){
	r.GET("/HelloWorld",Hello)
}
