package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"
)

func checkServerHealth(ctx context.Context, enableHttps bool) error {
	client := &http.Client{
		Timeout: 1 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				// If the server is serving HTTPS, it's probably got a self-signed certificate anyway
				InsecureSkipVerify: true,
			},
		},
	}

	host := net.JoinHostPort("localhost", port)

	scheme := "http"
	if enableHttps {
		scheme = "https"
	}
	url := fmt.Sprintf("%s://%s", scheme, host)

	slog.InfoContext(ctx, "Checking server health", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	b, err := client.Do(req)
	defer func() {
		_, _ = io.Copy(io.Discard, b.Body)
		_ = b.Body.Close()
	}()
	return err
}
