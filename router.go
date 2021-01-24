package main

import (
	"github.com/gin-gonic/gin"
)

func initRouter() {
	router := gin.Default()
	router.POST("/api", Handler)
	log.Fatal(router.Run("127.0.0.1:21125"))
}
