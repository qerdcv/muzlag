package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/qerdcv/muzlag.go/internal/bot"
	"github.com/qerdcv/muzlag.go/internal/bot/handler"
	"github.com/qerdcv/muzlag.go/internal/bot/queue"
	"github.com/qerdcv/muzlag.go/internal/config"
	"github.com/qerdcv/muzlag.go/internal/logger"
	"github.com/qerdcv/muzlag.go/internal/metrics"
)

func main() {
	cfg := config.New()
	l := logger.NewTextLogger()

	b, err := bot.New(
		l,
		cfg.BotToken,
		handler.NewHealthHandler(),
		handler.NewPlayerHandler(
			l,
			queue.New(),
		),
	)
	if err != nil {
		l.Error("bot new", "err", err)
		return
	}

	ctx := context.Background()
	signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	errG, ctx := errgroup.WithContext(ctx)
	l.Info("starting")
	errG.Go(func() error { return b.Run(ctx) })
	errG.Go(func() error {
		return metrics.RunMetricsHandler(ctx, l)
	})

	if err = errG.Wait(); err != nil {
		l.Error("error group wait", "err", err)
	}
}
