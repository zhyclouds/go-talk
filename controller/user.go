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
		"userId": data.UserId,
	})
}

func Login(c *gin.Context) {
	var u service.User
	login, err := u.Login(c)
	if err != nil {
		res.Error(c, res.Status{
			StatusCode: res.LoginErrorStatus.StatusCode,
			StatusMsg:  res.LoginErrorStatus.StatusMsg,
		})
		return
	}
	data := login.(service.UserLoginResp)
	res.Success(c, res.R{
		"userId": data.UserId,
	})
}
