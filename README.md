# smartctl-prometheus-exporter

Simple project to export smartctl metrics to Prometheus.

Docker image: `sokolimedia/smartctl-prometheus-exporter:latest`

Container has to be run with `privileged` option enabled.

Project exports http api on `:9000` with metrics at `/metrics` url.
