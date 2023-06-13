package main

import (
	"github.com/gin-gonic/gin"
	ctrl "go-talk/controller"
)

func handle(r *gin.Engine) {
	// 测试接口
	r.GET("/ping", ctrl.Ping)
}
