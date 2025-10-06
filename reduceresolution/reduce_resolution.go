// Copyright (C) 2025 Bang & Olufsen A/S, Denmark
//
// SPDX-License-Identifier: GPL-2.0-or-later

package reduceresolution

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type ReduceResolution struct {
	Logger *zap.Logger
}

// ProcessMetrics logs information about incoming metrics
func (p *ReduceResolution) ProcessMetrics(_ context.Context, metrics pmetric.Metrics) (pmetric.Metrics, error) {
	if metrics.ResourceMetrics().Len() == 0 {
		return metrics, nil
	}

	var aggregationTimeStamp pcommon.Timestamp = pcommon.NewTimestampFromTime(time.Now())

	var scopesMaps map[string]*ScopeContainer = make(map[string]*ScopeContainer)
	for i := 0; i < metrics.ResourceMetrics().Len(); i++ {
		resourceMetric := metrics.ResourceMetrics().At(i)
		for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
			scopeMetric := resourceMetric.ScopeMetrics().At(j)
			scopeKey := CreateScopeKey(scopeMetric)
			scopeContainer, ok := scopesMaps[scopeKey]

			if !ok {
				scopeContainer = CreateScopeContainer(scopeMetric)
				scopesMaps[scopeKey] = scopeContainer
			}

			for k := 0; k < scopeMetric.Metrics().Len(); k++ {
				metric := scopeMetric.Metrics().At(k)
				switch metric.Type() {
				// Deal with all gauges
				case pmetric.MetricTypeGauge:
					for l := 0; l < metric.Gauge().DataPoints().Len(); l++ {
						gauge := metric.Gauge().DataPoints().At(l)
						key := CreateMetricKey(metric, gauge.Attributes())
						if gauge.ValueType() == pmetric.NumberDataPointValueTypeInt {
							metricAggregate, ok := scopeContainer.intGaugeAggregate[key]
							if !ok {
								scopeContainer.intGaugeAggregate[key] = CreateGaugeAggregate(metric, gauge.Attributes(), gauge.StartTimestamp(), gauge.IntValue())
							} else {
								AggregateGauge(metricAggregate, gauge.StartTimestamp(), gauge.IntValue())
							}
						} else if gauge.ValueType() == pmetric.NumberDataPointValueTypeDouble {
							metricAggregate, ok := scopeContainer.floatGaugeAggregate[key]
							if !ok {
								scopeContainer.floatGaugeAggregate[key] = CreateGaugeAggregate(metric, gauge.Attributes(), gauge.StartTimestamp(), gauge.DoubleValue())
							} else {
								AggregateGauge(metricAggregate, gauge.StartTimestamp(), gauge.DoubleValue())
							}
						}
					}

				// Deal with all the counters
				case pmetric.MetricTypeSum:
					for l := 0; l < metric.Sum().DataPoints().Len(); l++ {
						counter := metric.Sum().DataPoints().At(l)
						key := CreateMetricKey(metric, counter.Attributes())

						if counter.ValueType() == pmetric.NumberDataPointValueTypeInt {
							metricAggregate, ok := scopeContainer.intCounterAggregate[key]
							if !ok {
								scopeContainer.intCounterAggregate[key] = CreateCounterAggregate(metric, counter.Attributes(), counter.StartTimestamp(), counter.Timestamp(), counter.IntValue())
							} else {
								AggregateCounter(metricAggregate, counter.StartTimestamp(), counter.Timestamp(), counter.IntValue())
							}
						} else if counter.ValueType() == pmetric.NumberDataPointValueTypeDouble {
							metricAggregate, ok := scopeContainer.floatCounterAggregate[key]
							if !ok {
								scopeContainer.floatCounterAggregate[key] = CreateCounterAggregate(metric, counter.Attributes(), counter.StartTimestamp(), counter.Timestamp(), counter.DoubleValue())
							} else {
								AggregateCounter(metricAggregate, counter.StartTimestamp(), counter.Timestamp(), counter.DoubleValue())
							}
						}
					}

				// For any non implemented metrics
				default:
					scopeContainer.leftoverMetric = append(scopeContainer.leftoverMetric, metric)
				}
			}
		}

	}

	firstResourceMetric := metrics.ResourceMetrics().At(0)
	metrics.ResourceMetrics().RemoveIf(func(m pmetric.ResourceMetrics) bool { return true })

	finalResourceMetric := metrics.ResourceMetrics().AppendEmpty()
	firstResourceMetric.Resource().CopyTo(finalResourceMetric.Resource())
	finalResourceMetric.SetSchemaUrl(firstResourceMetric.SchemaUrl())

	for _, scopeContainer := range scopesMaps {
		scope := finalResourceMetric.ScopeMetrics().AppendEmpty()
		scope.Scope().SetName(scopeContainer.scopeName)
		scope.Scope().SetVersion(scopeContainer.scopeVersion)

		for _, metricAggregate := range scopeContainer.intGaugeAggregate {
			CreateGaugeMetrics(scope, metricAggregate, aggregationTimeStamp)
		}
		for _, metricAggregate := range scopeContainer.floatGaugeAggregate {
			CreateGaugeMetrics(scope, metricAggregate, aggregationTimeStamp)
		}

		for _, metricAggregate := range scopeContainer.intCounterAggregate {
			CreateCounterMetrics(scope, metricAggregate, aggregationTimeStamp)
		}
		for _, metricAggregate := range scopeContainer.floatCounterAggregate {
			CreateCounterMetrics(scope, metricAggregate, aggregationTimeStamp)
		}
		for _, metric := range scopeContainer.leftoverMetric {
			metric.MoveTo(scope.Metrics().AppendEmpty())
		}
	}

	return metrics, nil
}
