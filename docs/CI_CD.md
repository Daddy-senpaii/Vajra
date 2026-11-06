# CI/CD

GitHub Actions is the CI provider because the repository remote is GitHub. Workflows are versioned under `.github/workflows/`.

- `backend-ci.yml`: Go dependency verification, formatting, vet, race-enabled tests, build, and an ephemeral MongoDB service for opt-in integration tests.
- `mobile-ci.yml`: reproducible npm install, TypeScript check, lint, Jest, and Android debug build.
- `release.yml`: version-tag builds Go and Android release artifacts. Android output is unsigned. No production deployment or Google Play upload is configured.

The MongoDB integration test runs when `VAYRA_TEST_MONGODB_URI` is set and is skipped in local unit-only runs. CI provides a temporary MongoDB service and does not use production credentials.

## Later deployment setup

When a hosting provider is selected, add a separate protected deployment job that consumes a versioned artifact, deploys only from an approved release/tag, checks the configured health endpoint, and reports status. The Android release job builds an unsigned AAB until all signing secrets are configured. To enable signing, add `ANDROID_KEYSTORE_BASE64`, `ANDROID_KEYSTORE_PASSWORD`, `ANDROID_KEY_ALIAS`, and `ANDROID_KEY_PASSWORD` as repository or protected-environment secrets. Never store signing keys or passwords in source control. Development, staging, and production settings must remain separate.

Production deployment is intentionally not configured at this foundation stage.
