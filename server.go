package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
)

func app() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		var response struct {
			URL struct {
				Host  string              `json:"host"`
				Path  string              `json:"path"`
				Query map[string][]string `json:"query"`
			}
			Headers map[string][]string `json:"headers"`
			Method  string              `json:"method"`
		}

		response.URL.Host = r.Host
		response.URL.Path = r.URL.Path
		response.URL.Query = r.URL.Query()
		response.Headers = r.Header
		response.Method = r.Method

		out, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(out)
	})
	return mux
}

func runApp(ctx context.Context, addr string, handler http.Handler) error {
	rootCtx, cancelRoot := signal.NotifyContext(ctx, signals...)
	defer cancelRoot()

	// In-flight requests get a context that won't be immediately cancelled on SIGINT/SIGTERM
	// so that they can be gracefully stopped.
	ongoingCtx, cancelOngoing := context.WithCancel(context.WithoutCancel(rootCtx))
	server := &http.Server{
		Addr: addr,
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
		Handler: handler,
	}

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		slog.InfoContext(rootCtx, "Server listening", "addr", addr)
		if err := server.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-rootCtx.Done():
		slog.InfoContext(rootCtx, "Received shutdown signal, shutting down")
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			cancelOngoing()
			return err
		}
	}

	slog.InfoContext(rootCtx, "Waiting for ongoing requests to finish")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.WithoutCancel(rootCtx), shutdownPeriod)
	defer cancelShutdown()
	err := server.Shutdown(shutdownCtx)
	cancelOngoing()
	if err != nil {
		slog.ErrorContext(rootCtx, "Failed to wait for ongoing requests to finish, waiting for forced cancellation")
		timeSleep(shutdownHardPeriod)
	}

	slog.InfoContext(rootCtx, "Server shut down")
	return err
}
