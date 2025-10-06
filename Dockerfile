# Copyright (C) 2025 Bang & Olufsen A/S, Denmark
#
# SPDX-License-Identifier: GPL-2.0-or-later

FROM golang:1.20.12

WORKDIR /app

# from version 0.88
RUN git clone https://github.com/open-telemetry/opentelemetry-collector && \
cd opentelemetry-collector && git checkout d42d7e80974b3ba9fd1e235065638fe6a4d5455e && \
cd cmd/builder && go build && mv builder ocb && mv ocb /app/ocb
RUN rm -r /app/opentelemetry-collector/
