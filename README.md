<p align="center">
  <a href="https://xyo.financial" target="_blank" rel="noopener noreferrer">
    <img alt="XYO Financial Go Mascot" width="380" src="https://raw.githubusercontent.com/xyo-financial/sdk-go/main/docs/mascot.png" />
  </a>
</p>

<h1 align="center">XYO Financial SDK for Go</h1>

<p align="center">
  <a href="https://github.com/xyo-financial/sdk-go/actions/workflows/makefile.yml"><img src="https://github.com/xyo-financial/sdk-go/actions/workflows/makefile.yml/badge.svg?branch=main" alt="CI Build Pipeline" /></a>
  <a href="https://github.com/xyo-financial/sdk-go/actions/workflows/release.yml"><img src="https://github.com/xyo-financial/sdk-go/actions/workflows/release.yml/badge.svg" alt="Release Pipeline" /></a>
  <a href="https://pkg.go.dev/github.com/xyo-financial/sdk-go/v2"><img src="https://pkg.go.dev/badge/github.com/xyo-financial/sdk-go/v2" alt="Go Reference" /></a>
  <a href="https://go.dev/doc/go1.22"><img src="https://img.shields.io/badge/Go-%3E%3D1.22-00ADD8?logo=go&logoColor=white" alt="Go Compatibility" /></a>
  <a href="https://datatracker.ietf.org/doc/html/rfc7807"><img src="https://img.shields.io/badge/RFC_7807-Compliant-brightgreen" alt="RFC 7807 Compliant" /></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License" /></a>
</p>

<p align="center">
  <strong>The official Go SDK for <a href="https://xyo.financial">XYO Financial</a>.</strong><br>
  Seamlessly enrich raw financial transactions into clean merchant profiles, intelligent business categorizations, high-res logos, and geolocated address metadata using AI-powered enrichment pipelines.
</p>

---

## 📖 Summary

The **XYO Financial SDK for Go** provides a high-performance, strongly typed client library for integrating XYO's AI-driven transaction enrichment engine into cloud-native Go microservices and enterprise financial backends. 

Engineered for Tier-1 banks, payment processors, and fintech platforms, this SDK transforms raw, cryptic merchant descriptions (e.g., `AMZN MKTP UK*1M23456`, `SQ *COSTA GREENWICH`) into structured, verified merchant profiles complete with official names, descriptions, merchant logos, standard industry categories, and geocoded location data.

Maintained by [Syniol Limited](https://syniol.com) as the official Go distribution for XYO.Financial.

---

## 🏗 Architectural Principles

1. **Context-Native & Cancelable**: Every client method accepts `context.Context` as its first parameter, supporting fine-grained request timeouts, deadline propagation, and instant cancellation across distributed call trees.
2. **Thread-Safe & Reusable**: The `xyo.Client` is fully stateless and safe for concurrent use across goroutines. Recommended for allocation as a long-lived application singleton.
3. **Structured RFC 7807 Error Handling**: API exceptions preserve RFC 7807 problem details (`Type`, `Title`, `Status`, `Detail`, `Instance`) accessible through idiomatic Go 1.13+ error unwrapping (`errors.As`).
4. **Resilient Streaming**: Bulk collection results are streamed and decompressed on-the-fly (`tar.gz`) directly into memory without requiring temporary disk buffers.
5. **Zero External Runtime Footprint**: Relies solely on standard Go libraries and deterministic OpenAPI client generation for minimal supply-chain risk.
6. **Decompression-Bomb Hardened**: The bulk download pipeline enforces three independent safety limits — maximum entry count, maximum entry size, and maximum total archive size — using `io.LimitReader` on both the compressed and decompressed streams, preventing adversarial archives from exhausting host memory.
7. **Credential-Scoped Downloads**: The `Authorization` header is only attached to requests whose destination host matches the configured API base URL. Downloads redirected to pre-signed S3, GCS, or CDN endpoints never receive the Bearer token.

---

## ⚙️ System Requirements

- **Go**: Version `1.22` or newer.
- **Network**: Outbound HTTPS connectivity to `api.xyo.financial` over port `443` (TLS 1.2+ mandatory).
- **Authentication**: A valid API Key obtained from the [XYO Dashboard](https://xyo.financial/dashboard).

---

## 📦 Installation

Add the SDK to your Go module:

```bash
go get github.com/xyo-financial/sdk-go/v2
```

---

## 🚀 Quickstart Guide

### 1. Client Initialization

Initialize the `xyo.Client` using `xyo.NewClient`. Pass a configuration struct specifying your API key:

```go
package main

import (
	"log"
	"os"

	"github.com/xyo-financial/sdk-go/v2"
)

func main() {
	client, err := xyo.NewClient(&xyo.Config{
		APIKey: os.Getenv("XYO_API_KEY"),
	})
	if err != nil {
		log.Fatalf("failed to initialize XYO client: %v", err)
	}

	_ = client
}
```

> **Note:** `xyo.Config` is an alias for `xyo.ClientConfig`. By default, `NewClient` initializes a production-ready `*http.Client` using `DefaultEnterpriseTransport` (`MaxIdleConnsPerHost = 100`, HTTP/2 enabled) and a 30-second timeout to prevent TCP socket exhaustion under high concurrency. You may also configure a custom `BaseURL` (e.g. for staging or mock testing) and a custom `*http.Client`.

---

### 2. Single Transaction Enrichment (`EnrichTransaction`)

Enrich a single payment transaction synchronously. Best suited for real-time authorization paths, checkout flows, and interactive user banking dashboards:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xyo-financial/sdk-go/v2"
)

func main() {
	client, err := xyo.NewClient(&xyo.Config{
		APIKey: os.Getenv("XYO_API_KEY"),
	})
	if err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.EnrichTransaction(ctx, &xyo.EnrichmentRequest{
		Content:     "COSTA PICKUP",
		CountryCode: "GB",
	})
	if err != nil {
		log.Fatalf("enrichment failed: %v", err)
	}

	fmt.Printf("Merchant:    %s\n", resp.Merchant)
	fmt.Printf("Description: %s\n", resp.Description)
	fmt.Printf("Categories:  %v\n", resp.Categories)
	fmt.Printf("Logo:        %s\n", resp.Logo)
	fmt.Printf("Location:    %s\n", resp.Location)
	fmt.Printf("Address:     %s\n", resp.Address)
}
```

---

### 3. Bulk Transaction Enrichment (`EnrichTransactions`)

Submit high-volume transaction batches asynchronously for ETL pipelines, end-of-day reconciliation, and historical statement enrichment:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xyo-financial/sdk-go/v2"
)

func main() {
	client, err := xyo.NewClient(&xyo.Config{
		APIKey: os.Getenv("XYO_API_KEY"),
	})
	if err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	requests := []*xyo.EnrichmentRequest{
		{
			Content:     "Costa PICKUP",
			CountryCode: "GB",
		},
		{
			Content:     "STRBUKS GREENWICH",
			CountryCode: "GB",
		},
		{
			Content:     "UBER TRIP HELP.UBER.COM",
			CountryCode: "US",
		},
	}

	batch, err := client.EnrichTransactions(ctx, requests)
	if err != nil {
		log.Fatalf("bulk enrichment submission failed: %v", err)
	}

	fmt.Printf("Batch Job ID:  %s\n", batch.ID)
	fmt.Printf("Download Link: %s\n", batch.Link)
}
```

---

### 4. Bulk Job Polling & Streaming Download (`GetEnrichmentStatus`)

Poll the status of an asynchronous bulk enrichment job and stream-decode the results directly into memory once ready:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xyo-financial/sdk-go/v2"
)

func main() {
	client, err := xyo.NewClient(&xyo.Config{
		APIKey: os.Getenv("XYO_API_KEY"),
	})
	if err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}

	jobID := "72c037df-d0d3-43ee-9470-323ff35a2e50"
	downloadURL := fmt.Sprintf("https://api.xyo.financial/ai/transactions/download/%s.tar.gz", jobID)

	// Poll with timeout and exponential backoff
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	backoff := 2 * time.Second
	maxBackoff := 15 * time.Second

	for {
		status, err := client.GetEnrichmentStatus(ctx, jobID)
		if err != nil {
			log.Fatalf("failed checking status: %v", err)
		}

		fmt.Printf("Job %s status: %s\n", jobID, status)

		if status == xyo.EnrichmentCollectionStatusReady {
			break
		}
		if status == xyo.EnrichmentCollectionStatusFailed {
			log.Fatalf("bulk enrichment job %s failed processing on server", jobID)
		}

		select {
		case <-ctx.Done():
			log.Fatalf("timed out waiting for bulk enrichment: %v", ctx.Err())
		case <-time.After(backoff):
			if backoff < maxBackoff {
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
			}
		}
	}

	// Stream and decompress results on-the-fly
	results, err := client.DownloadEnrichmentCollection(ctx, downloadURL)
	if err != nil {
		log.Fatalf("failed downloading collection: %v", err)
	}

	for i, item := range results {
		fmt.Printf("[%d] Merchant: %-20s | Categories: %v\n", i+1, item.Merchant, item.Categories)
	}
}
```

---

## 🚀 Framework & Architecture Integration

The XYO Go SDK is architected for zero-overhead integration into modern Go microservice frameworks, cloud-native API gateways, and serverless compute platforms.

### Microservice Handler Pattern (Gin)

Inject a long-lived `xyo.Client` singleton into your Gin HTTP handlers with deterministic context timeout propagation:

```go
package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xyo-financial/sdk-go/v2"
)

// EnrichHandler mounts the XYO transaction enrichment pipeline as a Gin handler.
func EnrichHandler(client xyo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req xyo.EnrichmentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Propagate HTTP request context with strict deadline
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		res, err := client.EnrichTransaction(ctx, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func main() {
	client, err := xyo.NewClient(&xyo.Config{
		APIKey: os.Getenv("XYO_API_KEY"),
	})
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.POST("/api/v1/enrich", EnrichHandler(client))
	_ = r.Run(":8080")
}
```

---

### Standard Library Microservice Pattern (`net/http`)

For lean, zero-dependency services, integrate directly with Go's standard `net/http` server (supporting Go 1.22+ enhanced routing patterns):

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/xyo-financial/sdk-go/v2"
)

// EnrichHTTPHandler returns an idiomatic standard library net/http handler.
func EnrichHTTPHandler(client xyo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req xyo.EnrichmentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("invalid request payload: %v", err), http.StatusBadRequest)
			return
		}

		// Propagate request context with 5s timeout budget
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		res, err := client.EnrichTransaction(ctx, &req)
		if err != nil {
			http.Error(w, fmt.Sprintf("enrichment error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(res)
	}
}

func main() {
	client, err := xyo.NewClient(&xyo.Config{
		APIKey: os.Getenv("XYO_API_KEY"),
	})
	if err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/enrich", EnrichHTTPHandler(client))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server listening on :8080")
	log.Fatal(server.ListenAndServe())
}
```

---

### Production Architecture Highlights

* **0 External Dependencies**: The SDK core utilizes standard Go libraries exclusively (`net/http`, `crypto/tls`, `encoding/json`), guaranteeing a completely flat dependency tree with zero transitive CVE exposure and minimal binary size.
* **Thread-Safe Singleton Reuse**: The `xyo.Client` is stateless and safe for concurrent execution across unlimited goroutines. Initialize a single client instance at application startup and share it across all request handlers to benefit from connection reuse via `DefaultEnterpriseTransport`.
* **~1ms Static Binary Cold-Start Execution**: Compiles into self-contained, statically linked binaries (`CGO_ENABLED=0`) with instantaneous **~1ms cold-start latency** and sub-10MB memory footprints—ideal for AWS Lambda, Google Cloud Run, Knative, and high-density Kubernetes pods.

---

## 🛡 RFC 7807 Structured Error Handling

The XYO API adheres to the **RFC 7807 (Problem Details for HTTP APIs)** specification. When an API call fails, the SDK wraps the structured error response in an `*xyo.ErrorResponse` containing one or more `*xyo.APIError` entries.

The SDK implements a **three-tier error parsing cascade** to ensure a structured error is always returned, regardless of how the server or upstream infrastructure (proxies, WAFs, CDN edge nodes) responds:

| Tier | Mechanism | Fires When |
|:---|:---|:---|
| **1st** | Typed `openapi.ErrorResponse` model decode | Server returns well-formed `application/json` error body |
| **2nd** | Raw `[]byte` JSON unmarshal fallback | Server returns valid JSON but with incorrect `Content-Type` |
| **3rd** | HTTP status code extraction | Server returns HTML, plain text, or an empty body (e.g. `502 Bad Gateway` from a proxy) |

Use `errors.As` for type-safe inspection of error codes and detailed diagnostic reasons:

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/xyo-financial/sdk-go/v2"
)

func enrichWithDiagnostics(client xyo.Client, rawContent, countryCode string) (*xyo.EnrichmentResponse, error) {
	resp, err := client.EnrichTransaction(context.Background(), &xyo.EnrichmentRequest{
		Content:     rawContent,
		CountryCode: countryCode,
	})
	if err != nil {
		var apiErr *xyo.ErrorResponse
		if errors.As(err, &apiErr) {
			fmt.Printf("XYO API Error (HTTP Status %d):\n", apiErr.HTTPStatusCode)
			for _, e := range apiErr.Errors {
				fmt.Printf("  - Type:     %s\n", e.Type)
				fmt.Printf("    Title:    %s\n", e.Title)
				fmt.Printf("    Detail:   %s\n", e.Detail)
				fmt.Printf("    Instance: %s\n", e.Instance)
			}

			switch apiErr.HTTPStatusCode {
			case http.StatusBadRequest:
				// Invalid request payload — verify content format and country code (ISO 3166-1 alpha-2)
			case http.StatusUnauthorized, http.StatusForbidden:
				// Authentication failure — verify API key validity and billing status
			case http.StatusTooManyRequests:
				// Rate limit exceeded — back off and retry
			case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
				// Server degradation — route to fallback pipeline
			}
			return nil, err
		}

		// Transport, connection, or context cancellation error
		return nil, fmt.Errorf("transport failure: %w", err)
	}

	return resp, nil
}
```

### HTTP Status Code Reference

| HTTP Status | Category | Description & Enterprise Mitigation |
|:---|:---|:---|
| `400` | Bad Request | Invalid parameter (e.g. content > 128 characters or invalid ISO country code). Route payload to Dead Letter Queue (DLQ). |
| `401` / `403` | Unauthorized / Forbidden | Invalid or revoked API key, or account subscription expired. Alert SecOps/DevOps immediately. |
| `422` | Unprocessable Entity | Content could not be parsed into valid enrichment data. |
| `429` | Rate Limited | Quota or throughput exceeded. Implement exponential backoff with jitter. |
| `500` / `502` / `503` | Server Error | Upstream service degradation. Open circuit breaker; fallback to raw payment record. |

---

## ⚙️ Advanced Client Configuration

```go
customHTTPClient := &http.Client{
	Timeout: 45 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	},
}

client, err := xyo.NewClient(&xyo.Config{
	APIKey:     os.Getenv("XYO_API_KEY"),
	BaseURL:    "https://api.xyo.financial", // Optional override
	HTTPClient: customHTTPClient,            // Custom transport configuration
})
```

---

## 🔄 Backward Compatibility

For backward compatibility with earlier integrations, the following alias methods are provided:

| Canonical Method | Backward-Compatible Alias | Purpose |
|:---|:---|:---|
| `EnrichTransactions(...)` | `EnrichTransactionCollection(...)` | Submit bulk transaction enrichment |
| `GetEnrichmentStatus(...)` | `EnrichTransactionCollectionStatus(...)` | Check bulk enrichment processing status |
| `xyo.Config` | `xyo.ClientConfig` | Configuration struct alias |

---

## 🔒 Security & Compliance

- **Zero Third-Party Runtime Dependencies**: The entire SDK is built on the Go standard library. The `go.mod` contains no external packages, producing a completely flat dependency graph — no transitive CVE exposure, no supply-chain audit overhead.
- **Data Minimisation**: Only raw payment descriptions and ISO 3166-1 country codes are transmitted. No card numbers (PAN), CVVs, sort codes, or personally identifiable information (PII) required.
- **TLS 1.2+ Enforcement**: All requests are routed exclusively over TLS-encrypted channels. Bearer token authorization headers are transmitted only over HTTPS.
- **Credential-Scoped Downloads**: The `Authorization` header is conditionally attached only when the download URL host matches the configured API base URL. Pre-signed S3, GCS, or CDN download links — which originate from the batch enrichment workflow — are served without the Bearer token, preventing API key leakage to third-party storage providers.
- **Release Provenance & SBOM**: Every release is published with a cryptographically signed **GitHub Artifact Attestation** (build provenance), an **SPDX Software Bill of Materials**, and a **SHA-256 checksum** — providing SLSA Level 2 supply-chain integrity out of the box.

---

## 📞 Support

For enterprise SLAs, custom rate tiers, or dedicated technical support:
- **Dashboard**: [https://xyo.financial/dashboard](https://xyo.financial/dashboard)
- **Email Support**: [support@syniol.com](mailto:support@syniol.com)
- **Maintained By**: [Syniol Limited](https://syniol.com)

---

## 📄 License

This project is licensed under the **Apache License, Version 2.0** - see the [LICENSE](LICENSE) file for details.

Copyright &copy; 2026 Syniol Limited. All rights reserved.
