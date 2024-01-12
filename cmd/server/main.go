package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ermanimer/apigateway/pkg/config"
	"github.com/ermanimer/apigateway/pkg/handlers/healthcheckhandler"
	"github.com/ermanimer/apigateway/pkg/handlers/upstreamhandler"
	"github.com/ermanimer/apigateway/pkg/server"
)

func main() {
	// create the l
	l := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// read the config
	c, err := config.ReadFromYaml("config.yaml")
	if err != nil {
		l.Error("failed to read config", "error", err)
	}

	// create the server
	s := server.New(c.Server)

	// register the health check handler
	s.RegisterHandler("/health-check", healthcheckhandler.New())

	// register upstream handlers
	for _, u := range c.Upstreams {
		s.RegisterHandler(u.Pattern, upstreamhandler.New(u))
		l.Info("upstream handler is registered", "pattern", u.Pattern, "url", u.URL)
	}

	// start the server
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		err = s.Start()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			l.Error("failed to start server", "error", err)
			cancel()
		}
	}()

	// wait for the shutdown signal
	<-ctx.Done()
	err = s.Shutdown()
	if err != nil {
		l.Error("failed to shutdown server", "error", err)
	}
}
