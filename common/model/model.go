package model

type User struct {
	Id       int64  `gorm:"primaryKey"`
	Name     string `gorm:"not null;unique;index"`
	Username string `gorm:"not null;unique;index"`
	Password string `gorm:"not null"`
}
