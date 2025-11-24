package main

import (
	"context"
	"net/http"
	"os"

	"github.com/UsatovPavel/PRAssign/internal/config"
)

func main() {
	client := http.Client{Timeout: config.HTTPClientTimeoutShort}

	ctx, cancel := context.WithTimeout(context.Background(), config.HTTPClientTimeoutShort)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:8080/health", nil)
	if err != nil {
		os.Exit(config.ExitCodeRuntimeError)
	}

	resp, err := client.Do(req)
	if err != nil {
		os.Exit(config.ExitCodeRuntimeError)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		os.Exit(0)
	}
	os.Exit(config.ExitCodeConfigError)
}
