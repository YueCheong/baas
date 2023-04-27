package main

import (
	"bit-bass/core"
	"bit-bass/webserver"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	core.InitCore()
	r := webserver.InitRouter()

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Listen:%s", err)
		}
	}()

	//接受ctrl+c信号，关闭系统
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutdown server")

	//关闭服务器前关闭所有网络
	for _, b := range core.CoreBlockChainManager.GetBlockchains() {
		core.CoreBlockChainManager.StopAndRemoveBlockchainById(b.GetId())
	}

	for _, n := range core.CoreBlockChainManager.GetNets() {
		n.RemoveNet()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		fmt.Println("Server shutdown err :", err)
	}
}
