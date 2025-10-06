// Copyright (C) 2025 Bang & Olufsen A/S, Denmark
//
// SPDX-License-Identifier: MIT License

package reduceresolution

import (
	"fmt"
	"sort"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type ScopeContainer struct {
	scopeName       string
	scopeVersion    string
	scopeAttributes pcommon.Map

	intGaugeAggregate   map[string]*GaugeAggregate[int64]
	floatGaugeAggregate map[string]*GaugeAggregate[float64]

	intCounterAggregate   map[string]*CounterAggregate[int64]
	floatCounterAggregate map[string]*CounterAggregate[float64]

	leftoverMetric []pmetric.Metric
}

func CreateScopeContainer(scopeMetric pmetric.ScopeMetrics) *ScopeContainer {
	return &ScopeContainer{
		scopeName:             scopeMetric.Scope().Name(),
		scopeVersion:          scopeMetric.Scope().Version(),
		scopeAttributes:       scopeMetric.Scope().Attributes(),
		intGaugeAggregate:     make(map[string]*GaugeAggregate[int64]),
		floatGaugeAggregate:   make(map[string]*GaugeAggregate[float64]),
		intCounterAggregate:   make(map[string]*CounterAggregate[int64]),
		floatCounterAggregate: make(map[string]*CounterAggregate[float64]),
		leftoverMetric:        make([]pmetric.Metric, 0),
	}
}

// Creates a unique deterministic key based on a scope's name, version, and its attributes
func CreateScopeKey(scopeMetric pmetric.ScopeMetrics) string {
	scope_keys := make([]string, 0, scopeMetric.Scope().Attributes().Len())
	for k := range scopeMetric.Scope().Attributes().AsRaw() {
		scope_keys = append(scope_keys, k)
	}

	sort.Strings(scope_keys)
	var attributeParts []string
	for _, k := range scope_keys {
		value, _ := scopeMetric.Scope().Attributes().Get(k)
		attributeParts = append(attributeParts, fmt.Sprintf("%s=%s", k, value.AsString()))
	}
	attributesStrings := strings.Join(attributeParts, ",")
	return fmt.Sprintf("%s|%s|%s", scopeMetric.Scope().Name(), scopeMetric.Scope().Version(), attributesStrings)

}

// Creates a unique deterministic key based on a metric's name, and its attributes
func CreateMetricKey(metric pmetric.Metric, attributes pcommon.Map) string {
	metric_keys := make([]string, 0, attributes.Len())
	for k := range attributes.AsRaw() {
		metric_keys = append(metric_keys, k)
	}

	sort.Strings(metric_keys)

	var attributeParts []string
	for _, k := range metric_keys {
		value, _ := attributes.Get(k)
		attributeParts = append(attributeParts, fmt.Sprintf("%s=%s", k, value.AsString()))
	}
	attributesStrings := strings.Join(attributeParts, ",")
	return fmt.Sprintf("%s@%s", metric.Name(), attributesStrings)
}
