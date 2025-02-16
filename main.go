package main

import (
	"github.com/gin-gonic/gin"
	"github.com/panedbay/test-task/api"
	"github.com/panedbay/test-task/db"
)

func main() {
	db.Init()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.POST("/api/auth", api.PostAPIAuth)
	r.GET("/api/buy/:item", api.GetAPIBuyItem)
	r.POST("/api/sendCoin", api.PostAPISendCoin)
	r.GET("/api/info", api.GetAPIInfo)
	r.Run()
}
