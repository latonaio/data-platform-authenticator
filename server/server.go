package main

import (
	"context"
	"data-platform-authenticator/configs"
	"data-platform-authenticator/pkg/db"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"data-platform-authenticator/pkg/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfgs, err := configs.New()
	if err != nil {
		log.Fatalf("failed to set configs: %v", err)
	}
	errC := make(chan error)
	quitC := make(chan os.Signal, 1)
	signal.Notify(quitC, syscall.SIGTERM, os.Interrupt)

	echoServer := server.New(ctx, cfgs)
	err = db.NewDBConPool(ctx, cfgs)
	if err != nil {
		panic(err)
	}
	go echoServer.Start(errC)

	select {
	case err := <-errC:
		panic(err)
	case <-quitC:
		if err := echoServer.Shutdown(ctx); err != nil {
			errC <- err
		}
		cancel()
		time.Sleep(1 * time.Second)
	}
}
