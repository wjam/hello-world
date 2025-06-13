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
	slog.InfoContext(ctx, "Checking server health")
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	host := net.JoinHostPort("127.0.0.1", port)

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s", host), nil)
	if err != nil {
		return err
	}

	_, err = client.Do(req)
	return err
}
