package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/qerdcv/muzlag.go/internal/bot"
	"github.com/qerdcv/muzlag.go/internal/config"
	"github.com/qerdcv/muzlag.go/internal/logger"
)

func main() {
	cfg := config.New()
	l := logger.NewTextLogger()

	b, err := bot.New(l, cfg.BotToken)
	if err != nil {
		l.Error("bot new", "err", err)
		return
	}

	l.Info("starting")
	if err = b.Run(); err != nil {
		l.Error("bot run", "err", err)
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	fmt.Println("shutting down")
	if err = b.Close(); err != nil {
		l.Error("bot close", "err", err)
		return
	}
}
