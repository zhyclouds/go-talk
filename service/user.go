package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-talk/common/log"
	"go-talk/common/model"
	"go-talk/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrUsernameExits = errors.New("username already exists")
var PasswordErr = errors.New("password error")
var ErrIdEmpty = errors.New("user account is empty")
var ErrIsFriend = errors.New("already friends")

type User struct{}

type UserLoginReq struct {
	Account  string `form:"account" binding:"required,min=1,max=32"`
	Password string `form:"password" binding:"required,min=6,max=32"`
}

type UserLoginResp struct {
	Token string `json:"token"`
}

type UserRegisterReq struct {
	Account  string `form:"account" binding:"required,min=1,max=32"`
	Password string `form:"password" binding:"required,min=6,max=32"`
	Email    string `form:"email"`
}

type UserRegisterResp struct {
	Token string `json:"token"`
}

type UserAddFriendResp struct {
	Msg string `json:"msg"`
}

// Register 用户注册
func (u *User) Register(c *gin.Context) (interface{}, error) {
	var req UserRegisterReq
	err := c.ShouldBindWith(&req, binding.Form)
	if err != nil {
		log.Logger.Error("validate err", zap.Error(err))
		return nil, err
	}

	// 检查是否已经注册
	cnt, err := model.GetUserBasicCountByAccount(req.Account)
	if err != nil {
		log.Logger.Error("MongoDB happen error", zap.Error(err))
		return nil, err
	}
	if cnt > 0 {
		return nil, ErrUsernameExits
	}

	// 加密
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := model.User{
		Identity: utils.GetUUID(),
		Account:  req.Account,
		Password: string(hash),
		Email:    req.Email,
	}

	err = model.InsertOneUserBasic(&user)
	if err != nil {
		log.Logger.Error("MongoDB happen error", zap.Error(err))
		return nil, err
	}

	token, err := utils.GenerateToken(user.Identity, user.Email)
	if err != nil {
		log.Logger.Error("GenerateToken error", zap.Error(err))
		return nil, err
	}

	return UserRegisterResp{
		Token: token,
	}, nil
}

// Login 用户登录
func (u *User) Login(c *gin.Context) (interface{}, error) {
	var req UserLoginReq

	// 解析参数
	err := c.ShouldBindWith(&req, binding.Form)
	if err != nil {
		log.Logger.Error("parse json error")
		return nil, err
	}

	ub, err := model.GetUserBasicByAccount(req.Account)
	if err != nil {
		log.Logger.Error("MongoDB happen error", zap.Error(err))
		return nil, err
	}

	// 检查密码
	err = bcrypt.CompareHashAndPassword([]byte(ub.Password), []byte(req.Password))
	if err != nil {
		log.Logger.Error("password error", zap.Any("user", ub))
		return nil, PasswordErr
	}

	// 生成token
	token, err := utils.GenerateToken(ub.Identity, ub.Email)
	if err != nil {
		log.Logger.Error("GenerateToken error", zap.Error(err))
		return nil, err
	}

	return UserLoginResp{
		Token: token,
	}, nil
}

// AddFriend 添加好友
func (u *User) AddFriend(c *gin.Context) (interface{}, error) {
	account := c.PostForm("account")
	if account == "" {
		return nil, ErrIdEmpty
	}

	friend, err := model.GetUserBasicByAccount(account)
	if err != nil {
		log.Logger.Error("[DB ERROR]", zap.Error(err))
		return nil, err
	}

	user := c.MustGet("user_claims").(*utils.UserClaims)

	if model.JudgeUserIsFriend(friend.Identity, user.Identity) {
		return nil, ErrIsFriend
	}
	// 保存房间记录
	rb := &model.RoomBasic{
		Identity:     utils.GetUUID(),
		UserIdentity: user.Identity,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	if err := model.InsertOneRoomBasic(rb); err != nil {
		log.Logger.Error("[DB ERROR]", zap.Error(err))
		return nil, err
	}
	// 保存用户与房间的关联记录
	ur := &model.UserRoom{
		UserIdentity: user.Identity,
		RoomIdentity: rb.Identity,
		RoomType:     1,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	if err := model.InsertOneUserRoom(ur); err != nil {
		log.Logger.Error("[DB ERROR]", zap.Error(err))
		return nil, err
	}
	ur = &model.UserRoom{
		UserIdentity: friend.Identity,
		RoomIdentity: rb.Identity,
		RoomType:     1,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	if err := model.InsertOneUserRoom(ur); err != nil {
		log.Logger.Error("[DB ERROR]", zap.Error(err))
		return nil, err
	}

	return UserAddFriendResp{
		Msg: "添加成功",
	}, nil
}
