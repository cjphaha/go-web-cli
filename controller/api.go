package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)
// @Summary HelloWorld
// @Description 这里是详细描述
// @Tags 获取信息
// @Accept application/json
// @Produce application/json
// @Success 200
// @Router /HelloWorld [get]
func Hello(c *gin.Context){
	c.JSON(http.StatusOK,gin.H{
		"message":"hello",
	})
}
