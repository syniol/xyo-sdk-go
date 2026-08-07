# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.1] - 2026-08-07
### Changed
- Updated repository URL and module path to `github.com/xyo-financial/sdk-go`.

## [1.0.0] - 2026-07-11
### Added
- Initial release of the Go SDK for XYO.Financial API.
- `NewClient` constructor accepting `ClientConfig` with API key.
- `EnrichTransaction` method for single payment transaction enrichment.
- `EnrichTransactionCollection` method for bulk enrichment requests.
- `EnrichTransactionCollectionStatus` method to query bulk enrichment status.
- Pluggable HTTP transport layer via `httpClient` for testability.
- Unit tests covering all enrichment endpoints (success and error paths).
- Docker-based CI pipeline with `golangci-lint` and `go test`.
- Example application in `example/` directory.
