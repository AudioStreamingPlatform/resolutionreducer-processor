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

func TestValidateIntGaugeAggregation(t *testing.T) {
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
							[]GaugeArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)),
									[]int64{
										3,
									},
								},
							},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{},
							[]HistogramArg{},
						},
					},
				},
			},
		},
	)

	t.Run("validate single gauge", func(t *testing.T) {
		finalMetrics, error := processor.ProcessMetrics(nil, mainMetrics)

		assert.NoError(t, error)
		var count bool = false
		var max bool = false
		var min bool = false
		var avg bool = false

		for i := 0; i < finalMetrics.ResourceMetrics().Len(); i++ {
			resourceMetric := finalMetrics.ResourceMetrics().At(i)
			assert.Equal(t, 1, resourceMetric.ScopeMetrics().Len())
			for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
				scope := resourceMetric.ScopeMetrics().At(j)
				assert.Equal(t, "testscope", scope.Scope().Name())
				for k := 0; k < scope.Metrics().Len(); k++ {
					metric := scope.Metrics().At(k)
					assert.Equal(t, pmetric.MetricTypeGauge, metric.Type())
					switch metric.Name() {
					case "testmetric_gauge_count":
						ValidateIntGauge(t, metric, &count, 1, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)))
					case "testmetric_gauge_max":
						ValidateIntGauge(t, metric, &max, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)))
					case "testmetric_gauge_min":
						ValidateIntGauge(t, metric, &min, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)))
					case "testmetric_gauge_avg":
						ValidateIntGauge(t, metric, &avg, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)))
					}
				}
				assert.True(t, count)
				assert.True(t, max)
				assert.True(t, min)
				assert.True(t, avg)
			}
		}
	})
}

func TestValidate2GaugeAggregationDifferentScope(t *testing.T) {
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
							[]GaugeArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)),
									[]int64{
										3,
									},
								},
							},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{},
							[]HistogramArg{},
						},
						{
							"testscope",
							"1.0",
							[]GaugeArg[float64]{},
							[]GaugeArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									[]int64{
										5,
									},
								},
							},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{},
							[]HistogramArg{},
						},
					},
				},
			},
		},
	)

	t.Run("validate aggregated separated gauge", func(t *testing.T) {
		finalMetrics, error := processor.ProcessMetrics(nil, mainMetrics)

		assert.NoError(t, error)
		var count bool = false
		var max bool = false
		var min bool = false
		var avg bool = false

		for i := 0; i < finalMetrics.ResourceMetrics().Len(); i++ {
			resourceMetric := finalMetrics.ResourceMetrics().At(i)
			assert.Equal(t, 1, resourceMetric.ScopeMetrics().Len())
			for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
				scope := resourceMetric.ScopeMetrics().At(j)
				assert.Equal(t, "testscope", scope.Scope().Name())
				for k := 0; k < scope.Metrics().Len(); k++ {
					metric := scope.Metrics().At(k)
					assert.Equal(t, pmetric.MetricTypeGauge, metric.Type())
					switch metric.Name() {
					case "testmetric_gauge_count":
						ValidateIntGauge(t, metric, &count, 2, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)))
					case "testmetric_gauge_max":
						ValidateIntGauge(t, metric, &max, 5, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)))
					case "testmetric_gauge_min":
						ValidateIntGauge(t, metric, &min, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)))
					case "testmetric_gauge_avg":
						ValidateIntGauge(t, metric, &avg, 4, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)))
					}
				}
				assert.True(t, count)
				assert.True(t, max)
				assert.True(t, min)
				assert.True(t, avg)
			}
		}
	})
}

func TestValidate2GaugeAggregationSameScope(t *testing.T) {
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
							[]GaugeArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									[]int64{
										3, 5,
									},
								},
							},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{},
							[]HistogramArg{},
						},
					},
				},
			},
		},
	)

	t.Run("validate aggregated gauge", func(t *testing.T) {
		finalMetrics, error := processor.ProcessMetrics(nil, mainMetrics)

		assert.NoError(t, error)
		var count bool = false
		var max bool = false
		var min bool = false
		var avg bool = false

		for i := 0; i < finalMetrics.ResourceMetrics().Len(); i++ {
			resourceMetric := finalMetrics.ResourceMetrics().At(i)
			assert.Equal(t, 1, resourceMetric.ScopeMetrics().Len())
			for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
				scope := resourceMetric.ScopeMetrics().At(j)
				assert.Equal(t, "testscope", scope.Scope().Name())
				assert.Equal(t, "1.0", scope.Scope().Version())
				for k := 0; k < scope.Metrics().Len(); k++ {
					metric := scope.Metrics().At(k)
					assert.Equal(t, pmetric.MetricTypeGauge, metric.Type())
					switch metric.Name() {
					case "testmetric_gauge_count":
						ValidateIntGauge(t, metric, &count, 2, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
					case "testmetric_gauge_max":
						ValidateIntGauge(t, metric, &max, 5, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
					case "testmetric_gauge_min":
						ValidateIntGauge(t, metric, &min, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
					case "testmetric_gauge_avg":
						ValidateIntGauge(t, metric, &avg, 4, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
					}
				}
				assert.True(t, count)
				assert.True(t, max)
				assert.True(t, min)
				assert.True(t, avg)
			}
		}
	})
}

func TestValidate2GaugeDoubleAggregationSameScope(t *testing.T) {
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
							[]GaugeArg[float64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									[]float64{
										3, 5,
									},
								},
							},
							[]GaugeArg[int64]{},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{},
							[]HistogramArg{},
						},
					},
				},
			},
		},
	)

	t.Run("validate aggregated gauge", func(t *testing.T) {
		finalMetrics, error := processor.ProcessMetrics(nil, mainMetrics)

		assert.NoError(t, error)
		var count bool = false
		var max bool = false
		var min bool = false
		var avg bool = false

		for i := 0; i < finalMetrics.ResourceMetrics().Len(); i++ {
			resourceMetric := finalMetrics.ResourceMetrics().At(i)
			assert.Equal(t, 1, resourceMetric.ScopeMetrics().Len())
			for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
				scope := resourceMetric.ScopeMetrics().At(j)
				assert.Equal(t, "testscope", scope.Scope().Name())
				assert.Equal(t, "1.0", scope.Scope().Version())
				for k := 0; k < scope.Metrics().Len(); k++ {
					metric := scope.Metrics().At(k)
					assert.Equal(t, pmetric.MetricTypeGauge, metric.Type())
					switch metric.Name() {
					case "testmetric_gauge_count":
						ValidateIntGauge(t, metric, &count, 2, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
					case "testmetric_gauge_max":
						ValidateDoubleGauge(t, metric, &max, 5.0, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
					case "testmetric_gauge_min":
						ValidateDoubleGauge(t, metric, &min, 3.0, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
					case "testmetric_gauge_avg":
						ValidateDoubleGauge(t, metric, &avg, 4.0, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
					}
				}
				assert.True(t, count)
				assert.True(t, max)
				assert.True(t, min)
				assert.True(t, avg)
			}
		}
	})
}

func TestValidate2GaugeAggregationDiffScopeDiffName(t *testing.T) {
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
							[]GaugeArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									[]int64{
										3,
									},
								},
							},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{},
							[]HistogramArg{},
						},
						{
							"testscope2",
							"1.0",
							[]GaugeArg[float64]{},
							[]GaugeArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									[]int64{
										5,
									},
								},
							},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{},
							[]HistogramArg{},
						},
					},
				},
			},
		},
	)

	t.Run("validate double gauge in different scopes", func(t *testing.T) {
		finalMetrics, error := processor.ProcessMetrics(nil, mainMetrics)

		assert.NoError(t, error)
		var scope1 bool = false
		var scope2 bool = false
		for i := 0; i < finalMetrics.ResourceMetrics().Len(); i++ {
			resourceMetric := finalMetrics.ResourceMetrics().At(i)
			assert.Equal(t, 2, resourceMetric.ScopeMetrics().Len())
			for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
				scope := resourceMetric.ScopeMetrics().At(j)
				switch scope.Scope().Name() {
				case "testscope":
					assert.False(t, scope1)
					scope1 = true
					var count bool = false
					var max bool = false
					var min bool = false
					var avg bool = false
					for k := 0; k < scope.Metrics().Len(); k++ {
						metric := scope.Metrics().At(k)
						assert.Equal(t, pmetric.MetricTypeGauge, metric.Type())
						switch metric.Name() {
						case "testmetric_gauge_count":
							ValidateIntGauge(t, metric, &count, 1, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)))
						case "testmetric_gauge_max":
							ValidateIntGauge(t, metric, &max, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)))
						case "testmetric_gauge_min":
							ValidateIntGauge(t, metric, &min, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)))
						case "testmetric_gauge_avg":
							ValidateIntGauge(t, metric, &avg, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)))
						}
					}
					assert.True(t, count)
					assert.True(t, max)
					assert.True(t, min)
					assert.True(t, avg)
				case "testscope2":
					assert.False(t, scope2)
					scope2 = true
					var count bool = false
					var max bool = false
					var min bool = false
					var avg bool = false
					for k := 0; k < scope.Metrics().Len(); k++ {
						metric := scope.Metrics().At(k)
						assert.Equal(t, pmetric.MetricTypeGauge, metric.Type())
						switch metric.Name() {
						case "testmetric_gauge_count":
							ValidateIntGauge(t, metric, &count, 1, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
						case "testmetric_gauge_max":
							ValidateIntGauge(t, metric, &max, 5, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
						case "testmetric_gauge_min":
							ValidateIntGauge(t, metric, &min, 5, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
						case "testmetric_gauge_avg":
							ValidateIntGauge(t, metric, &avg, 5, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
						}
					}
					assert.True(t, count)
					assert.True(t, max)
					assert.True(t, min)
					assert.True(t, avg)
				}
			}
		}
		assert.True(t, scope1)
		assert.True(t, scope2)
	})
}

func TestValidate2GaugeAggregationDiffScopeDiffVersion(t *testing.T) {
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
							[]GaugeArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									[]int64{
										3,
									},
								},
							},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{},
							[]HistogramArg{},
						},
						{
							"testscope",
							"2.0",
							[]GaugeArg[float64]{},
							[]GaugeArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									[]int64{
										5,
									},
								},
							},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{},
							[]HistogramArg{},
						},
					},
				},
			},
		},
	)

	t.Run("validate double gauge with same scope diff versions", func(t *testing.T) {
		finalMetrics, error := processor.ProcessMetrics(nil, mainMetrics)

		assert.NoError(t, error)
		var scope1 bool = false
		var scope2 bool = false
		for i := 0; i < finalMetrics.ResourceMetrics().Len(); i++ {
			resourceMetric := finalMetrics.ResourceMetrics().At(i)
			assert.Equal(t, 2, resourceMetric.ScopeMetrics().Len())
			for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
				scope := resourceMetric.ScopeMetrics().At(j)
				switch scope.Scope().Version() {
				case "1.0":
					assert.Equal(t, "testscope", scope.Scope().Name())
					assert.False(t, scope1)
					scope1 = true
					var count bool = false
					var max bool = false
					var min bool = false
					var avg bool = false
					for k := 0; k < scope.Metrics().Len(); k++ {
						metric := scope.Metrics().At(k)
						assert.Equal(t, pmetric.MetricTypeGauge, metric.Type())
						switch metric.Name() {
						case "testmetric_gauge_count":
							ValidateIntGauge(t, metric, &count, 1, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
						case "testmetric_gauge_max":
							ValidateIntGauge(t, metric, &max, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
						case "testmetric_gauge_min":
							ValidateIntGauge(t, metric, &min, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
						case "testmetric_gauge_avg":
							ValidateIntGauge(t, metric, &avg, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
						}
					}
					assert.True(t, count)
					assert.True(t, max)
					assert.True(t, min)
					assert.True(t, avg)
				case "2.0":
					assert.Equal(t, "testscope", scope.Scope().Name())
					assert.False(t, scope2)
					scope2 = true
					var count bool = false
					var max bool = false
					var min bool = false
					var avg bool = false
					for k := 0; k < scope.Metrics().Len(); k++ {
						metric := scope.Metrics().At(k)
						assert.Equal(t, pmetric.MetricTypeGauge, metric.Type())
						switch metric.Name() {
						case "testmetric_gauge_count":
							ValidateIntGauge(t, metric, &count, 1, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)))
						case "testmetric_gauge_max":
							ValidateIntGauge(t, metric, &max, 5, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)))
						case "testmetric_gauge_min":
							ValidateIntGauge(t, metric, &min, 5, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)))
						case "testmetric_gauge_avg":
							ValidateIntGauge(t, metric, &avg, 5, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)))
						}
					}
					assert.True(t, count)
					assert.True(t, max)
					assert.True(t, min)
					assert.True(t, avg)
				}
			}
		}
		assert.True(t, scope1)
		assert.True(t, scope2)
	})
}
