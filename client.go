package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

func checkServerHealth(ctx context.Context) error {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	host := net.JoinHostPort("localhost", port)

	slog.InfoContext(ctx, "Checking server health", "host", host)

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s", host), nil)
	if err != nil {
		return err
	}

	_, err = client.Do(req)
	return err
}
