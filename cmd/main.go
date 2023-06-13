package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-talk/common/config"
	"go-talk/common/db"
	"go-talk/common/log"
	"sync"
)

var once sync.Once

// 初始化函数
func init() {
	once.Do(func() {
		config.ReadCfg()
		config.Init()
		log.Init()
		db.Init()
	})
}

// main 函数
func main() {
	r := gin.Default()

	handle(r)

	r.Run(fmt.Sprintf("%s:%s", config.AppCfg.Host, config.AppCfg.Port))
}
