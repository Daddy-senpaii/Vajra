# Architecture

## Components

- `mobile/`: React Native 0.88 with TypeScript, Android as the first target.
- `backend/`: Go HTTP API, with transport, business services, repositories, and database setup kept as separate layers as features are added.
- MongoDB: official Go Driver v2.9.0; a single reusable client is created at startup and closed during graceful shutdown.
- `docker-compose.yml`: local MongoDB service for development. CI uses an ephemeral MongoDB service container.

The API initializes checked-in indexes idempotently during startup and exposes process health (`/healthz`) and database readiness (`/readyz`). Times are stored in UTC; challenge dates will be derived using each user's configured timezone.

## Local and Android Studio workflow

The working copy is under `/mnt/d/Vajra` in WSL, corresponding to `D:\Vajra` in Windows. Use WSL for Git, Go, npm, and the Metro bundler. Open `D:\Vajra\mobile\android` in Android Studio on Windows to install Android SDK packages and run an emulator or physical device. Do not copy or maintain a second checkout.

Keep Android signing material and environment secrets out of Git. Set local SDK paths in ignored `mobile/android/local.properties` only if needed.
