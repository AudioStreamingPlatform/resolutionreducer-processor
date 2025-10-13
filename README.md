# Reduce Resolution processor

This is a custom processor that reduces all received metrics within one time interval into a set of statics that describe what happened in that same period.

## How are metrics reduced

This processor grabs a group of received metrics, and summarizes each metric with a group of statistics. The processor does not batch by itself. Therefore, to aggregate metrics for a window of 30 seconds, for example, it is necessary to add a `batch` processor first with a timeout of 30 seconds before the processor. For example:

```yaml
...

processors:
  batch:
    timeout: 30s
  
  reduceresolution:

...

service:
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [batch, reduceresolution]
      exporters: [debug]
...
```

### Gauge
The end result of an aggregated gauge is the following metrics:
- _gauge_avg (the average value)
- _gauge_max (the maximum value found within the sample)
- _gauge_min (the minimum value found within the sample)
- _gauge_sum (the sum of all the values found within the sample)
- _gauge_count (the total number values found in the sample)

When a gauge is received, it is converted into 5 different gauges, represented above. Therefore, if a gauge called `speed` is emitted from a client, it will result in 5 gauges called `speed_gauge_avg`, `speed_gauge_max`, `speed_gauge_max`, `speed_gauge_min`, `speed_gauge_sum` and `speed_gauge_count`.

### Counter and UpDownCounter
Both the Counter and the UpDownCounter are just summed together and emitted with a single value. The name of the counter or the UpDownCounter are not changed.

### Linear Histogram
Just like the Counter, the Linear Histogram is summed together and emitted a single set of buckets. The name of the Histogram is not changed.

## How to work locally

Install the go version 1.20.12 locally.
Obtain the Open Telemetry Collector builder with `Dockerfile`.

## How to build a collector

To build a collector run the following command:

`./ocb --config manifest.yaml`

## How to run the built collector

To run the built collector, run the following command:

`./_build/opentelemetry-collector --config config.yaml`

## How to run the unit tests

Change to the `reduceresolution` folder and run `go test`. For more details on the test run, execute `go test -v`.