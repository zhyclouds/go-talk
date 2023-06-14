package main

import (
	"github.com/gin-gonic/gin"
	ctrl "go-talk/controller"
)

func handle(r *gin.Engine) {
	// 测试接口
	r.GET("/ping", ctrl.Ping)

	basic := r.Group("/gotalk")

	// 用户相关接口
	userGroup := basic.Group("/user")
	{
		// 用户登录
		userGroup.POST("/register", ctrl.Register)
		// 新用户注册
		userGroup.POST("/login", ctrl.Login)
		// 添加好友
		userGroup.POST("/add/friend", ctrl.AddFriend)
	}

}
