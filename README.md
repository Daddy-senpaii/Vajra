# Vajra

Vajra is a cross-platform adult-content blocking application focused on helping users maintain commitments through reliable content blocking.

## Current Status

- **Phase 1: Core Blocking Engine** (implemented; sample data only)

See [AGENT.md](AGENT.md) for architecture, guidelines, and roadmap.

## Core engine

The standalone Go engine lives in `packages/blocking-engine`. It normalizes and
validates domain inputs, supports IDNs, matches blocklist domains and their
subdomains, and gives allowlist entries precedence. It works offline and reads
local text lists; no production adult-domain feed is included. The checked-in
lists under `packages/blocking-engine/testdata` are fixtures for development and
testing only.

Run the engine tests from its module directory:

```sh
cd packages/blocking-engine
go test ./...
```

Evaluate one domain using the sample blocklist:

```sh
go run ./cmd/blocklist-tool \
  -blocklist testdata/sample_blocklist.txt \
  -allowlist testdata/sample_allowlist.txt \
  -domain www.adult-example.com
```

Omit `-domain` to read domains interactively from standard input. Blocklist and
allowlist files use one domain per line; blank lines and `#` comments are
ignored. Run lookup benchmarks with `go test ./benchmarks -bench .`.
