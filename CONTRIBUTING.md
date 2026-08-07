# Contributing to XYO.Financial Go SDK

Thank you for your interest in the XYO.Financial Go SDK. 

**Please note that we do not accept external code contributions or Pull Requests to this repository.** 

Development, code changes, and publishing are exclusively maintained by employees of **Syniol Limited**. 

## Raising Issues

While we do not accept external Pull Requests, we heavily rely on our community and users to report bugs, request features, and provide feedback. 

If you encounter an issue or have a feature request, please open an issue on our GitHub repository. Provide as much detail as possible, including:
- SDK version
- Go environment details
- A minimal reproducible example

## Internal Development (Syniol Limited Employees Only)

The following sections are for internal engineering use.

### Development Setup

1. Make sure you have Go 1.18 or newer installed.
2. Clone the repository:
   ```sh
   git clone https://github.com/xyo-financial/sdk-go.git
   cd sdk-go
   ```
3. Run the tests to verify everything works:
   ```sh
   go test ./... -v
   ```

### Coding Standards

- Follow the standard Go conventions (`gofmt`, `go vet`).
- Use `golangci-lint` for linting — the CI pipeline enforces this.
- Keep the public API surface minimal and well-documented with GoDoc comments.
- Use the `internal/` package for implementation details that should not be exposed.
- Write table-driven tests using the standard `testing` package.

### Submitting a Pull Request

1. Branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes: `go test ./... -v`
5. Ensure the linter passes: `golangci-lint run ./...`
6. Open your pull request for internal review.

### Publishing a New Version

1. Push all the new changes to the `main` branch.
2. `git tag v1.x.x`
3. `git push origin v1.x.x`
4. The `release.yml` pipeline will automatically run tests, generate artifacts (SBOM, checksums, attestation), and create a GitHub Release.
