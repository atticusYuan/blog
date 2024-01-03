package main

import (
	"context"
	"fmt"
	"gogin/models"
	"gogin/pkg/gredis"
	"gogin/pkg/logging"
	"gogin/pkg/setting"
	"gogin/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen: %s\n", err)
		}
	}()
	fmt.Println("Addr:", setting.ServerSetting.HttpPort)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	//context.Backgroud()用于生成一个空的根上下文，位于整个上下文链的最上层，其目的在于：
	//1.保持所有上下文链的一致性
	//2.可以通过不同的函数来添加一些特定的属性，例如context.WithTimeout，生成一个带有5秒超时限制的上下文，一旦相关操作超过5秒，
	//便会触发取消(cancel)操作，从而及时结束或清理相关操作，避免长时间等待而被阻塞。
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
