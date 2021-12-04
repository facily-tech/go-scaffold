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

	coreLog "github.com/facily-tech/go-core/log"
	"github.com/facily-tech/go-scaffold/internal/config"
	"github.com/facily-tech/go-scaffold/pkg/core/types"
)

type Config struct {
	Addr             string
	GracefulDuration time.Duration
}

func Run(ctx context.Context, cnf Config, handler http.Handler, log coreLog.Logger) {
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
				log.Fatal(ctx, "graceful shutdown timed out, forcing exit")
			}
		}()

		// Trigger graceful shutdown
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(ctx, "shutdown error", coreLog.Error(err))
		}
		serverStopCtx()
	}()

	startingMessage(ctx, cnf.Addr, log)

	if err := server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		log.Warn(ctx, "HTTP server requested to stop")
	} else {
		log.Error(ctx, "HTTP server stopped with error", coreLog.Error(err))
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

func startingMessage(ctx context.Context, where string, log coreLog.Logger) {
	v, ok := ctx.Value(types.ContextKey(types.Version)).(*config.Version)
	if !ok {
		log.Warn(ctx, "could not get version, received")
	}

	t, ok := ctx.Value(types.ContextKey(types.StartedAt)).(time.Time)
	if !ok {
		fmt.Println(reflect.TypeOf(ctx.Value(types.ContextKey(types.StartedAt))))
		log.Warn(ctx, "could not get startedTime time")
	}

	log.Info(ctx,
		"Starting API Server",
		coreLog.Any("bind", where),
		coreLog.Any("start time", t),
		coreLog.Any("version", v.GitCommitHash),
	)
}
