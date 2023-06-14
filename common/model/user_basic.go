package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"not null;unique;index"`
	Username string `gorm:"not null;unique;index"`
	Password string `gorm:"not null"`
	Avatar   string `gorm:"default:'https://p3.itc.cn/images03/20200516/0346a117a87b453fbd6d7b1d6698923d.jpeg'"`
}
