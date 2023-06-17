package controller

import (
	"github.com/gin-gonic/gin"
	res "go-talk/common/result"
	"go-talk/service"
)

func ChatList(c *gin.Context) {
	var chat service.Chat
	chatList, err := chat.ChatList(c)
	if err != nil {
		if err == service.ErrRoomIdentityEmpty {
			res.Error(c, res.Status{
				StatusCode: res.RoomIdentityEmptyStatus.StatusCode,
				StatusMsg:  res.RoomIdentityEmptyStatus.StatusMsg,
			})
			return
		} else if err == service.ErrNoValid {
			res.Error(c, res.Status{
				StatusCode: res.NoValidErrorStatus.StatusCode,
				StatusMsg:  res.NoValidErrorStatus.StatusMsg,
			})
			return
		} else {
			res.Error(c, res.Status{
				StatusCode: res.ChatListErrorStatus.StatusCode,
				StatusMsg:  res.ChatListErrorStatus.StatusMsg,
			})
			return
		}
	}
	data := chatList.(service.ChatListResp)
	res.Success(c, res.R{
		"data": data,
	})
}
