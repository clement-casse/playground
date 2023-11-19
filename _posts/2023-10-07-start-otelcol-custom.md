---
layout: single
categories: [ otelcol-custom ]
title: "Building a Custom OpenTelemetry Collector with Nix - A Starting Point"
tags: [ opentelemetry, nix-flake ]
---

In this post, I'll present how to combine *Nix Flakes* with *OpenTelemetry Collector Builder* to create a custom Collector binary only with the components you need.
<!--more-->

# Context and Motivation

## The Big Picture of OpenTelemetry Collector

OpenTelemetry is a big open source project where multiple developers contribute on a daily basis.
It is one of the [CNCF incupating projects](https://www.cncf.io/projects/opentelemetry/), 
Unlike [Prometheus](https://www.cncf.io/projects/prometheus/) or [Jaeger](https://www.cncf.io/projects/jaeger/) which propose solutions to store and query respectively metrics and traces emited by software in Cloud environments, OpenTelemetry takes a complementary approach by proposing a standardized approach for collecting and processing monitoring data before sending them to a particular backend.

Indeed, OpenTelemetry describes itself as:

> OpenTelemetry is a collection of APIs, SDKs, and tools.
> Use it to instrument, generate, collect, and export telemetry data (metrics, logs, and traces) to help you analyze your software’s performance and behavior.

And, among the *tools* category, lies the *OpenTelemetry Collector*: a standalone process described as a vendor-agnostic way to receive, process and export telemetry data.
This *collector* materializes as different kind of modules chained together to form a processing pipeline for the different monitoring signals (traces, metrics, logs):

- `receivers` are modules that ingest data in the pipeline, their role is to transpose some specific format in the OpenTelemetry format so they can be processed by further stages in the pipeline.
- `processors` are modules that alter telemetry data that already is in the OpenTelemetry format, this alteration can be like adding/removing some attribute to the data, or simply to propagate the data further in the pipeline or not ...
- `exporters` are modules that convert OpenTelemetry data to an other format to send data to a custom backend like, e.g. Jaeger, Amazon CloudWatch, Prometheus, a file in the JSON or CSV format, ...

Other modules exists like `extensions` and `connectors` which will not be detailed in this post but, for the scope of what is presented here, have a similar behaviour.

The number of modules for OpenTelemetry Collector is high and ever increasing as demonstrates the GitHub repository [`open-telemetry/opentelemetry-collector-contrib`](https://github.com/open-telemetry/opentelemetry-collector-contrib) which lists modules contributions from the community.

## 