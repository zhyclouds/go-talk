package main

import (
	"github.com/gin-gonic/gin"
	ctrl "go-talk/controller"
	"go-talk/middlewares"
)

func handle(r *gin.Engine) {
	// 测试接口
	r.GET("/ping", ctrl.Ping)

	basic := r.Group("/gotalk")

	// 用户登录
	basic.POST("/register", ctrl.Register)
	// 新用户注册
	basic.POST("/login", ctrl.Login)

	// 用户相关接口
	userGroup := basic.Group("/user", middlewares.AuthCheck())
	{
		// 添加好友
		userGroup.POST("/add/friend", ctrl.AddFriend)

		// 删除好友
		userGroup.POST("/delete/friend", ctrl.DeleteFriend)
	}

	// 聊天相关接口
	chatGroup := basic.Group("/chat", middlewares.AuthCheck())
	{
		// 获取用户聊天消息列表
		chatGroup.GET("/chatlist", ctrl.ChatList)

		// 发送消息
		chatGroup.GET("/websocket/message", ctrl.WebsocketMessage)
	}
}
