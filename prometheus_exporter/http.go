package prometheus_exporter

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
)

func RunHTTPServer(logger *slog.Logger) {
	go CollectSmartCtlStats(logger)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/dashboard.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/monitoring/dashboard.json")
	})
	http.HandleFunc("/prometheusRule.yml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/monitoring/prometheusRule.yml")
	})

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		logger.Error("failed to run http server", "error", err)
		return
	}
}
