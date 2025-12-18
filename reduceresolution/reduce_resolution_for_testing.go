// Copyright (C) 2025 Bang & Olufsen A/S, Denmark
//
// SPDX-License-Identifier: GPL-2.0-or-later

package reduceresolution

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func ValidateGauge(t *testing.T, metric pmetric.Metric, wasChecked *bool, startTs pcommon.Timestamp) {
	assert.False(t, *wasChecked)
	*wasChecked = true
	assert.Equal(t, pmetric.MetricTypeGauge, metric.Type())
	assert.Equal(t, 1, metric.Gauge().DataPoints().Len())
	assert.Equal(t, startTs, metric.Gauge().DataPoints().At(0).StartTimestamp())
}

func ValidateDoubleGauge(t *testing.T, metric pmetric.Metric, wasChecked *bool, value float64, startTs pcommon.Timestamp) {
	ValidateGauge(t, metric, wasChecked, startTs)
	assert.Equal(t, pmetric.NumberDataPointValueTypeDouble, metric.Gauge().DataPoints().At(0).ValueType())
	assert.Equal(t, value, metric.Gauge().DataPoints().At(0).DoubleValue())
}

func ValidateIntGauge(t *testing.T, metric pmetric.Metric, wasChecked *bool, value int64, startTs pcommon.Timestamp) {
	ValidateGauge(t, metric, wasChecked, startTs)
	assert.Equal(t, pmetric.NumberDataPointValueTypeInt, metric.Gauge().DataPoints().At(0).ValueType())
	assert.Equal(t, value, metric.Gauge().DataPoints().At(0).IntValue())
}

func ValidateCounter(t *testing.T, metric pmetric.Metric, wasChecked *bool, isCumulative bool, isMonotonic bool, startTs pcommon.Timestamp) {
	assert.False(t, *wasChecked)
	*wasChecked = true
	assert.Equal(t, pmetric.MetricTypeSum, metric.Type())
	assert.Equal(t, 1, metric.Sum().DataPoints().Len())
	if isCumulative {
		assert.Equal(t, pmetric.AggregationTemporalityCumulative, metric.Sum().AggregationTemporality())
	} else {
		assert.Equal(t, pmetric.AggregationTemporalityDelta, metric.Sum().AggregationTemporality())
	}
	assert.Equal(t, isMonotonic, metric.Sum().IsMonotonic())
	assert.Equal(t, startTs, metric.Sum().DataPoints().At(0).StartTimestamp())
}

func ValidateDoubleCounter(t *testing.T, metric pmetric.Metric, wasChecked *bool, isCumulative bool, isMonotonic bool, value float64, startTs pcommon.Timestamp) {
	ValidateCounter(t, metric, wasChecked, isCumulative, isMonotonic, startTs)
	assert.Equal(t, pmetric.NumberDataPointValueTypeDouble, metric.Sum().DataPoints().At(0).ValueType())
	assert.Equal(t, value, metric.Sum().DataPoints().At(0).DoubleValue())
}

func ValidateIntCounter(t *testing.T, metric pmetric.Metric, wasChecked *bool, isCumulative bool, isMonotonic bool, value int64, startTs pcommon.Timestamp) {
	ValidateCounter(t, metric, wasChecked, isCumulative, isMonotonic, startTs)
	assert.Equal(t, pmetric.NumberDataPointValueTypeInt, metric.Sum().DataPoints().At(0).ValueType())
	assert.Equal(t, value, metric.Sum().DataPoints().At(0).IntValue())
}

func ValidateHistogram(t *testing.T, metric pmetric.Metric, wasChecked *bool, isCumulative bool, boundaries []float64, value HistogramValue, startTs pcommon.Timestamp) {
	assert.Equal(t, value.count, metric.Histogram().DataPoints().At(0).Count())
	assert.Equal(t, value.sum, metric.Histogram().DataPoints().At(0).Sum())
	assert.Equal(t, value.max, metric.Histogram().DataPoints().At(0).Max())
	assert.Equal(t, value.min, metric.Histogram().DataPoints().At(0).Min())
	assert.Equal(t, isCumulative, metric.Histogram().AggregationTemporality() == pmetric.AggregationTemporalityCumulative)
	assert.Equal(t, startTs, metric.Histogram().DataPoints().At(0).StartTimestamp())

	assert.Equal(t, len(boundaries), metric.Histogram().DataPoints().At(0).ExplicitBounds().Len())
	for i := 0; i < len(boundaries); i++ {
		assert.Equal(t, boundaries[i], metric.Histogram().DataPoints().At(0).ExplicitBounds().At(i))
	}

	assert.Equal(t, len(value.buckets), metric.Histogram().DataPoints().At(0).BucketCounts().Len())
	for i := 0; i < len(value.buckets); i++ {
		assert.Equal(t, value.buckets[i], metric.Histogram().DataPoints().At(0).BucketCounts().At(i))
	}

	*wasChecked = true
}

type GaugeArg[T GaugeValue] struct {
	name    string
	startTS pcommon.Timestamp
	ts      pcommon.Timestamp
	values  []T
}

type CounterArg[T CounterValue] struct {
	name       string
	startTS    pcommon.Timestamp
	ts         pcommon.Timestamp
	cumulative bool
	monotonic  bool
	values     []T
}

type HistogramValue struct {
	count   uint64
	sum     float64
	max     float64
	min     float64
	buckets []uint64
}

type HistogramArg struct {
	name       string
	startTS    pcommon.Timestamp
	ts         pcommon.Timestamp
	cumulative bool
	boundaries []float64
	values     []HistogramValue
}

type ScopeArg struct {
	name    string
	version string

	doubleGauges []GaugeArg[float64]
	intGauges    []GaugeArg[int64]

	doubleCounters []CounterArg[float64]
	intCounters    []CounterArg[int64]

	histograms []HistogramArg
}

type ResourceMetricsArg struct {
	scopes []ScopeArg
}

type MetricArg struct {
	resourceMetrics []ResourceMetricsArg
}

func CreateArgument(metricArg MetricArg) pmetric.Metrics {
	var res pmetric.Metrics = pmetric.NewMetrics()

	for _, resourceMetricArg := range metricArg.resourceMetrics {
		resourceMetric := res.ResourceMetrics().AppendEmpty()
		for _, scopeMetricArg := range resourceMetricArg.scopes {
			scopeMetric := resourceMetric.ScopeMetrics().AppendEmpty()
			scopeMetric.Scope().SetName(scopeMetricArg.name)
			scopeMetric.Scope().SetVersion(scopeMetricArg.version)

			for _, metricArg := range scopeMetricArg.intGauges {
				metric := scopeMetric.Metrics().AppendEmpty()
				metric.SetName(metricArg.name)
				gauge := metric.SetEmptyGauge()
				for _, value := range metricArg.values {
					dp := gauge.DataPoints().AppendEmpty()
					dp.SetIntValue(value)
					dp.SetStartTimestamp(metricArg.startTS)
					dp.SetTimestamp(metricArg.ts)
				}
			}

			for _, metricArg := range scopeMetricArg.doubleGauges {
				metric := scopeMetric.Metrics().AppendEmpty()
				metric.SetName(metricArg.name)
				gauge := metric.SetEmptyGauge()
				for _, value := range metricArg.values {
					dp := gauge.DataPoints().AppendEmpty()
					dp.SetDoubleValue(value)
					dp.SetStartTimestamp(metricArg.startTS)
					dp.SetTimestamp(metricArg.ts)
				}
			}

			for _, metricArg := range scopeMetricArg.intCounters {
				metric := scopeMetric.Metrics().AppendEmpty()
				metric.SetName(metricArg.name)
				counter := metric.SetEmptySum()
				if metricArg.cumulative {
					counter.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
				} else {
					counter.SetAggregationTemporality(pmetric.AggregationTemporalityDelta)
				}
				counter.SetIsMonotonic(metricArg.monotonic)
				for _, value := range metricArg.values {
					dp := counter.DataPoints().AppendEmpty()
					dp.SetIntValue(value)
					dp.SetStartTimestamp(metricArg.startTS)
					dp.SetTimestamp(metricArg.ts)
				}
			}

			for _, metricArg := range scopeMetricArg.doubleCounters {
				metric := scopeMetric.Metrics().AppendEmpty()
				metric.SetName(metricArg.name)
				counter := metric.SetEmptySum()
				if metricArg.cumulative {
					counter.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
				} else {
					counter.SetAggregationTemporality(pmetric.AggregationTemporalityDelta)
				}
				counter.SetIsMonotonic(metricArg.monotonic)
				for _, value := range metricArg.values {
					dp := counter.DataPoints().AppendEmpty()
					dp.SetDoubleValue(value)
					dp.SetStartTimestamp(metricArg.startTS)
					dp.SetTimestamp(metricArg.ts)
				}
			}

			for _, metricArg := range scopeMetricArg.histograms {
				metric := scopeMetric.Metrics().AppendEmpty()
				metric.SetName(metricArg.name)
				histogram := metric.SetEmptyHistogram()
				if metricArg.cumulative {
					histogram.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
				} else {
					histogram.SetAggregationTemporality(pmetric.AggregationTemporalityDelta)
				}
				for _, value := range metricArg.values {
					dp := histogram.DataPoints().AppendEmpty()
					dp.SetCount(value.count)
					dp.SetSum(value.sum)
					dp.SetMax(value.max)
					dp.SetMin(value.min)
					dp.ExplicitBounds().Append(metricArg.boundaries...)
					dp.BucketCounts().Append(value.buckets...)
					dp.SetStartTimestamp(metricArg.startTS)
					dp.SetTimestamp(metricArg.ts)
				}
			}
		}
	}

	return res
}
