package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-talk/common/db"
	"go-talk/common/log"
	"go-talk/common/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var ErrUsernameExits = errors.New("username already exists")
var ErrEmpty = errors.New("username or password is empty")

type User struct{}

type UserLoginReq struct {
	Username string `form:"username" binding:"required,min=1,max=32"`
	Password string `form:"password" binding:"required,min=6,max=32"`
}

type UserLoginResp struct {
	UserId uint `json:"user_id"`
}

type UserRegisterReq struct {
	Username string `form:"username" binding:"required,min=1,max=32"`
	Password string `form:"password" binding:"required,min=6,max=32"`
}

type UserRegisterResp struct {
	UserId uint `json:"user_id"`
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
	var count int64
	err = db.MySQL.Debug().Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error
	if err != nil {
		log.Logger.Error("mysql happen error", zap.Error(err))
		return nil, err
	}
	if count != 0 {
		log.Logger.Error("user already exist")
		return nil, ErrUsernameExits
	}

	// 加密
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := model.User{
		Name:     req.Username,
		Username: req.Username,
		Password: string(hash),
	}

	db.MySQL.Debug().Create(&user)

	return UserRegisterResp{
		UserId: user.ID,
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

	var user model.User
	err = db.MySQL.Debug().Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		log.Logger.Error("mysql happen error")
		return nil, err
	}

	// 检查密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Logger.Error("password error", zap.Any("user", user))
		return nil, err
	}

	return UserLoginResp{
		UserId: user.ID,
	}, nil
}
