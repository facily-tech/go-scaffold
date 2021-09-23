package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/facily-tech/go-scaffold/internal/config"
	"github.com/facily-tech/go-scaffold/pkg/core/types"
	"go.uber.org/zap"
)

type Config struct {
	Addr             string
	GracefulDuration time.Duration
}

func Run(ctx context.Context, cnf Config, handler http.Handler, log *zap.Logger) {
	server := &http.Server{Addr: cnf.Addr, Handler: handler}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(ctx)

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, cnf.GracefulDuration)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out, forcing exit")
			}
		}()

		// Trigger graceful shutdown
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal("shutdown error", zap.Error(err))
		}
		serverStopCtx()
	}()

	startingMessage(ctx, cnf.Addr, log)

	if err := server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		log.Warn("HTTP server requested to stop")
	} else {
		log.Error("HTTP server stopped with error", zap.Error(err))
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

func startingMessage(ctx context.Context, where string, log *zap.Logger) {
	v, ok := ctx.Value(types.ContextKey(types.Version)).(*config.Version)
	if !ok {
		log.Warn("could not get version, received")
	}

	t, ok := ctx.Value(types.ContextKey(types.StartedAt)).(time.Time)
	if !ok {
		fmt.Println(reflect.TypeOf(ctx.Value(types.ContextKey(types.StartedAt))))
		log.Warn("could not get startedTime time")
	}

	log.Info(
		"Starting API Server",
		zap.String("bind", where),
		zap.Time("start time", t),
		zap.String("version", v.GitCommitHash),
	)
}
