# Vajra

Vajra is a cross-platform adult-content blocking application focused on helping users maintain commitments through reliable content blocking.

## Current Status

- **Phase 1: Core Blocking Engine** (offline engine with bundled adult-domain snapshot)

## Core engine

The standalone Go engine lives in `packages/blocking-engine`. It normalizes and
validates domain inputs, supports IDNs, matches blocklist domains and their
subdomains, and gives allowlist entries precedence. It works offline and ships
with a snapshot of a maintained adult-domain list. See
[`blocklist/data/README.md`](packages/blocking-engine/blocklist/data/README.md)
for its source, license, snapshot date, and limitations. The small lists under
`packages/blocking-engine/testdata` remain test fixtures.

The Chrome extension packages 30,000 static blocking rules, using Chrome's
guaranteed static-rule minimum. Major adult domains are prioritized; the
remaining browser rules are a deterministic sample of the full list. The Go
engine itself loads the full bundled list.

Run the engine tests from its module directory:

```sh
cd packages/blocking-engine
go test ./...
```

Evaluate one domain using the bundled adult-domain list:

```sh
go run ./cmd/blocklist-tool \
  -domain www.pornhub.com
```

To use custom data, pass `-blocklist path/to/domains.txt` and/or `-allowlist
path/to/domains.txt`. List files use one domain per line; blank lines and `#`
comments are ignored. Omit `-domain` to read domains interactively from standard
input. Run lookup benchmarks with `go test ./benchmarks -bench .`.

Refresh the bundled adult-domain snapshot and regenerate browser rules with
`scripts/update-adult-blocklist.sh`, then review the updated data, metadata, and
rules.

## Browser extension

Build the extension from `apps/extension` with `npm ci && npm run build`, then
load `apps/extension/dist` as an unpacked extension in Chrome. After rebuilding,
reload the extension from `chrome://extensions`. The packaged browser rules
cover 30,000 domains; the standalone Go engine uses the full 426,211-domain
snapshot.
