package metrics

import (
	"context"
	"errors"
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

		cancelCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		go func() {
			<-cancelCtx.Done()
			if err := s.Close(); err != nil {
				l.Error("server close", "err", err)
			}
		}()

		if err := s.Shutdown(cancelCtx); err != nil {
			l.Error("server shutdown", "err", err)
		}
	}()

	l.Info("running metrics server")
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		l.Error("metrics listen and serve", "err", err)
	}

	return nil
}
