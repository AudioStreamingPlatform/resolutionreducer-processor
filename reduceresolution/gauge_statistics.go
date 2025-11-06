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

func Abs[T GaugeValue](value T) T {
	if value < 0 {
		return -value
	}
	return value
}

type GaugeAggregate[T GaugeValue] struct {
	count       int64
	average     T
	sum         T
	max         T
	min         T
	max_abs     T
	min_abs     T
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
		max_abs:     Abs(value),
		min_abs:     Abs(value),
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

	if aggregate.min_abs > Abs(value) {
		aggregate.min_abs = value
	}

	if aggregate.max_abs < Abs(value) {
		aggregate.max_abs = value
	}

	if startTS < aggregate.startTS {
		aggregate.startTS = startTS
	}
}

func CreateGaugeMetrics[T GaugeValue](scope pmetric.ScopeMetrics, aggregate *GaugeAggregate[T], aggregationTS pcommon.Timestamp, Config ProcessedConfig) {

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

	// The average is commented in order to reduce the number of metrics generated
	// from one gauge, however this show an example on the gauge can be broken down
	// into more or less metrics depending on what is required
	//  createSpecificMetric(scope, aggregate, "_gauge_avg", aggregate.average)

	if Config.RealMaxMinAggregation[aggregate.name] {
		createSpecificMetric(scope, aggregate, "_gauge_max", aggregate.max)
		createSpecificMetric(scope, aggregate, "_gauge_min", aggregate.min)
	} else {
		createSpecificMetric(scope, aggregate, "_gauge_abs_max", aggregate.max_abs)
		createSpecificMetric(scope, aggregate, "_gauge_abs_min", aggregate.min_abs)
	}

}
