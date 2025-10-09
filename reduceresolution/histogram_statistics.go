// Copyright (C) 2025 Bang & Olufsen A/S, Denmark
//
// SPDX-License-Identifier: GPL-2.0-or-later

package reduceresolution

import (
	"slices"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type HistogramAggregate struct {
	count          uint64
	sum            float64
	max            float64
	min            float64
	bucketCounts   []uint64
	explicitBounds []float64
	name           string
	unit           string
	description    string
	attributes     pcommon.Map
	startTS        pcommon.Timestamp
	lastTS         pcommon.Timestamp
	aggregation    pmetric.AggregationTemporality
}

func CreateHistogramAggregate(metric pmetric.Metric, value pmetric.HistogramDataPoint) *HistogramAggregate {
	return &HistogramAggregate{
		count:          value.Count(),
		sum:            value.Sum(),
		max:            value.Max(),
		min:            value.Min(),
		bucketCounts:   value.BucketCounts().AsRaw(),
		explicitBounds: value.ExplicitBounds().AsRaw(),
		name:           metric.Name(),
		unit:           metric.Unit(),
		description:    metric.Description(),
		attributes:     value.Attributes(),
		startTS:        value.StartTimestamp(),
		lastTS:         value.Timestamp(),
		aggregation:    metric.Histogram().AggregationTemporality(),
	}
}

func AggregateHistogram(aggregate *HistogramAggregate, value pmetric.HistogramDataPoint) int16 {
	switch aggregate.aggregation {
	case pmetric.AggregationTemporalityCumulative:
		if aggregate.lastTS < value.Timestamp() {
			aggregate.count = value.Count()
			aggregate.sum = value.Sum()
			aggregate.max = value.Max()
			aggregate.min = value.Min()
			aggregate.bucketCounts = value.BucketCounts().AsRaw()
			aggregate.explicitBounds = value.ExplicitBounds().AsRaw()
			if value.StartTimestamp() < aggregate.startTS {
				aggregate.startTS = value.StartTimestamp()
			}
			aggregate.lastTS = value.Timestamp()
		}
	case pmetric.AggregationTemporalityDelta:
		if slices.Equal(aggregate.explicitBounds, value.ExplicitBounds().AsRaw()) &&
			len(aggregate.bucketCounts) == value.BucketCounts().Len() {
			for i := 0; i < value.BucketCounts().Len(); i++ {
				aggregate.bucketCounts[i] = aggregate.bucketCounts[i] + value.BucketCounts().At(i)
			}
		} else {
			return 1
		}
		aggregate.count += value.Count()
		aggregate.sum += value.Sum()
		if aggregate.max < value.Max() {
			aggregate.max = value.Max()
		}
		if aggregate.min > value.Min() {
			aggregate.min = value.Min()
		}
		aggregate.explicitBounds = value.ExplicitBounds().AsRaw()
		if value.StartTimestamp() < aggregate.startTS {
			aggregate.startTS = value.StartTimestamp()
		}
	}
	return 0
}

func CreateHistogramMetrics(scope pmetric.ScopeMetrics, aggregate *HistogramAggregate, aggregationTS pcommon.Timestamp) {
	metric_value := scope.Metrics().AppendEmpty()
	metric_value.SetName(aggregate.name)
	metric_value.SetUnit(aggregate.unit)
	metric_value.SetDescription(aggregate.description)
	histogram := metric_value.SetEmptyHistogram()
	histogram.SetAggregationTemporality(aggregate.aggregation)
	histogram_dp := histogram.DataPoints().AppendEmpty()
	histogram_dp.SetStartTimestamp(aggregate.startTS)
	histogram_dp.SetTimestamp(aggregationTS)
	aggregate.attributes.CopyTo(histogram_dp.Attributes())
	histogram_dp.SetCount(aggregate.count)
	histogram_dp.SetSum(aggregate.sum)
	histogram_dp.SetMax(aggregate.max)
	histogram_dp.SetMin(aggregate.min)

	for i := 0; i < len(aggregate.explicitBounds); i++ {
		histogram_dp.ExplicitBounds().Append(aggregate.explicitBounds[i])
	}

	for i := 0; i < len(aggregate.bucketCounts); i++ {
		histogram_dp.BucketCounts().Append(aggregate.bucketCounts[i])
	}
}
