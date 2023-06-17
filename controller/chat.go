package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-talk/common/log"
	"go-talk/common/model"
	res "go-talk/common/result"
	"go-talk/service"
	"go-talk/utils"
	"go.uber.org/zap"
	"time"
)

var upgrader = websocket.Upgrader{}
var wc = make(map[string]*websocket.Conn)

type MessageStruct struct {
	Message      string `json:"message"`
	RoomIdentity string `json:"room_identity"`
}

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

func WebsocketMessage(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Logger.Error("websocket error", zap.Error(err))
		res.Error(c, res.Status{
			StatusCode: res.WebSocketErrorStatus.StatusCode,
			StatusMsg:  res.WebSocketErrorStatus.StatusMsg,
		})
		return
	}
	defer conn.Close()
	uc := c.MustGet("user_claims").(*utils.UserClaims)
	wc[uc.Identity] = conn
	for {
		ms := new(MessageStruct)
		err := conn.ReadJSON(ms)
		if err != nil {
			log.Logger.Error("Read Error", zap.Error(err))
			return
		}
		// 判断用户是否属于消息体的房间
		_, err = model.GetUserRoomByUserIdentityRoomIdentity(uc.Identity, ms.RoomIdentity)
		if err != nil {
			log.Logger.Info("UserIdentity with RoomIdentity Not Exits", zap.Error(err))
			return
		}
		// 保存消息
		mb := &model.MessageBasic{
			UserIdentity: uc.Identity,
			RoomIdentity: ms.RoomIdentity,
			Data:         ms.Message,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}
		err = model.InsertOneMessageBasic(mb)
		if err != nil {
			log.Logger.Error("[DB ERROR]", zap.Error(err))
			return
		}
		// 获取在特定房间的在线用户
		userRooms, err := model.GetUserRoomByRoomIdentity(ms.RoomIdentity)
		if err != nil {
			log.Logger.Error("[DB ERROR]", zap.Error(err))
			return
		}
		for _, room := range userRooms {
			if cc, ok := wc[room.UserIdentity]; ok {
				err := cc.WriteMessage(websocket.TextMessage, []byte(ms.Message))
				if err != nil {
					log.Logger.Error("Write Message Error", zap.Error(err))
					return
				}
			}
		}
	}
}
