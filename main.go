package main

import (
	"fmt"
	"myProject/gin-blog/setting"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	server := http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()
}
