# 📝 Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Exported sentinel validation errors (`ErrNilRequest`, `ErrEmptyContent`, `ErrContentTooLong`, `ErrEmptyCountryCode`, `ErrInvalidCountryCode`) enabling `errors.Is` matching.

### Fixed
- Resolved connection leak in `Client.Close()` when a custom `*http.Client` with `nil` Transport is provided.
- Enforced strict 2-letter ISO 3166-1 alpha-2 validation (`[A-Za-z]{2}`) and automatic uppercase normalization for `CountryCode`.
- Added strict `http`/`https` scheme and host validation for `BaseURL` during client initialization.
- Upgraded release workflow example verification step to Go 1.22.

### Changed
- Automated documentation cross-repository version dispatch synchronization integration.

## [2.1.0] - 2026-08-23

### Added
- Dynamic API key rotation via `APIKeySupplier func() string` in `ClientConfig` for zero-downtime key rotation.
- Multi-tenant routing support with `x-api-user` header via `BulkEnrichmentOptions` in `EnrichTransactionsWithOptions` and `GetEnrichmentStatusWithOptions`.
- Configurable zero-trust archive download host allowlist via `TrustedDownloadHosts` in `ClientConfig`.
- Client `Close()` method implementing `io.Closer` to close idle HTTP transport connections.
- Release workflow with SLSA build provenance attestations, SPDX SBOM generation, SHA-256 checksums, and automated GitHub release publishing.

### Security
- Strict zero-trust domain allowlisting for `DownloadEnrichmentCollection` to prevent Server-Side Request Forgery (SSRF).
- Eliminated response body preview in Content-Type error messages to prevent sensitive HTML/SSRF data leaks.
- Active stream overflow erroring with `maxBytesReader` on compressed and decompressed archive streams.
- Zip-Slip and path traversal mitigations during archive tarball extraction.
- CRLF injection validation for `x-api-user` header values.
- Tar entry filename sanitization for safe error reporting without control character injection.
- Enforced TLS 1.2+ minimum protocol version with modern cipher suites in `DefaultEnterpriseTransport`.

### Fixed
- Switched input length validation in `EnrichmentRequest.Validate()` to `utf8.RuneCountInString` (max 128 Unicode runes), supporting multi-byte international characters and currency symbols.

## [2.0.2] - 2026-08-15

### Added
- In-memory streaming tarball decompression for bulk collection downloads (`DownloadEnrichmentCollection`).
- Decompression bomb mitigations: aggregate stream byte bounds (`DefaultMaxArchiveBytes = 100MiB`), per-entry size limits (`DefaultMaxEntryBytes = 10MiB`), and entry count bounds (`DefaultMaxTarEntries = 50,000`).
- Strict URL scheme validation (`http`/`https`) for bulk download endpoints.
- Authorization header cross-host leak prevention on downloads.
- Automatic secret token redaction in OpenAPI client debug logging (`Authorization`, `Proxy-Authorization`, `Cookie`).
- Bounded response body reading across all OpenAPI endpoints (`maxResponseBytes = 32MiB`).
- Enterprise connection-pooled transport (`DefaultEnterpriseTransport` / `NewDefaultHTTPClient`) with `MaxIdleConnsPerHost = 100` to prevent TCP `TIME_WAIT` socket exhaustion under heavy concurrent load.
- Dedicated SDK version constants in `version.go` (`Version`, `DefaultUserAgent`).
- Client-side input validation guards for enrichment requests (`EnrichmentRequest.Validate()`).

### Changed
- Removed `decoder.DisallowUnknownFields()` across models to guarantee forward-compatible JSON decoding for API additive schema updates.
- Upgraded CI workflows and Docker environments to Go 1.22 with SLSA attestations and SBOM generation.
- Expanded static analysis suite with `gosec` and `bodyclose` linters.
- Updated mascot and documentation branding from Go Financial to XYO Financial.

## [2.0.1] - 2026-08-14

### Changed
- Aligned `LICENSE` with exact standard Apache 2.0 text for `licensecheck` compliance.

## [2.0.0] - 2026-08-12

### Added
- Migrated module path to `github.com/xyo-financial/sdk-go/v2`.
- Full OpenAPI-driven client architecture with strongly-typed interfaces.
- Structured RFC 7807 problem details error handling via `ErrorResponse` and `APIError`, supporting `errors.As` unwrapping.
- Modern idiomatic Go client with `EnrichTransaction`, `EnrichTransactions`, `GetEnrichmentStatus`, and `DownloadEnrichmentCollection` methods.
- Comprehensive test suite covering success paths, error handling, and edge cases.

## [1.2.4] - 2026-08-07

### Changed
- Updated repository URL and module path to `github.com/xyo-financial/sdk-go`.

## [1.2.3] - 2026-07-20

### Changed
- Updated SDK license.

## [1.2.2] - 2026-07-17

### Fixed
- Removed AI session ID from documentation.

## [1.2.1] - 2026-07-17

### Added
- Added pkg.go.dev badge link.

## [1.2.0] - 2026-07-17

### Added
- Maintenance release for July 2026.

## [1.1.0] - 2025-10-19

### Added
- Documentation updates and two new enrichment fields.

## [1.0.4] - 2025-09-03

### Changed
- Made `ApiBasePath` private and added base path configuration to client config.

## [1.0.3] - 2025-09-03

### Changed
- Updated tests and created internal utility method to set request headers.

## [1.0.2] - 2025-09-03

### Changed
- Documentation updates.

## [1.0.1] - 2025-09-02

### Changed
- Documentation updates.

## [1.0.0] - 2025-08-22

### Added
- Initial release of the Go SDK for XYO Financial API.
- `NewClient` constructor accepting `ClientConfig` with API key.
- `EnrichTransaction` method for single payment transaction enrichment.
- `EnrichTransactionCollection` method for bulk enrichment requests.
- `EnrichTransactionCollectionStatus` method to query bulk enrichment status.
- Pluggable HTTP transport layer via `httpClient` for testability.
- Unit tests covering all enrichment endpoints (success and error paths).
- Docker-based CI pipeline with `golangci-lint` and `go test`.
- Example application in `example/` directory.

[Unreleased]: https://github.com/xyo-financial/sdk-go/compare/v2.1.0...HEAD
[2.1.0]: https://github.com/xyo-financial/sdk-go/compare/v2.0.2...v2.1.0
[2.0.2]: https://github.com/xyo-financial/sdk-go/compare/v2.0.1...v2.0.2
[2.0.1]: https://github.com/xyo-financial/sdk-go/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/xyo-financial/sdk-go/compare/v1.2.4...v2.0.0
[1.2.4]: https://github.com/xyo-financial/sdk-go/compare/v1.2.3...v1.2.4
[1.2.3]: https://github.com/xyo-financial/sdk-go/compare/v1.2.2...v1.2.3
[1.2.2]: https://github.com/xyo-financial/sdk-go/compare/v1.2.1...v1.2.2
[1.2.1]: https://github.com/xyo-financial/sdk-go/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/xyo-financial/sdk-go/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/xyo-financial/sdk-go/compare/v1.0.4...v1.1.0
[1.0.4]: https://github.com/xyo-financial/sdk-go/compare/v1.0.3...v1.0.4
[1.0.3]: https://github.com/xyo-financial/sdk-go/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/xyo-financial/sdk-go/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/xyo-financial/sdk-go/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/xyo-financial/sdk-go/releases/tag/v1.0.0
