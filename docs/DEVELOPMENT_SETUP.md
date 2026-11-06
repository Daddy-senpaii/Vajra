# Development setup (WSL 2 + Android Studio)

The supported workflow keeps source in `/mnt/d/Vajra`, uses Ubuntu 24.04 in WSL for command-line development, and uses Android Studio on Windows for Android SDK management, emulators, and native project inspection. Windows sees the project as `D:\Vajra`.

## Prerequisites

- Go 1.26.1 (see `backend/go.mod`).
- Node.js 22.13.0 or newer (see `.nvmrc`; React Native 0.88 and Metro require it).
- JDK 17 for Android builds.
- Android Studio with Android SDK Platform, Build Tools, Platform Tools, and an emulator image installed.
- MongoDB 8.0, either through Docker Desktop with WSL integration, another local MongoDB service, or a MongoDB URI configured in `.env`.

The current WSL environment has Go 1.26.1 and Java 21, but Node 20.20.0; it needs Node 22.13.0+ and JDK 17 selected for local Android builds. MongoDB and Docker are not installed in WSL at present. Android Studio's bundled JBR may be used by Android Studio; for WSL Gradle builds, configure `JAVA_HOME` to a JDK 17 installation.

## One-time WSL setup

1. Install/use Node 22.13.0 or newer and verify `node --version`.
2. Install JDK 17 and set `JAVA_HOME` to it for WSL Gradle builds.
3. Install Android Studio on Windows. In **Settings → Languages & Frameworks → Android SDK**, note the SDK path and install the SDK/platform tools and emulator required by the Gradle project.
4. For Android Studio, open `D:\Vajra\mobile\android`. Keep WSL and Windows on this same checkout.
5. For WSL Gradle/ADB commands, make the Windows SDK reachable from WSL and set `ANDROID_HOME` and `PATH` for that installation. If file watching or Gradle performance is poor on `/mnt/d`, move the repository to the WSL Linux filesystem only after arranging a single shared workflow; do not keep divergent copies.
6. Install Docker Desktop on Windows and enable WSL integration for Ubuntu 24.04, or configure another MongoDB instance. Docker Compose is only the local convenience path.

## Start MongoDB and API

```sh
cp .env.example .env
docker compose up -d mongodb
set -a && . ./.env && set +a
cd backend
go mod download
go run ./cmd/api
```

In another WSL terminal, verify `curl http://localhost:8080/healthz` and `curl http://localhost:8080/readyz`.

## Start Android app

```sh
cd mobile
npm ci
npm start
```

Run the app from Android Studio's Run action after selecting an emulator/device, or from WSL with `npx react-native run-android` once Android SDK/ADB paths are configured. Keep Metro running in a separate terminal. The app is currently the generated React Native launch screen; product screens come in later phases.

## Common checks

```sh
(cd backend && go test ./... && go vet ./...)
(cd mobile && npx tsc --noEmit && npm run lint && npm test -- --runInBand)
```

CI runs these checks and builds the Android debug APK. Device/emulator testing still needs Android Studio or a connected Android device.
