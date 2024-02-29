package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"

	"github.com/qerdcv/muzlag.go/internal/bot"
	"github.com/qerdcv/muzlag.go/internal/bot/handler"
	"github.com/qerdcv/muzlag.go/internal/bot/player"
	"github.com/qerdcv/muzlag.go/internal/config"
	"github.com/qerdcv/muzlag.go/internal/logger"
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
			player.New(),
		),
	)
	if err != nil {
		l.Error("bot new", "err", err)
		return
	}

	ctx := context.Background()
	signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	s := http.Server{
		Addr: ":8009",
	}

	errG, ctx := errgroup.WithContext(ctx)
	l.Info("starting")
	errG.Go(b.Run)
	errG.Go(func() error {
		http.Handle("/metrics", promhttp.Handler())
		if serveErr := s.ListenAndServe(); serveErr != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http listen and serve: %w", serveErr)
		}

		return nil
	})

	<-ctx.Done()

	l.Info("shutting down")

	if err = b.Close(); err != nil {
		l.Error("bot close", "err", err)
		return
	}

	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-cancelCtx.Done()
		if closeErr := s.Close(); closeErr != nil {
			l.Error("server close", "err", closeErr)
		}
	}()

	if err = s.Shutdown(cancelCtx); err != nil {
		l.Error("server shutdown: %w", err)
	}
}
