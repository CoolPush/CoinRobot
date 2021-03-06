package main

import (
	"github.com/gin-gonic/gin"
)

func initRouter() {
	router := gin.New()
	router.Use(gin.Recovery())
	router.POST("/api/coin", HandlerCoin)
	router.POST("/api/lsp", HandlerLSP)
	log.Fatal(router.Run("127.0.0.1:21125"))
}
