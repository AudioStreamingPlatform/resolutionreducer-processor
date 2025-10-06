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

func ValidateCounter(t *testing.T, metric pmetric.Metric, wasChecked *bool, isCumulative bool, startTs pcommon.Timestamp) {
	assert.False(t, *wasChecked)
	*wasChecked = true
	assert.Equal(t, pmetric.MetricTypeSum, metric.Type())
	assert.Equal(t, 1, metric.Sum().DataPoints().Len())
	if isCumulative {
		assert.Equal(t, pmetric.AggregationTemporalityCumulative, metric.Sum().AggregationTemporality())
	} else {
		assert.Equal(t, pmetric.AggregationTemporalityDelta, metric.Sum().AggregationTemporality())
	}
	assert.Equal(t, startTs, metric.Sum().DataPoints().At(0).StartTimestamp())
}

func ValidateDoubleCounter(t *testing.T, metric pmetric.Metric, wasChecked *bool, isCumulative bool, value float64, startTs pcommon.Timestamp) {
	ValidateCounter(t, metric, wasChecked, isCumulative, startTs)
	assert.Equal(t, pmetric.NumberDataPointValueTypeDouble, metric.Sum().DataPoints().At(0).ValueType())
	assert.Equal(t, value, metric.Sum().DataPoints().At(0).DoubleValue())
}

func ValidateIntCounter(t *testing.T, metric pmetric.Metric, wasChecked *bool, isCumulative bool, value int64, startTs pcommon.Timestamp) {
	ValidateCounter(t, metric, wasChecked, isCumulative, startTs)
	assert.Equal(t, pmetric.NumberDataPointValueTypeInt, metric.Sum().DataPoints().At(0).ValueType())
	assert.Equal(t, value, metric.Sum().DataPoints().At(0).IntValue())
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
	values     []T
}

type ScopeArg struct {
	name    string
	version string

	doubleGauges []GaugeArg[float64]
	intGauges    []GaugeArg[int64]

	doubleCounters []CounterArg[float64]
	intCounters    []CounterArg[int64]
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
				for _, value := range metricArg.values {
					dp := counter.DataPoints().AppendEmpty()
					dp.SetDoubleValue(value)
					dp.SetStartTimestamp(metricArg.startTS)
					dp.SetTimestamp(metricArg.ts)
				}
			}
		}
	}

	return res
}
