# Gemini CLI Brief: Vajra Backend with CleanBrowsing

You are taking over implementation of Vajra's adult-site blocking project. The user approved a Go backend that centralizes Vajra's policy feed and integrates with CleanBrowsing. CleanBrowsing's DNS resolver performs broad adult-category filtering; its REST API manages custom profiles and explicit domain entries. The backend must keep provider credentials private and serve the extension and future Android client.

The user wants a practical filter with broad coverage, similar in purpose to NordVPN's adult-site blocking. Do not promise “99%” coverage without a defined test corpus and measured results. Prioritize reliable blocking, honest protection status, offline operation, and low false positives.

## 1. Read and preserve the project

Before changing files:

1. Read the full root `AGENTS.md`, `README.md`, this brief, and all relevant engine and extension code.
2. Run `git status --short` and inspect diffs. Preserve all existing user changes and untracked files. Never reset, clean, checkout over, or delete user data.
3. Inspect the blocklist source, attribution, update script, curated data, extension rule generation, current tests, and the repository's actual app structure. The architecture document describes future structure; do not assume missing directories or modules exist.
4. Report a short architecture and implementation plan before code. The root `AGENTS.md` requires the user to approve the plan before implementation. Do not stop at the plan after approval; carry the work through implementation, verification, and documentation.

Known repository context at handoff:

- `packages/blocking-engine` is a Go engine with normalization, validation, allowlist, blocklist, and offline embedded data.
- `packages/blocking-engine/blocklist/data/curated-domains.txt` contains user-supplied labeled domains and source-listed variants. Load it in the engine and extension rules; retain its provenance.
- `apps/extension` is a Manifest V3 browser extension. Its generated rules currently use the supported static-rule budget and do not provide system-wide DNS filtering.
- `apps/mobile` may not exist yet. Verify first. The intended mobile stack in `AGENTS.md` is Expo/React Native with a Kotlin native Android module.
- `sites_labeled_dataset.json` and `site_domain_variants.json` are user-provided research inputs and may be untracked. Preserve them. Treat `generated` variants as guesses, never silently add them wholesale.
- Prior changes may be uncommitted. Do not overwrite them while implementing this feature.

## 2. Define the product behavior

Build the backend and client contracts with these behaviors:

- The Go engine remains the local policy source for Vajra's curated blocklist and extension rules.
- The API serves a versioned, ETag-enabled curated-domain feed and CleanBrowsing DNS configuration.
- The extension refreshes the small feed periodically, validates it, caches the last good copy, and keeps working offline with packaged rules.
- Android consumes `GET /v1/config` and uses the documented CleanBrowsing Adult Filter or the user's configured custom-profile DoT/DoH endpoints.
- A CleanBrowsing profile sync adds only Vajra-curated additions. Never try to download/upload CleanBrowsing's private category corpus or push the whole upstream snapshot.
- No per-navigation or per-DNS request is sent to Vajra's API. No query logs, browsing history, complete URLs, or user identifiers are collected by Vajra.
- Explicit allowlist precedence remains in Vajra's Go engine and browser extension. Explain that a provider's category-level block can have different allowlist behavior unless the CleanBrowsing custom profile is configured accordingly.

Treat this as domain-level filtering. DNS cannot tell whether a safe or adult page is being served from the same shared hostname. Do not claim to block page paths, image/video content, or every app's traffic with DNS alone.

## 3. Use CleanBrowsing without confusing the resolver and its REST API

The REST API is a management API, not a synchronous “is this domain adult?” API. It can manage profiles, explicit blocklist/allowlist entries, and network configuration. The DNS resolver applies CleanBrowsing's category filter. Keep that boundary clear in code and docs.

The implementation must address:

1. **Provider secret:** `CLEANBROWSING_API_KEY` and `VAJRA_ADMIN_TOKEN` are server-only environment values. Never place them in extension/mobile source, public configuration responses, logs, or committed files. CleanBrowsing requires its API key in a query parameter; do not log outbound request URLs or return URL-bearing transport errors.
2. **Provider profile:** sync only validated `curated-domains.txt` entries to a preconfigured CleanBrowsing profile. The Android DoT hostname/DoH URL must correspond to the profile receiving those entries. Public Adult Filter defaults do not automatically include account-specific custom domains.
3. **Safe sync:** list provider entries first, add only missing domains, make synchronization idempotent, preserve provider entries, report partial errors, and do not issue deletes. Do not sync all 426k upstream domains. Never invoke provider log endpoints.
4. **Public client API:** serve only safe resolver configuration and Vajra's curated feed. The admin sync endpoint must use a separate bearer secret and should be restricted at the network layer. Use HTTPS in deployment.
5. **Extension behavior:** fetch the feed on startup and periodically, use ETag/304, validate and deduplicate every domain, keep the last cached feed on network failure, and respect browser rule limits. Do not send each visited domain to the backend.
6. **Android behavior:** Android 9+ Private DNS can consume a DoT hostname. Use Android's explicit supported setup path; do not silently alter system security settings. If implementing a VPN-mode fallback, follow Android's official `VpnService` consent/lifecycle requirements and route only what the design can safely support.
7. **Coverage and privacy:** CleanBrowsing sees DNS queries sent to its resolver. Vajra does not store or proxy per-site queries. DNS cannot distinguish adult and safe pages on the same hostname or stop arbitrary DoH, direct-IP, or VPN bypasses.

If the CleanBrowsing API or profile configuration lacks a capability, do not invent one. Explain the limitation and use only supported profile management and resolver endpoints.

## 4. Improve blocklist coverage responsibly

- Keep the current upstream blocklist snapshot and its license/attribution intact. Do not edit the upstream snapshot while leaving its metadata checksum, line count, or source claims stale.
- Keep user-supplied curated domains in a separate file with clear provenance. Ensure the Go engine and generated extension rules both include them.
- Store only hostname/domain values. Reject URLs, paths, comments embedded in values, and malformed strings. Normalize IDNs consistently through the shared engine.
- Include candidates explicitly labeled as user-supplied/known or as source-listed variants only after validating their format. Preserve a route to remove false positives.
- Exclude pattern-generated guesses unless the user separately approves a reviewed set. Do not broaden a domain into its parent or unrelated siblings.
- Keep allowlist precedence and test false suffix matches such as `notexample.com` versus `example.com`.
- Make list updates atomic and validated (size, parse, checksum/signature where applicable) before replacing working data. Never execute downloaded list content. Preserve offline behavior.
- Record list version, update time, source, counts, and checksum for each independently maintained source. Do not state that the upstream snapshot includes curated additions.
- Do not add Supabase, user accounts, analytics, or cloud query reporting. The user specifically authorized the Go backend and CleanBrowsing integration only.

## 5. Implement in reviewable stages

After the user approves the plan, proceed in these stages and report progress:

### Stage A — Go backend

- Add `services/api` using the existing blocking engine and curated-domain data.
- Implement health, DNS configuration, versioned blocklist, and protected provider-sync endpoints.
- Add request timeouts, bounded responses, graceful shutdown, CORS configuration, and secret-safe errors.
- Document local setup, custom-profile setup, production HTTPS, and environment variables.

### Stage B — Client integration

- Connect the extension to the public curated feed using a build-time API base URL.
- Cache the latest good response and merge remote entries with user block entries while preserving allowlist precedence and DNR limits.
- Keep the packaged extension and Go engine functional when the API is not configured or offline.
- Provide the Android client contract through `/v1/config`; do not claim a mobile app is integrated if `apps/mobile` does not exist.

### Stage C — Provider synchronization and protection status

- Make CleanBrowsing synchronization an explicit protected operation, not an unexpected startup side effect.
- Clearly report unavailable API configuration and partial sync failures.
- Ensure clients can tell the configured resolver and blocklist version; do not show a false “protected” state.
- Do not expose provider logs or query history.

### Stage D — Tests, device verification, and documentation

- Test CleanBrowsing API behavior with a fake HTTP transport; never use a real API key in automated tests.
- Test HTTP status codes, ETag handling, auth requirements, provider error handling, invalid domains, and idempotent sync.
- Run Go tests, extension tests/build, and available lint/type checks. Do not report an unrun test as passed.
- Verify live provider sync only when server credentials are configured locally. Never ask the user to paste secrets into chat.
- Update README, API docs, and OpenAPI contract. Document that Android app/device verification is pending if no Android client or emulator exists.

## 6. Acceptance criteria

The backend integration is complete only when:

- The API starts without CleanBrowsing credentials and serves health, DNS config, and the curated feed.
- The API key stays server-side; provider sync requires the admin bearer token, is idempotent, and only adds curated domains.
- The extension refreshes/caches the feed without per-navigation backend calls and continues offline with packaged rules.
- Android has a documented `/v1/config` contract for CleanBrowsing DoT/DoH setup. Do not claim the missing mobile app is complete.
- Automated checks pass, and provider/live-device verification gaps are reported plainly.
- Privacy statements match actual behavior; Vajra does not collect or persist DNS queries or browsing history.

## 7. What to report at completion

Give the user:

1. What architecture was implemented and why.
2. Which files/modules changed.
3. Which automated and device checks passed, with exact commands and results.
4. Any unverified platform behavior, known DNS bypasses, or remaining false-positive risks.
5. The exact Android build/install steps.
6. The next most valuable action, such as testing user-provided miss domains or reviewing false positives.

Do not claim NordVPN parity, universal blocking, or a numeric accuracy rate unless measured against a published, reproducible test corpus.
