// Copyright (C) 2025 Bang & Olufsen A/S, Denmark
//
// SPDX-License-Identifier: GPL-2.0-or-later

package reduceresolution

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

func TestValidateHistogramAggregation(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	var processor *ReduceResolution = &ReduceResolution{
		Logger: logger,
	}

	var mainMetrics pmetric.Metrics = CreateArgument(
		MetricArg{
			[]ResourceMetricsArg{
				{
					[]ScopeArg{
						{
							"testscope",
							"1.0",
							[]GaugeArg[float64]{},
							[]GaugeArg[int64]{},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{},
							[]HistogramArg{
								{
									"testhistogram",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									true,
									[]float64{0.0, 5.0, 10.0},
									[]HistogramValue{
										{
											1, 2.0, 2.0, 2.0, []uint64{0, 1, 0, 0},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)

	t.Run("validate single histogram value", func(t *testing.T) {
		finalMetrics, error := processor.ProcessMetrics(nil, mainMetrics)

		assert.NoError(t, error)
		var histogram bool = false

		for i := 0; i < finalMetrics.ResourceMetrics().Len(); i++ {
			resourceMetric := finalMetrics.ResourceMetrics().At(i)
			assert.Equal(t, 1, resourceMetric.ScopeMetrics().Len())
			for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
				scope := resourceMetric.ScopeMetrics().At(j)
				assert.Equal(t, "testscope", scope.Scope().Name())
				for k := 0; k < scope.Metrics().Len(); k++ {
					metric := scope.Metrics().At(k)
					assert.Equal(t, pmetric.MetricTypeHistogram, metric.Type())
					assert.Equal(t, "testhistogram", metric.Name())
					ValidateHistogram(t, metric, &histogram, true, []float64{0.0, 5.0, 10.0}, HistogramValue{1, 2.0, 2.0, 2.0, []uint64{0, 1, 0, 0}}, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
				}
				assert.True(t, histogram)
			}
		}
	})
}

func TestValidateHistogramSameScopeAggregationCumulative(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	var processor *ReduceResolution = &ReduceResolution{
		Logger: logger,
	}

	var mainMetrics pmetric.Metrics = CreateArgument(
		MetricArg{
			[]ResourceMetricsArg{
				{
					[]ScopeArg{
						{
							"testscope",
							"1.0",
							[]GaugeArg[float64]{},
							[]GaugeArg[int64]{},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{},
							[]HistogramArg{
								{
									"testhistogram",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									true,
									[]float64{0.0, 5.0, 10.0},
									[]HistogramValue{
										{
											1, 2.0, 2.0, 2.0, []uint64{0, 1, 0, 0},
										},
									},
								},
								{
									"testhistogram",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 40, 0, time.UTC)),
									true,
									[]float64{0.0, 5.0, 10.0},
									[]HistogramValue{
										{
											2, 6.0, 4.0, 2.0, []uint64{0, 2, 0, 0},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)

	t.Run("validate double histogram value", func(t *testing.T) {
		finalMetrics, error := processor.ProcessMetrics(nil, mainMetrics)

		assert.NoError(t, error)
		var histogram bool = false

		for i := 0; i < finalMetrics.ResourceMetrics().Len(); i++ {
			resourceMetric := finalMetrics.ResourceMetrics().At(i)
			assert.Equal(t, 1, resourceMetric.ScopeMetrics().Len())
			for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
				scope := resourceMetric.ScopeMetrics().At(j)
				assert.Equal(t, "testscope", scope.Scope().Name())
				for k := 0; k < scope.Metrics().Len(); k++ {
					metric := scope.Metrics().At(k)
					assert.Equal(t, pmetric.MetricTypeHistogram, metric.Type())
					assert.Equal(t, "testhistogram", metric.Name())
					ValidateHistogram(t, metric, &histogram, true, []float64{0.0, 5.0, 10.0}, HistogramValue{2, 6.0, 4.0, 2.0, []uint64{0, 2, 0, 0}}, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
				}
				assert.True(t, histogram)
			}
		}
	})
}

func TestValidateHistogramSameScopeAggregationDelta(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	var processor *ReduceResolution = &ReduceResolution{
		Logger: logger,
	}

	var mainMetrics pmetric.Metrics = CreateArgument(
		MetricArg{
			[]ResourceMetricsArg{
				{
					[]ScopeArg{
						{
							"testscope",
							"1.0",
							[]GaugeArg[float64]{},
							[]GaugeArg[int64]{},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{},
							[]HistogramArg{
								{
									"testhistogram",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									false,
									[]float64{0.0, 5.0, 10.0},
									[]HistogramValue{
										{
											1, 2.0, 2.0, 2.0, []uint64{0, 1, 0, 0},
										},
									},
								},
								{
									"testhistogram",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 40, 0, time.UTC)),
									false,
									[]float64{0.0, 5.0, 10.0},
									[]HistogramValue{
										{
											1, 4.0, 4.0, 4.0, []uint64{0, 1, 0, 0},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)

	t.Run("validate double histogram value", func(t *testing.T) {
		finalMetrics, error := processor.ProcessMetrics(nil, mainMetrics)

		assert.NoError(t, error)
		var histogram bool = false

		for i := 0; i < finalMetrics.ResourceMetrics().Len(); i++ {
			resourceMetric := finalMetrics.ResourceMetrics().At(i)
			assert.Equal(t, 1, resourceMetric.ScopeMetrics().Len())
			for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
				scope := resourceMetric.ScopeMetrics().At(j)
				assert.Equal(t, "testscope", scope.Scope().Name())
				for k := 0; k < scope.Metrics().Len(); k++ {
					metric := scope.Metrics().At(k)
					assert.Equal(t, pmetric.MetricTypeHistogram, metric.Type())
					assert.Equal(t, "testhistogram", metric.Name())
					ValidateHistogram(t, metric, &histogram, false, []float64{0.0, 5.0, 10.0}, HistogramValue{2, 6.0, 4.0, 2.0, []uint64{0, 2, 0, 0}}, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
				}
				assert.True(t, histogram)
			}
		}
	})
}

func TestValidateHistogramSameScopeAggregationWithDifferentBounds(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	var processor *ReduceResolution = &ReduceResolution{
		Logger: logger,
	}

	var mainMetrics pmetric.Metrics = CreateArgument(
		MetricArg{
			[]ResourceMetricsArg{
				{
					[]ScopeArg{
						{
							"testscope",
							"1.0",
							[]GaugeArg[float64]{},
							[]GaugeArg[int64]{},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{},
							[]HistogramArg{
								{
									"testhistogram",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									false,
									[]float64{0.0, 5.0, 10.0},
									[]HistogramValue{
										{
											1, 2.0, 2.0, 2.0, []uint64{0, 1, 0, 0},
										},
									},
								},
								{
									"testhistogram",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 40, 0, time.UTC)),
									false,
									[]float64{0.0, 5.0},
									[]HistogramValue{
										{
											1, 4.0, 4.0, 4.0, []uint64{0, 1, 0, 0},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)

	t.Run("validate double histogram value where second metric is ignored", func(t *testing.T) {
		finalMetrics, error := processor.ProcessMetrics(nil, mainMetrics)

		assert.NoError(t, error)
		var histogram bool = false

		for i := 0; i < finalMetrics.ResourceMetrics().Len(); i++ {
			resourceMetric := finalMetrics.ResourceMetrics().At(i)
			assert.Equal(t, 1, resourceMetric.ScopeMetrics().Len())
			for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
				scope := resourceMetric.ScopeMetrics().At(j)
				assert.Equal(t, "testscope", scope.Scope().Name())
				for k := 0; k < scope.Metrics().Len(); k++ {
					metric := scope.Metrics().At(k)
					assert.Equal(t, pmetric.MetricTypeHistogram, metric.Type())
					assert.Equal(t, "testhistogram", metric.Name())
					ValidateHistogram(t, metric, &histogram, false, []float64{0.0, 5.0, 10.0}, HistogramValue{1, 2.0, 2.0, 2.0, []uint64{0, 1, 0, 0}}, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
				}
				assert.True(t, histogram)
			}
		}
	})
}
