// Copyright (C) 2025 Bang & Olufsen A/S, Denmark
//
// SPDX-License-Identifier: GPL-2.0-or-later

package reduceresolution

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type CounterValue interface {
	int64 | float64
}

type CounterAggregate[T CounterValue] struct {
	value       T
	name        string
	description string
	unit        string
	attributes  pcommon.Map
	startTS     pcommon.Timestamp
	lastTS      pcommon.Timestamp
	aggregation pmetric.AggregationTemporality
	monotonic   bool
}

func CreateCounterAggregate[T CounterValue](metric pmetric.Metric, attributes pcommon.Map, startTS pcommon.Timestamp, lastTS pcommon.Timestamp, value T) *CounterAggregate[T] {
	return &CounterAggregate[T]{
		value:       value,
		name:        metric.Name(),
		description: metric.Description(),
		unit:        metric.Unit(),
		attributes:  attributes,
		startTS:     startTS,
		lastTS:      lastTS,
		aggregation: metric.Sum().AggregationTemporality(),
		monotonic:   metric.Sum().IsMonotonic(),
	}
}

func AggregateCounter[T CounterValue](aggregate *CounterAggregate[T], startTS pcommon.Timestamp, lastTS pcommon.Timestamp, value T) {
	switch aggregate.aggregation {
	case pmetric.AggregationTemporalityCumulative:
		if aggregate.lastTS < lastTS {
			aggregate.value = value
			if startTS < aggregate.startTS {
				aggregate.startTS = startTS
			}
			aggregate.lastTS = lastTS
		}
	case pmetric.AggregationTemporalityDelta:
		aggregate.value += value
		if startTS < aggregate.startTS {
			aggregate.startTS = startTS
		}
	}
}

func CreateCounterMetrics[T GaugeValue](scope pmetric.ScopeMetrics, aggregate *CounterAggregate[T], aggregationTS pcommon.Timestamp) {
	metric_value := scope.Metrics().AppendEmpty()
	metric_value.SetName(aggregate.name)
	metric_value.SetUnit(aggregate.unit)
	metric_value.SetDescription(aggregate.description)
	counter := metric_value.SetEmptySum()
	counter.SetAggregationTemporality(aggregate.aggregation)
	counter.SetIsMonotonic(aggregate.monotonic)
	counter_dp := counter.DataPoints().AppendEmpty()
	counter_dp.SetStartTimestamp(aggregate.startTS)
	counter_dp.SetTimestamp(aggregationTS)
	aggregate.attributes.CopyTo(counter_dp.Attributes())
	switch v := any(aggregate.value).(type) {
	case int64:
		counter_dp.SetIntValue(v)
	case float64:
		counter_dp.SetDoubleValue(v)
	}
}
