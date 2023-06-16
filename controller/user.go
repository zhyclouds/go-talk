package controller

import (
	"github.com/gin-gonic/gin"
	res "go-talk/common/result"
	"go-talk/service"
)

func Register(c *gin.Context) {
	var u service.User
	register, err := u.Register(c)
	if err != nil {
		if err == service.ErrUsernameExits {
			res.Error(c, res.Status{
				StatusCode: res.UsernameExitErrorStatus.StatusCode,
				StatusMsg:  res.UsernameExitErrorStatus.StatusMsg,
			})
		} else {
			res.Error(c, res.Status{
				StatusCode: res.RegisterErrorStatus.StatusCode,
				StatusMsg:  res.RegisterErrorStatus.StatusMsg,
			})
		}
		return
	}
	data := register.(service.UserRegisterResp)
	res.Success(c, res.R{
		"token": data.Token,
	})
}

func Login(c *gin.Context) {
	var u service.User
	login, err := u.Login(c)
	if err != nil {
		if err == service.PasswordErr {
			res.Error(c, res.Status{
				StatusCode: res.PasswordErrorStatus.StatusCode,
				StatusMsg:  res.PasswordErrorStatus.StatusMsg,
			})
			return
		} else {
			res.Error(c, res.Status{
				StatusCode: res.LoginErrorStatus.StatusCode,
				StatusMsg:  res.LoginErrorStatus.StatusMsg,
			})
			return
		}
	}
	data := login.(service.UserLoginResp)
	res.Success(c, res.R{
		"token": data.Token,
	})
}

func AddFriend(c *gin.Context) {
	var u service.User
	addFriend, err := u.AddFriend(c)
	if err != nil {
		if err == service.ErrIdEmpty {
			res.Error(c, res.Status{
				StatusCode: res.IdEmptyErrorStatus.StatusCode,
				StatusMsg:  res.IdEmptyErrorStatus.StatusMsg,
			})
		} else if err == service.ErrIsFriend {
			res.Error(c, res.Status{
				StatusCode: res.IsFriendErrorStatus.StatusCode,
				StatusMsg:  res.IsFriendErrorStatus.StatusMsg,
			})
		} else {
			res.Error(c, res.Status{
				StatusCode: res.AddFriendErrorStatus.StatusCode,
				StatusMsg:  res.AddFriendErrorStatus.StatusMsg,
			})
		}
		return
	}
	data := addFriend.(service.UserAddFriendResp)
	res.Success(c, res.R{
		"message": data.Msg,
	})
}

func DeleteFriend(c *gin.Context) {
	var u service.User
	delFriend, err := u.DeleteFriend(c)
	if err != nil {
		if err == service.ErrIdEmpty {
			res.Error(c, res.Status{
				StatusCode: res.IdEmptyErrorStatus.StatusCode,
				StatusMsg:  res.IdEmptyErrorStatus.StatusMsg,
			})
		} else if err == service.ErrNotFriend {
			res.Error(c, res.Status{
				StatusCode: res.NotFriendErrorStatus.StatusCode,
				StatusMsg:  res.NotFriendErrorStatus.StatusMsg,
			})
		} else {
			res.Error(c, res.Status{
				StatusCode: res.DeleteFriendErrorStatus.StatusCode,
				StatusMsg:  res.DeleteFriendErrorStatus.StatusMsg,
			})
		}
	}

	data := delFriend.(service.UserDelFriendResp)
	res.Success(c, res.R{
		"message": data.Msg,
	})
}
