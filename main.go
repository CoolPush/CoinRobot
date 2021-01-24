package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("程序启动:" + "CoinRobot")
	gin.SetMode(gin.ReleaseMode)

	//服务运行
	initRouter()
}
