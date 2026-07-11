# Contributing to XYO.Financial Go SDK

First, thank you for taking the time to contribute.

## Development Setup

1. Make sure you have Go 1.18 or newer installed.
2. Clone the repository:
   ```sh
   git clone https://github.com/syniol/xyo-sdk-go.git
   cd xyo-sdk-go
   ```
3. Run the tests to verify everything works:
   ```sh
   go test ./... -v
   ```

## Coding Standards

- Follow the standard Go conventions (`gofmt`, `go vet`).
- Use `golangci-lint` for linting — the CI pipeline enforces this.
- Keep the public API surface minimal and well-documented with GoDoc comments.
- Use the `internal/` package for implementation details that should not be exposed.
- Write table-driven tests using the standard `testing` package.

## Submitting a Pull Request

1. Fork the repository and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes: `go test ./... -v`
5. Ensure the linter passes: `golangci-lint run ./...`
6. Open your pull request.

## Publishing a New Version

1. Push all the new changes to the `main` branch.
2. `git tag v1.x.x`
3. `git push origin v1.x.x`
4. The `release.yml` pipeline will automatically run tests, generate artifacts (SBOM, checksums, attestation), and create a GitHub Release.

Thank you
