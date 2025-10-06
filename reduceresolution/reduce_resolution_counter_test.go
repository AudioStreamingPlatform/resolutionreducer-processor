// Copyright (C) 2025 Bang & Olufsen A/S, Denmark
//
// SPDX-License-Identifier: MIT License

package reduceresolution

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

func TestValidateIntCounterAggregation(t *testing.T) {
	var processor *ReduceResolution = &ReduceResolution{
		Logger: &zap.Logger{},
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
							[]CounterArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									true,
									[]int64{
										3,
									},
								},
							},
						},
					},
				},
			},
		},
	)

	t.Run("validate single counter", func(t *testing.T) {
		finalMetrics, error := processor.ProcessMetrics(nil, mainMetrics)

		assert.NoError(t, error)
		var counter bool = false

		for i := 0; i < finalMetrics.ResourceMetrics().Len(); i++ {
			resourceMetric := finalMetrics.ResourceMetrics().At(i)
			assert.Equal(t, 1, resourceMetric.ScopeMetrics().Len())
			for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
				scope := resourceMetric.ScopeMetrics().At(j)
				assert.Equal(t, "testscope", scope.Scope().Name())
				for k := 0; k < scope.Metrics().Len(); k++ {
					metric := scope.Metrics().At(k)
					assert.Equal(t, pmetric.MetricTypeSum, metric.Type())
					assert.Equal(t, "testmetric", metric.Name())
					ValidateIntCounter(t, metric, &counter, true, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
				}
				assert.True(t, counter)
			}
		}
	})
}

func TestValidate2CounterAggregationDifferentScope(t *testing.T) {
	var processor *ReduceResolution = &ReduceResolution{
		Logger: &zap.Logger{},
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
							[]CounterArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									true,
									[]int64{
										3,
									},
								},
							},
						},
						{
							"testscope",
							"1.0",
							[]GaugeArg[float64]{},
							[]GaugeArg[int64]{},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 40, 0, time.UTC)),
									true,
									[]int64{
										5,
									},
								},
							},
						},
					},
				},
			},
		},
	)

	t.Run("validate aggregated counter", func(t *testing.T) {
		finalMetrics, error := processor.ProcessMetrics(nil, mainMetrics)

		assert.NoError(t, error)
		var counter bool = false

		for i := 0; i < finalMetrics.ResourceMetrics().Len(); i++ {
			resourceMetric := finalMetrics.ResourceMetrics().At(i)
			assert.Equal(t, 1, resourceMetric.ScopeMetrics().Len())
			for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
				scope := resourceMetric.ScopeMetrics().At(j)
				assert.Equal(t, "testscope", scope.Scope().Name())
				for k := 0; k < scope.Metrics().Len(); k++ {
					metric := scope.Metrics().At(k)
					assert.Equal(t, pmetric.MetricTypeSum, metric.Type())
					assert.Equal(t, "testmetric", metric.Name())
					ValidateIntCounter(t, metric, &counter, true, 5, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
				}
				assert.True(t, counter)
			}
		}
	})
}

func TestValidate2CounterAggregationSameScope(t *testing.T) {
	var processor *ReduceResolution = &ReduceResolution{
		Logger: &zap.Logger{},
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
							[]CounterArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									true,
									[]int64{
										3, 5,
									},
								},
							},
						},
					},
				},
			},
		},
	)

	t.Run("validate single counter", func(t *testing.T) {
		finalMetrics, error := processor.ProcessMetrics(nil, mainMetrics)

		assert.NoError(t, error)
		var counter bool = false

		for i := 0; i < finalMetrics.ResourceMetrics().Len(); i++ {
			resourceMetric := finalMetrics.ResourceMetrics().At(i)
			assert.Equal(t, 1, resourceMetric.ScopeMetrics().Len())
			for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
				scope := resourceMetric.ScopeMetrics().At(j)
				assert.Equal(t, "testscope", scope.Scope().Name())
				for k := 0; k < scope.Metrics().Len(); k++ {
					metric := scope.Metrics().At(k)
					assert.Equal(t, pmetric.MetricTypeSum, metric.Type())
					assert.Equal(t, "testmetric", metric.Name())
					ValidateIntCounter(t, metric, &counter, true, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
				}
				assert.True(t, counter)
			}
		}
	})
}

func TestValidate2DoubleCounterAggregationSameScope(t *testing.T) {
	var processor *ReduceResolution = &ReduceResolution{
		Logger: &zap.Logger{},
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
							[]CounterArg[float64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									true,
									[]float64{
										3.0, 5.0,
									},
								},
							},
							[]CounterArg[int64]{},
						},
					},
				},
			},
		},
	)

	t.Run("validate single counter", func(t *testing.T) {
		finalMetrics, error := processor.ProcessMetrics(nil, mainMetrics)

		assert.NoError(t, error)
		var counter bool = false

		for i := 0; i < finalMetrics.ResourceMetrics().Len(); i++ {
			resourceMetric := finalMetrics.ResourceMetrics().At(i)
			assert.Equal(t, 1, resourceMetric.ScopeMetrics().Len())
			for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
				scope := resourceMetric.ScopeMetrics().At(j)
				assert.Equal(t, "testscope", scope.Scope().Name())
				for k := 0; k < scope.Metrics().Len(); k++ {
					metric := scope.Metrics().At(k)
					assert.Equal(t, pmetric.MetricTypeSum, metric.Type())
					assert.Equal(t, "testmetric", metric.Name())
					ValidateDoubleCounter(t, metric, &counter, true, 3.0, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
				}
				assert.True(t, counter)
			}
		}
	})
}

func TestValidate2CounterAggregationDiffScopeDiffName(t *testing.T) {
	var processor *ReduceResolution = &ReduceResolution{
		Logger: &zap.Logger{},
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
							[]CounterArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									true,
									[]int64{
										3,
									},
								},
							},
						},
						{
							"testscope2",
							"1.0",
							[]GaugeArg[float64]{},
							[]GaugeArg[int64]{},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									true,
									[]int64{
										5,
									},
								},
							},
						},
					},
				},
			},
		},
	)

	t.Run("validate single counter in each scope", func(t *testing.T) {
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
					var counter bool = false
					for k := 0; k < scope.Metrics().Len(); k++ {
						metric := scope.Metrics().At(k)
						assert.Equal(t, pmetric.MetricTypeSum, metric.Type())
						assert.Equal(t, "testmetric", metric.Name())
						ValidateIntCounter(t, metric, &counter, true, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
					}
					assert.True(t, counter)

				case "testscope2":
					assert.False(t, scope2)
					scope2 = true
					var counter bool = false
					for k := 0; k < scope.Metrics().Len(); k++ {
						metric := scope.Metrics().At(k)
						assert.Equal(t, pmetric.MetricTypeSum, metric.Type())
						assert.Equal(t, "testmetric", metric.Name())
						ValidateIntCounter(t, metric, &counter, true, 5, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)))

					}
					assert.True(t, counter)
				}
			}
		}
		assert.True(t, scope1)
		assert.True(t, scope2)
	})
}

func TestValidate2CounterAggregationDiffScopeDiffVersion(t *testing.T) {
	var processor *ReduceResolution = &ReduceResolution{
		Logger: &zap.Logger{},
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
							[]CounterArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									true,
									[]int64{
										3,
									},
								},
							},
						},
						{
							"testscope",
							"2.0",
							[]GaugeArg[float64]{},
							[]GaugeArg[int64]{},
							[]CounterArg[float64]{},
							[]CounterArg[int64]{
								{
									"testmetric",
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)),
									pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 30, 0, time.UTC)),
									true,
									[]int64{
										5,
									},
								},
							},
						},
					},
				},
			},
		},
	)

	t.Run("validate single counter in diff scopes", func(t *testing.T) {
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
					var counter bool = false
					for k := 0; k < scope.Metrics().Len(); k++ {
						metric := scope.Metrics().At(k)
						assert.Equal(t, pmetric.MetricTypeSum, metric.Type())
						assert.Equal(t, "testmetric", metric.Name())
						ValidateIntCounter(t, metric, &counter, true, 3, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 20, 0, time.UTC)))
					}
					assert.True(t, counter)
				case "2.0":
					assert.Equal(t, "testscope", scope.Scope().Name())
					assert.False(t, scope2)
					scope2 = true
					var counter bool = false
					for k := 0; k < scope.Metrics().Len(); k++ {
						metric := scope.Metrics().At(k)
						assert.Equal(t, pmetric.MetricTypeSum, metric.Type())
						assert.Equal(t, "testmetric", metric.Name())
						ValidateIntCounter(t, metric, &counter, true, 5, pcommon.NewTimestampFromTime(time.Date(2025, time.January, 1, 12, 0, 10, 0, time.UTC)))
					}
					assert.True(t, counter)
				}
			}
		}
		assert.True(t, scope1)
		assert.True(t, scope2)
	})
}
