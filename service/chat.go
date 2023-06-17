package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-talk/common/log"
	"go-talk/common/model"
	"go-talk/utils"
	"go.uber.org/zap"
	"strconv"
)

type Chat struct{}

var ErrRoomIdentityEmpty = errors.New("房间号为空")
var ErrNoValid = errors.New("非法访问")

type ChatListResp struct {
	List []*model.MessageBasic `json:"list"`
}

func (chat *Chat) ChatList(c *gin.Context) (interface{}, error) {
	roomIdentity := c.Query("room_identity")
	if roomIdentity == "" {
		log.Logger.Info("房间号不能为空")
		return nil, ErrRoomIdentityEmpty
	}

	// 判断用户是否属于该房间
	uc := c.MustGet("user_claims").(*utils.UserClaims)
	_, err := model.GetUserRoomByUserIdentityRoomIdentity(uc.Identity, roomIdentity)
	if err != nil {
		log.Logger.Error("非法访问")
		return nil, ErrNoValid
	}

	pageIndex, _ := strconv.ParseInt(c.Query("page_index"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 32)
	skip := (pageIndex - 1) * pageSize

	// 聊天记录查询
	data, err := model.GetMessageListByRoomIdentity(roomIdentity, &pageSize, &skip)
	if err != nil {
		log.Logger.Error("MongoDB ERROR", zap.Error(err))
		return nil, err
	}

	return ChatListResp{
		List: data,
	}, nil
}
