# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.3.0] - 2021-02-03

### Changed

- Renamed host and guest respectively to management and workload.

## [1.2.0] - 2020-12-18

### Removed

- Removed high cardinality `nginx-ingress-controller` metrics from managed apps job.

## [1.1.2] - 2020-12-07

### Added

- Add `provider` label to StaticConfigs for etcd metrics.

## [1.1.1] - 2020-11-16

### Added

- Add `provider` label to tenant cluster etcd metrics.

## [1.1.0] - 2020-11-05

### Added

- Add `provider` label to tenant cluster metrics.

## [1.0.2] - 2020-10-12

- Drop monitoring of non-managed kiam pods.

## [0.2.0] - 2020-09-02

- Add kube-proxy metrics

## [0.1.0] - 2020-07-14

### Added

- This CHANGELOG file
- Tagging first version

[unreleased]: https://github.com/giantswarm/prometheus-config-controller/compare/v1.3.0...HEAD
[1.3.0]: https://github.com/giantswarm/prometheus-config-controllera/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/giantswarm/prometheus-config-controllera/compare/v1.1.2...v1.2.0
[1.1.2]: https://github.com/giantswarm/prometheus-config-controllera/compare/v1.1.1...v1.1.2
[1.1.1]: https://github.com/giantswarm/prometheus-config-controllera/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/giantswarm/prometheus-config-controllera/compare/v1.0.2...v1.1.0
[1.0.2]: https://github.com/giantswarm/prometheus-config-controllera/compare/v0.2.0...v1.0.2
[0.2.0]: https://github.com/giantswarm/prometheus-config-controller/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/giantswarm/prometheus-config-controller/tag/v0.1.0
