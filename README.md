# Reduce Resolution processor

This is a custom processor that reduces all received metrics within one time interval into a set of statics that describe what happened in that same period.

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