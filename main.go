package main

import (
	"fmt"
	"gin-example/gin-blog/routers"
	"gin-example/gin-blog/setting"
	"log"
	"syscall"

	"github.com/fvbock/endless"
)

func main() {
	endless.DefaultReadTimeOut = setting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20

	server := endless.NewServer(
		fmt.Sprintf(":%d", setting.HttpPort),
		routers.InitRouter(),
	)

	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid = %d\n", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server error: %v\n", err)
	}
}
