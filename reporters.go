package nano

import (
	"github.com/nut-game/nano/config"
	"github.com/nut-game/nano/logger"
	"github.com/nut-game/nano/metrics"
	"github.com/nut-game/nano/metrics/models"
)

// CreatePrometheusReporter create a Prometheus reporter instance
func CreatePrometheusReporter(serverType string, config config.MetricsConfig, customSpecs models.CustomMetricsSpec) (*metrics.PrometheusReporter, error) {
	logger.Log.Infof("prometheus is enabled, configuring reporter on port %d", config.Prometheus.Port)
	prometheus, err := metrics.GetPrometheusReporter(serverType, config, customSpecs)
	if err != nil {
		logger.Log.Errorf("failed to start prometheus metrics reporter, skipping %v", err)
	}
	return prometheus, err
}

// CreateStatsdReporter create a Statsd reporter instance
func CreateStatsdReporter(serverType string, config config.MetricsConfig) (*metrics.StatsdReporter, error) {
	logger.Log.Infof(
		"statsd is enabled, configuring the metrics reporter with host: %s",
		config.Statsd.Host,
	)
	metricsReporter, err := metrics.NewStatsdReporter(
		config,
		serverType,
	)
	if err != nil {
		logger.Log.Errorf("failed to start statds metrics reporter, skipping %v", err)
	}
	return metricsReporter, err
}
