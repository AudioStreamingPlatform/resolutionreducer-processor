// Copyright (C) 2025 Bang & Olufsen A/S, Denmark
//
// SPDX-License-Identifier: GPL-2.0-or-later

package reduceresolution

type Config struct {
	RealMaxMinAggregation []string `mapstructure:"real-max-min-aggregation"`
}

type ProcessedConfig struct {
	RealMaxMinAggregation map[string]bool
}

// Validate checks if the receiver configuration is valid
func (cfg *Config) Validate() error {
	return nil
}
