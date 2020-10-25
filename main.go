package main

import (
	"github.com/gin-gonic/gin"
	"go-web-cli/controller"
	"go-web-cli/database"
	"go-web-cli/middleware"
	swagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "go-web-cli/docs"
	)
// @title MyAPI
// @version 0.0.1
// @description 接口文档

// @contact.name 氕氘氚
// @contact.email caorcjp@gmail.com

// @host http://localhost:5800/
// @BasePath api/
func main(){
	database.InintDatabase()
	Router := gin.Default()
	Router.Use(middleware.Cors())
	Router.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))
	v1 := Router.Group("/controller")
	controller.HelloGroup(v1.Group("/Hello"))
	Router.Run(":5800")
}
