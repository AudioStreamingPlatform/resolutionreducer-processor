// Copyright (C) 2025 Bang & Olufsen A/S, Denmark
//
// SPDX-License-Identifier: GPL-2.0-or-later

package reduceresolution

type Config struct {
	MetricStatistics map[string][]string `mapstructure:"gauge-aggregations"`
}

type ProcessedConfig struct {
	MetricsStatistics map[string][]string
}

// Validate checks if the receiver configuration is valid
func (cfg *Config) Validate() error {
	return nil
}
