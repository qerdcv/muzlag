package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/qerdcv/muzlag.go/internal/bot"
	"github.com/qerdcv/muzlag.go/internal/config"
)

func main() {
	cfg, err := config.New(os.Args[1])
	if err != nil {
		panic(err)
	}

	b, err := bot.New(cfg.BotToken)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting bot...")
	if err = b.Run(); err != nil {
		panic(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	fmt.Println("stopping bot...")

	if err = b.Close(); err != nil {
		panic(err)
	}
}
