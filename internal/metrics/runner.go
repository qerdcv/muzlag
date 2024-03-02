package metrics

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/qerdcv/muzlag.go/internal/logger"
)

func RunMetricsHandler(ctx context.Context, l logger.Logger[*slog.Logger]) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	s := http.Server{
		Addr:    ":8009",
		Handler: mux,
	}

	go func() {
		<-ctx.Done()

		cancelCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		go func() {
			<-cancelCtx.Done()
			if err := s.Close(); err != nil {
				l.Error("server close", "err", err)
			}
		}()

		l.Info("shutting down metrics server")
		if err := s.Shutdown(cancelCtx); err != nil {
			l.Error("server shutdown", "err", err)
		}
	}()

	l.Info("running metrics server")
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server listen and serve: %w", err)
	}

	return nil
}
