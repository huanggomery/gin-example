package main

import (
	"fmt"
	"gin-example/gin-blog/routers"
	"gin-example/gin-blog/setting"
	"net/http"
)

func main() {
	router := routers.InitRouter()

	server := http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()
}
