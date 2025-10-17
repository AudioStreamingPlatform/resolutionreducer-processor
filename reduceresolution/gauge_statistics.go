// Copyright (C) 2025 Bang & Olufsen A/S, Denmark
//
// SPDX-License-Identifier: GPL-2.0-or-later

package reduceresolution

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type GaugeValue interface {
	int64 | float64
}

type GaugeAggregate[T GaugeValue] struct {
	count       int64
	average     T
	sum         T
	max         T
	min         T
	name        string
	description string
	unit        string
	attributes  pcommon.Map
	startTS     pcommon.Timestamp
}

func CreateGaugeAggregate[T GaugeValue](metric pmetric.Metric, attributes pcommon.Map, startTS pcommon.Timestamp, value T) *GaugeAggregate[T] {
	return &GaugeAggregate[T]{
		count:       1,
		max:         value,
		min:         value,
		sum:         value,
		name:        metric.Name(),
		description: metric.Description(),
		unit:        metric.Unit(),
		attributes:  attributes,
		startTS:     startTS,
	}
}

func AggregateGauge[T GaugeValue](aggregate *GaugeAggregate[T], startTS pcommon.Timestamp, value T) {
	aggregate.count++
	aggregate.sum += value
	if aggregate.min > value {
		aggregate.min = value
	}
	if aggregate.max < value {
		aggregate.max = value
	}
	if startTS < aggregate.startTS {
		aggregate.startTS = startTS
	}
}

func CreateGaugeMetrics[T GaugeValue](scope pmetric.ScopeMetrics, aggregate *GaugeAggregate[T], aggregationTS pcommon.Timestamp) {

	createSpecificMetric := func(scope pmetric.ScopeMetrics, aggregate *GaugeAggregate[T], sufix string, value T) {
		metric := scope.Metrics().AppendEmpty()
		metric.SetName(aggregate.name + sufix)
		metric.SetUnit(aggregate.unit)
		metric.SetDescription(aggregate.description)
		gauge := metric.SetEmptyGauge()
		gauge_dp := gauge.DataPoints().AppendEmpty()
		gauge_dp.SetStartTimestamp(aggregate.startTS)
		gauge_dp.SetTimestamp(aggregationTS)
		aggregate.attributes.CopyTo(gauge_dp.Attributes())
		switch v := any(value).(type) {
		case int64:
			gauge_dp.SetIntValue(v)
		case float64:
			gauge_dp.SetDoubleValue(v)
		}
	}
	aggregate.average = aggregate.sum / T(aggregate.count)

	createSpecificMetric(scope, aggregate, "_gauge_avg", aggregate.average)
	createSpecificMetric(scope, aggregate, "_gauge_max", aggregate.max)
	createSpecificMetric(scope, aggregate, "_gauge_min", aggregate.min)

	metric := scope.Metrics().AppendEmpty()
	metric.SetName(aggregate.name + "_gauge_count")
	metric.SetDescription(aggregate.description)
	gauge := metric.SetEmptyGauge()
	gauge_dp := gauge.DataPoints().AppendEmpty()
	aggregate.attributes.CopyTo(gauge_dp.Attributes())
	gauge_dp.SetIntValue(aggregate.count)
	gauge_dp.SetStartTimestamp(aggregate.startTS)
	gauge_dp.SetTimestamp(aggregationTS)
}
