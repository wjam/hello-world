package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"syscall"
	"time"
)

var (
	signals            = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	shutdownPeriod     = 15 * time.Second
	shutdownHardPeriod = 3 * time.Second
	timeSleep          = time.Sleep
	port               = "8080"
)

func main() {
	// Reminder: `defer` doesn't behave as expected in functions with log.Fatal, os.Exit, etc.
	rootCtx := context.Background()

	healthCheck := flag.Bool("health", false, "Check if server is healthy")
	flag.Parse()

	if *healthCheck {
		if err := checkServerHealth(rootCtx); err != nil {
			slog.ErrorContext(rootCtx, "health check failed", "error", err)
			os.Exit(1)
		}
		return
	}

	app := app()

	if err := runApp(rootCtx, fmt.Sprintf(":%s", port), app); err != nil {
		slog.ErrorContext(rootCtx, "failed to run app", "error", err)
		os.Exit(1)
	}
}
