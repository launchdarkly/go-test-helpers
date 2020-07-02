# Change log

All notable changes to the project will be documented in this file. This project adheres to [Semantic Versioning](http://semver.org).

## [1.2.0] - 2020-07-01
### Added:
- New package `testbox` for running a Go test in a temporary environment.

## [1.1.3] - 2020-06-04
### Added:
- Added `go.mod` so this package can be consumed as a module. This does not affect code that is currently consuming it via `go get`, `dep`, or `govendor`.

## [1.1.2] - 2020-04-01
### Fixed:
- Patch event type for client-side streams.

## [1.1.1] - 2020-04-01
### Fixed:
- In `ldservices`, fixed JSON property names for simulated client-side flag data.

## [1.1.0] - 2020-04-01
### Added:
- Method `HandlerForPathRegex` in `httphelpers`.
- New methods and types in `ldservices` to simulate LaunchDarkly client-side streaming endpoints.

## [1.0.0] - 2020-03-16
Initial release.
