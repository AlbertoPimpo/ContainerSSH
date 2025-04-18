package metricsintegration

import (
    "go.containerssh.io/containerssh/config"
    "go.containerssh.io/containerssh/internal/metrics"
    "go.containerssh.io/containerssh/internal/sshserver"
)

func NewHandler(
	cfg config.MetricsConfig,
	metricsCollector metrics.Collector,
	backend sshserver.Handler,
) (sshserver.Handler, error) {
	if !cfg.Enable {
		return backend, nil
	}

	connectionsMetric := metricsCollector.MustCreateCounterGeo(
		MetricNameConnections,
		"connections_total",
		MetricHelpConnections,
	)
	currentConnectionsMetric := metricsCollector.MustCreateGaugeGeo(
		MetricNameCurrentConnections,
		"connections",
		MetricHelpCurrentConnections,
	)

	handshakeSuccessfulMetric := metricsCollector.MustCreateCounterGeo(
		MetricNameSuccessfulHandshake,
		"handshakes_total",
		MetricHelpSuccessfulHandshake,
	)
	handshakeFailedMetric := metricsCollector.MustCreateCounterGeo(
		MetricNameFailedHandshake,
		"handshakes_total",
		MetricHelpFailedHandshake,
	)

	return &metricsHandler{
		backend:                   backend,
		metricsCollector:          metricsCollector,
		connectionsMetric:         connectionsMetric,
		handshakeSuccessfulMetric: handshakeSuccessfulMetric,
		handshakeFailedMetric:     handshakeFailedMetric,
		currentConnectionsMetric:  currentConnectionsMetric,
	}, nil
}
