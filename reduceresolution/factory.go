// Copyright (C) 2025 Bang & Olufsen A/S, Denmark
//
// SPDX-License-Identifier: GPL-2.0-or-later

package reduceresolution

import (
	"context"
	"strings"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

// NewFactory creates a factory for the metric logging processor
func NewFactory() processor.Factory {
	return processor.NewFactory(
		"reduceresolution",
		createDefaultConfig,
		processor.WithMetrics(createMetricsProcessor, component.StabilityLevelBeta))
}

// createDefaultConfig creates the default configuration for the processor
func createDefaultConfig() component.Config {
	return &Config{}
}

// createMetricsProcessor creates a new instance of the metric logging processor
func createMetricsProcessor(
	ctx context.Context,
	settings processor.CreateSettings,
	config component.Config,
	nextConsumer consumer.Metrics,
) (processor.Metrics, error) {
	c := config.(*Config)

	var processedConfig ProcessedConfig
	processedConfig.MetricsStatistics = map[string][]string{}
	for metricName, statisticsList := range c.MetricStatistics {
		processedConfig.MetricsStatistics[strings.ToLower(metricName)] = statisticsList
	}

	logProcessor := &ReduceResolution{
		Logger: settings.Logger,
		Config: processedConfig,
	}

	return processorhelper.NewMetricsProcessor(
		ctx,
		settings,
		config,
		nextConsumer,
		logProcessor.ProcessMetrics,
		processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}))
}
