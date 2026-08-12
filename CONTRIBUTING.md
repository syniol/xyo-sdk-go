# Contributing to XYO.Financial Go SDK

Thank you for your interest in the XYO.Financial Go SDK (`github.com/xyo-financial/sdk-go/v2`).

This document outlines the architecture, contribution workflow, code generation instructions, quality gates, and release standards required for maintaining and extending the SDK.

---

## 1. Contribution Policy & Issue Reporting

### Contribution Policy
Development, pull request reviews, and package publishing are maintained by **Syniol Limited** engineers. 

External pull requests submitted directly to this repository without prior coordination may be rejected. However, feedback, bug reports, and enhancement proposals from the developer community are welcome.

### Reporting Issues & Requesting Features
If you discover a bug, have questions, or wish to propose an enhancement:
1. Search existing [GitHub Issues](https://github.com/xyo-financial/sdk-go/issues) to see if the topic is already tracked.
2. Open a new issue with:
   - Go environment details (`go version`, `GOOS`, `GOARCH`).
   - SDK version or commit hash (`v2.x.x`).
   - Clear description of expected vs. actual behavior.
   - Minimal, reproducible Go code snippet reproducing the problem.

---

## 2. Two-Layer Architecture

The XYO Go SDK employs a **Two-Layer Architecture** designed to separate raw HTTP protocol serialization from developer ergonomics, reliability, and idiomatic Go idioms.

```
┌────────────────────────────────────────────────────────────────────────┐
│                        XYO Go SDK Architecture                         │
└────────────────────────────────────────────────────────────────────────┘

  Application Code
        │
        ▼
  ┌────────────────────────────────────────────────────────────────────┐
  │  WRAPPER LAYER (Hand-Crafted, Idiomatic Go)                        │
  │  Files: client.go, enrichment.go, errors.go                        │
  │  ────────────────────────────────────────────────────────────────  │
  │  • First-class context.Context cancellation & timeouts             │
  │  • High-level Client & Enrichment interfaces                       │
  │  • Idiomatic structs (EnrichmentRequest, EnrichmentResponse)       │
  │  • Structured error unmarshaling & RFC 7807 problem details        │
  │  • Safe default timeout (30s) & custom http.Client injection       │
  │  • Tarball download & decompression helpers (bulk enrichment)      │
  └────────────────────────────────────────────────────────────────────┘
        │
        ▼
  ┌────────────────────────────────────────────────────────────────────┐
  │  GENERATED LAYER (OpenAPI Generator Artifacts)                     │
  │  Directory: openapi/                                               │
  │  ────────────────────────────────────────────────────────────────  │
  │  • Auto-generated from xyo-financial/specs (openapi.yml)           │
  │  • Low-level HTTP requests, payload serialization, URL routing     │
  │  • Pure transport layer: openapi.APIClient, openapi.EnrichmentAPI  │
  │  • ⚠️ READ-ONLY: DO NOT EDIT OR FORMAT MANUALLY                    │
  └────────────────────────────────────────────────────────────────────┘
        │
        ▼
  XYO RESTful API (https://api.xyo.financial)
```

### Layer 1: Low-Level Generated Layer (`openapi/`)
- **Source:** Automatically generated from the canonical OpenAPI specification maintained in [`xyo-financial/specs`](https://github.com/xyo-financial/specs).
- **Purpose:** Handles low-level HTTP wire protocol serialization, JSON marshaling/unmarshaling, API endpoint routing, and raw response parsing (`openapi.APIClient`, `openapi.EnrichmentAPI`, `openapi.EnrichmentRequest`, `openapi.GenericOpenAPIError`).
- **Policy:** **READ-ONLY & IMMUTABLE**. Under no circumstances should files inside the `openapi/` directory be edited or reformatted manually. Any manual changes will be overwritten during the next code generation cycle.

### Layer 2: Ergonomic Wrapper Layer (`client.go`, `enrichment.go`, `errors.go`)
- **Source:** Hand-crafted Go code written and maintained directly in this repository.
- **Purpose:** Provides a clean, ergonomic, and idiomatic Go interface for developers:
  - [`client.go`](file:///Users/hadi/dev/start-ups/xyo/sdks/golang/client.go): Configures and instantiates the [`Client`](file:///Users/hadi/dev/start-ups/xyo/sdks/golang/client.go) interface with [`ClientConfig`](file:///Users/hadi/dev/start-ups/xyo/sdks/golang/client.go), default production timeouts (30s), Bearer token injection, and customizable HTTP transports.
  - [`enrichment.go`](file:///Users/hadi/dev/start-ups/xyo/sdks/golang/enrichment.go): Exposes the [`Enrichment`](file:///Users/hadi/dev/start-ups/xyo/sdks/golang/enrichment.go) interface with methods like `EnrichTransaction`, `EnrichTransactions`, `GetEnrichmentStatus`, and `DownloadEnrichmentCollection` (handling `tar.gz` streaming decompressions), with full `context.Context` cancellation support.
  - [`errors.go`](file:///Users/hadi/dev/start-ups/xyo/sdks/golang/errors.go): Translates low-level OpenAPI and HTTP errors into structured [`ErrorResponse`](file:///Users/hadi/dev/start-ups/xyo/sdks/golang/errors.go) and [`APIError`](file:///Users/hadi/dev/start-ups/xyo/sdks/golang/errors.go) objects conforming to RFC 7807 problem details, preserving HTTP status codes and actionable error messages.

---

## 3. Contribution Workflow & Decision Matrix

To maintain integrity across the XYO SDK ecosystem, determine where your proposed change belongs:

| Type of Change | Target Repository | Action Required |
| :--- | :--- | :--- |
| **API Endpoints & Contracts**<br>(New endpoints, URL changes, HTTP methods) | [`xyo-financial/specs`](https://github.com/xyo-financial/specs) | Submit PR to update `openapi.yml`. Once merged and tagged, the Go SDK regenerates automatically. |
| **Data Models & Schema Properties**<br>(New fields, type changes, validation rules) | [`xyo-financial/specs`](https://github.com/xyo-financial/specs) | Update the OpenAPI schema definitions in `specs`. |
| **SDK Ergonomics & Helpers**<br>(Higher-level methods, batch helpers, convenience wrappers) | [`xyo-financial/sdk-go`](https://github.com/xyo-financial/sdk-go) (This repo) | Implement directly in `client.go`, `enrichment.go`, or new wrapper files. |
| **Context & Transport Handling**<br>(Timeouts, retry middleware, HTTP client settings) | [`xyo-financial/sdk-go`](https://github.com/xyo-financial/sdk-go) (This repo) | Enhance wrapper layer and configuration options. |
| **Error Handling & Diagnostics**<br>(RFC 7807 mapping, status code extraction, unwrapping) | [`xyo-financial/sdk-go`](https://github.com/xyo-financial/sdk-go) (This repo) | Update `errors.go` and add corresponding test coverage. |
| **Unit & Integration Tests**<br>(Mock server tests, table-driven tests, benchmarks) | [`xyo-financial/sdk-go`](https://github.com/xyo-financial/sdk-go) (This repo) | Add tests in `*_test.go` files using standard Go `testing` packages. |

---

## 4. Code Generation

### Automated Upstream Synchronization
When a new release tag is pushed to [`xyo-financial/specs`](https://github.com/xyo-financial/specs), a GitHub Actions workflow automatically sends a dispatch event to this repository. The [`.github/workflows/generate.yml`](.github/workflows/generate.yml) workflow:
1. Checks out `xyo-financial/specs` at the tagged release.
2. Runs `openapi-generator-cli` to regenerate `openapi/`.
3. Cleans up generator scaffolding artifacts.
4. Compiles, vets, and commits the updated client.

### Manual / Local Code Generation
If you need to regenerate the low-level `openapi/` layer locally:

#### Prerequisites
- Node.js (v18+) with `npx`
- Go (v1.18+)
- Sibling clone of `xyo-financial/specs` (or direct path to `openapi.yml`)

#### Command
Run from the root of the Go SDK repository:

```bash
npx @openapitools/openapi-generator-cli generate \
  -i ../specs/openapi.yml \
  -g go \
  -o ./openapi \
  --additional-properties=packageName=openapi,withGoMod=false,hideGenerationTimestamp=true
```

#### Parameter Breakdown
- `-i ../specs/openapi.yml`: Path to the canonical OpenAPI specification.
- `-g go`: Targets the official OpenAPI Go client generator.
- `-o ./openapi`: Specifies the destination directory for generated artifacts.
- `--additional-properties=...`:
  - `packageName=openapi`: Declares `package openapi` across all generated files.
  - `withGoMod=false`: Suppresses generating an isolated `go.mod` file inside `openapi/`, ensuring it seamlessly forms a package of the parent `github.com/xyo-financial/sdk-go/v2`.
  - `hideGenerationTimestamp=true`: Prevents timestamp headers in generated comments, eliminating Git diff churn.

#### Post-Generation Clean-Up
After code generation completes, remove unnecessary generator metadata files (generated code in `openapi/` should remain untouched):

```bash
rm -f openapi/git_push.sh openapi/.travis.yml openapi/README.md
rm -rf openapi/test openapi/docs openapi/api
```

---

## 5. Quality Gates & Local Verification

All contributions must satisfy the following quality gates before being merged:

### 1. Unified Local Quality Check
Run the comprehensive pre-commit quality target:
```bash
make check
```
This runs format checks on non-generated code, `go vet`, and the complete test suite.

### 2. Code Formatting (Excluding Generated Code)
Format hand-crafted SDK code according to standard Go conventions:
```bash
make fmt
```
*(Note: Generated files in `openapi/` are ignored by `.golangci.yml` and the formatting rules).*

### 3. Build & Static Analysis
Ensure the entire module compiles cleanly and passes vetting:
```bash
go build ./...
go vet ./...
```

### 4. Test Suite & Race Detector
Execute all unit and integration tests with verbose output:
```bash
go test -v -race ./...
```
Ensure all test cases pass without panics, data races, or unhandled errors.

### 5. Example Application Verification
Confirm that the sample application compiles and runs without issues:
```bash
cd example && go run main.go && cd ..
```

### 6. Containerized Build (CI Parity)
Verify the production Docker build:
```bash
make build
```

---

## 6. Pull Request & Release Process

### Branching and Commits
1. Branch from `main` using descriptive branch names (e.g. `feat/enrichment-retry`, `fix/error-unmarshaling`).
2. Adhere to [Conventional Commits](https://www.conventionalcommits.org/) (e.g. `feat:`, `fix:`, `docs:`, `chore:`).
3. Ensure every new feature or bug fix includes corresponding table-driven unit tests in `*_test.go`.
4. Ensure all public types, interfaces, and functions include descriptive GoDoc comments.

### Release & Versioning Workflow
Releases follow [Semantic Versioning](https://semver.org/) (`vMAJOR.MINOR.PATCH`):

1. Ensure all changes are merged into `main` and all quality gates pass.
2. Create and push an annotated Git tag:
   ```bash
   git tag v2.0.0
   git push origin v2.0.0
   ```
3. The automated CI/CD release pipeline ([`.github/workflows/release.yml`](.github/workflows/release.yml)) will:
   - Verify tag integrity on the `main` branch.
   - Run full test suites (`go test ./... -v`).
   - Create release archive packages (`sdk-go-vX.Y.Z.tar.gz`).
   - Generate SHA-256 checksums and SPDX Software Bill of Materials (SBOM).
   - Sign and publish GitHub Artifact Attestations (build provenance).
   - Publish a GitHub Release with all cryptographic verification assets.
   - Execute verification tests on the example application.
