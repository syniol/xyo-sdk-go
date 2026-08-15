# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-08-15
### Added
- Migrated module path to `github.com/xyo-financial/sdk-go/v2`.
- Full OpenAPI-driven client architecture with strongly-typed interfaces.
- Structured RFC 7807 problem details error handling via `ErrorResponse` and `APIError`, supporting `errors.As` unwrapping.
- In-memory streaming tarball decompression for bulk collection downloads (`DownloadEnrichmentCollection`).
- Decompression bomb mitigations: aggregate stream byte bounds (`DefaultMaxArchiveBytes = 100MiB`), per-entry size limits (`DefaultMaxEntryBytes = 10MiB`), and entry count bounds (`DefaultMaxTarEntries = 50,000`).
- Strict URL scheme validation (`http`/`https`) for bulk download endpoints.
- Authorization header cross-host leak prevention on downloads.
- Automatic secret token redaction in OpenAPI client debug logging (`Authorization`, `Proxy-Authorization`, `Cookie`).
- Bounded response body reading across all OpenAPI endpoints (`maxResponseBytes = 32MiB`).
- Enterprise connection-pooled transport (`DefaultEnterpriseTransport` / `NewDefaultHTTPClient`) with `MaxIdleConnsPerHost = 100` to prevent TCP `TIME_WAIT` socket exhaustion under heavy concurrent load.
- Dedicated SDK version constants in `version.go` (`Version`, `DefaultUserAgent`).
- Modern Go 1.22 minimum toolchain support across `go.mod`, CI workflows, Docker, and examples.

### Changed
- Removed `decoder.DisallowUnknownFields()` across models to guarantee forward-compatible JSON decoding for API additive schema updates.
- Upgraded CI pipelines to Go 1.22 with SLSA attestations and SBOM generation.
- Expanded static analysis suite with `gosec` and `bodyclose` linters.
- Enhanced polling example documentation with exponential backoff and `context.WithTimeout` deadline cancellation.

## [1.2.4] - 2026-08-07
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
