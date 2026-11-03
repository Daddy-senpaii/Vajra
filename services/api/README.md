# Vajra API

This Go service exposes a small versioned policy/config API for the browser
extension and future Android client. It does not receive per-navigation or
per-DNS queries. The extension can cache the curated domain feed and continue
using its packaged rules offline. Android can read the CleanBrowsing DNS
configuration and configure Android Private DNS or its approved local VPN
mode.

CleanBrowsing's REST API manages account profiles and explicit domain entries;
the provider's DNS resolver performs category-based adult filtering. The
backend sync endpoint sends only the 243 project-curated domains, not the
upstream 426k-domain snapshot. Use an existing CleanBrowsing profile that is
connected to the custom resolver hostname configured below. Public Adult Filter
defaults do not include account-specific custom entries.

## Run locally

```sh
cp .env.example .env
# Set a random VAJRA_ADMIN_TOKEN with at least 32 characters. Add the
# CleanBrowsing API key only if you will use the protected sync endpoint.
# Never commit actual secrets.
set -a
. ./.env
set +a
go run ./cmd/api
```

The server listens on `:8080` by default. Put it behind HTTPS before exposing
it to devices. Set `VAJRA_API_ALLOWED_ORIGINS` to the exact extension origin(s)
that should read the public endpoints.

## Endpoints

- `GET /healthz`: process health.
- `GET /v1/config`: adult filter provider and DoT/DoH resolver configuration.
- `GET /v1/blocklist`: sorted project-curated delta, SHA-256 version, ETag, and
  count. Conditional requests return `304` when unchanged.
- `POST /v1/admin/cleanbrowsing/sync`: idempotently adds missing curated
  domains to the configured CleanBrowsing profile. Requires
  `Authorization: Bearer $VAJRA_ADMIN_TOKEN`; returns counts and any domains
  that failed. It never removes provider entries.

Example sync command:

```sh
curl -X POST http://localhost:8080/v1/admin/cleanbrowsing/sync \
  -H "Authorization: Bearer $VAJRA_ADMIN_TOKEN"
```

CleanBrowsing requires its API key as an HTTPS query parameter. The client
keeps it server-side, does not log request URLs, and suppresses URL-bearing
transport errors. Configure a long random Vajra admin token and restrict access
to the admin endpoint at the network layer as well.

## Extension connection

Set `VITE_VAJRA_API_URL` to this service's HTTPS base URL when building the
extension. Without it, the extension uses its packaged rules only. When
configured, it fetches the small curated delta on startup and every six hours,
uses ETags, validates domains, and keeps its last cached list if offline.

## Android connection

`GET /v1/config` gives the Android client the configured DoT hostname and DoH
URL. Android 9+ supports Private DNS. The client must clearly display whether
filtering is active and request explicit consent before using local VPN mode.
The current repository does not yet contain `apps/mobile`; this API endpoint is
the contract for that future client, not a claim that a mobile app is already
integrated.

## CleanBrowsing custom profile setup

1. Create or select the CleanBrowsing profile in the dashboard and configure
   its custom DoT/DoH hostnames.
2. Set `CLEANBROWSING_PROFILE_NAME`, `CLEANBROWSING_DOT_HOST`, and
   `CLEANBROWSING_DOH_URL` to that profile's values.
3. Set `CLEANBROWSING_API_KEY` only in the backend environment.
4. Call the protected sync endpoint. The API adds domains but does not delete
   existing CleanBrowsing entries.

The public Adult Filter endpoints are usable without an account, but they do
not use the account-specific additions from the sync endpoint. For mobile,
CleanBrowsing receives DNS queries that are sent to its resolver. Vajra itself
does not store those queries. DNS filtering cannot distinguish safe and adult
pages on the same hostname and cannot stop arbitrary DoH, direct-IP, or VPN
bypasses by itself.
