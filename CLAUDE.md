# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is Hoverfly

Hoverfly is an HTTP/HTTPS proxy for simulating APIs and services. It can capture, simulate, modify, synthesize, spy on, or diff HTTP traffic. It ships two binaries: `hoverfly` (the proxy server) and `hoverctl` (the CLI to manage it).

## Commands

```bash
# Build both binaries to target/
make build

# Run all tests (unit + functional + vet)
make test

# Hoverfly unit tests only
make hoverfly-test

# Hoverctl unit tests only
make hoverctl-test

# Functional tests (requires built binaries)
make hoverfly-functional-test
make hoverctl-functional-test

# Run a single test (from the relevant package directory)
cd core && go test -v ./... -run TestFunctionName

# Format and vet
make fmt
make vet
```

Tests use the **Ginkgo/Gomega** BDD framework. Functional tests in `functional-tests/` spin up actual hoverfly instances.

## Architecture

### Entry Points

- **`core/cmd/hoverfly/main.go`** — parses 50+ flags, wires up Hoverfly, starts proxy + admin API
- **`hoverctl/main.go`** — Cobra-based CLI; communicates with a running hoverfly instance via HTTP API

### Core Request Flow

1. **goproxy** intercepts incoming HTTP(S) requests (`core/proxy.go`)
2. The active **Mode** processes the request (`core/modes/`)
3. **Matching** finds a recorded response against the simulation store (`core/matching/`)
4. **Templating** renders dynamic response bodies if configured (`core/templating/`)
5. **Middleware** optionally transforms request/response (local binary or remote HTTP) (`core/middleware/`)
6. **Delay** is applied (`core/delay/`)
7. **Post-serve actions** fire (webhooks, scripts) (`core/action/`)
8. **Journal** records the exchange (`core/journal/`)

### Mode System

Six modes implement the `Mode` interface with a `Process()` method (`core/modes/`):

| Mode | Behavior |
|------|----------|
| **Simulate** | Return recorded responses from simulation store |
| **Capture** | Forward requests and record request/response pairs |
| **Spy** | Simulate if matched, else pass through |
| **Modify** | Pass through but run middleware on traffic |
| **Synthesize** | Generate responses entirely via middleware |
| **Diff** | Forward requests and diff responses against simulation |

Mode can be switched dynamically at runtime via the admin API.

### Key Packages

- **`core/hoverfly.go`** — `Hoverfly` struct, the central coordinator; holds references to all subsystems
- **`core/models/`** — `Simulation`, `RequestResponsePair`, `RequestMatcher`, `ResponseDetails` — the core data model
- **`core/matching/`** — two strategies: `StrongestMatchStrategy` (weighted field scoring) and `FirstMatchStrategy`; wrapped by `CacheMatcher` (LRU)
- **`core/handlers/v2/`** — 54 REST handler files for the admin API (v1 is legacy)
- **`core/templating/`** — Handlebars templating via `raymond`; supports CSV/SQL data sources and journal variable injection
- **`core/middleware/`** — executes local binaries or calls remote HTTP endpoints; passes JSON request/response payloads
- **`core/authentication/`** — opt-in JWT/basic auth; disabled by default
- **`hoverctl/wrapper/`** — HTTP client wrapper that hoverctl uses to call the admin API

### Simulation Format

The canonical data format is the `Simulation` JSON structure (`core/models/simulation.go`). It contains an array of `RequestResponsePair` entries, each with a `RequestMatcher` (supporting exact, regex, glob, xpath, jsonpath matchers per field) and a `ResponseDetails`. Test fixtures are in `functional-tests/testdata/`.

### Matching

Matchers operate per HTTP field (method, path, query, headers, body, scheme, destination). Multiple matchers on the same field are ANDed. `StrongestMatchStrategy` picks the most specific match; `FirstMatchStrategy` picks the first. The `CacheMatcher` wraps either strategy with an LRU cache keyed on the request fingerprint.

## Tech Stack

- **Go 1.26.1**, modules in `go.mod`
- **Proxy:** `github.com/SpectoLabs/goproxy` (custom MITM fork)
- **CLI:** `cobra` + `viper`
- **Routing (admin API):** `gorilla/mux`, `go-zoo/bone`
- **Testing:** `ginkgo` + `gomega`
- **Logging:** `logrus`
- **Templating:** `github.com/SpectoLabs/raymond` (Handlebars)
- **Fake data:** `gofakeit/v6`
