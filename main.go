package main

import (
	"log/slog"
	"os"
	"smartctl-prometheus-exporter/prometheus_exporter"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	prometheus_exporter.RunHTTPServer(logger)
}
