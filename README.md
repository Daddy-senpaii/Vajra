# Vayra

Vayra is an Android-first self-improvement app for challenges, progress tracking, analytics, and accountability with friends.

## Project layout

- `mobile/` — React Native + TypeScript Android application.
- `backend/` — Go HTTP API, MongoDB persistence, and index setup.
- `docs/` — product, architecture, API, data model, workflow, and setup documentation.

## Development

See [docs/DEVELOPMENT_SETUP.md](docs/DEVELOPMENT_SETUP.md) for the WSL 2 + Windows Android Studio workflow. Start with [AGENT.md](AGENT.md) for product scope and implementation phases.

```sh
cp .env.example .env
docker compose up -d mongodb
set -a && . ./.env && set +a
cd backend && go run ./cmd/api
```

In another WSL terminal:

```sh
cd mobile
npm ci
npm start
```

Open `mobile/android` in Android Studio on Windows to install SDK components, start an emulator, and build/run the app. The local source is under `/mnt/d/Vajra`, so it is visible to both WSL and Windows.

The API health endpoints are `GET /healthz` and `GET /readyz`. See [docs/API.md](docs/API.md).
