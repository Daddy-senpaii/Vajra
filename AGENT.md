# AGENTS.md — Adult Content Blocking Project

## 1. Project Vision

Build a cross-platform adult-content blocking application focused on helping users maintain commitments through reliable content blocking.

The product will eventually consist of:

1. **Go Core Blocking Engine** — shared domain classification and blocking logic.
2. **Browser Extension** — blocks adult domains while browsing.
3. **Website Dashboard** — account management, commitment tracking, settings, and progress.
4. **Android Application** — device-level protection and commitment management.
5. **Go Backend API** — authentication integration, commitment management, blocklist distribution, and synchronization.
6. **Supabase** — PostgreSQL database and authentication.

The architecture should allow these components to work together while keeping their responsibilities separate.

**Current development priority: Build the core blocking engine first.**

---

## 2. Complete Project Structure

Use the following structure as the long-term architectural reference.

```text
adult-blocker/
│
├── AGENTS.md
├── README.md
├── .gitignore
│
├── apps/
│   │
│   ├── extension/                  # Browser extension
│   │   ├── public/
│   │   ├── src/
│   │   │   ├── background/
│   │   │   ├── content/
│   │   │   ├── popup/
│   │   │   ├── blocked-page/
│   │   │   ├── rules/
│   │   │   ├── services/
│   │   │   ├── storage/
│   │   │   └── types/
│   │   ├── manifest.json
│   │   ├── package.json
│   │   └── tsconfig.json
│   │
│   ├── web/                        # Website dashboard
│   │   ├── app/
│   │   ├── components/
│   │   ├── features/
│   │   │   ├── auth/
│   │   │   ├── commitment/
│   │   │   ├── dashboard/
│   │   │   ├── progress/
│   │   │   └── settings/
│   │   ├── lib/
│   │   ├── services/
│   │   ├── types/
│   │   ├── package.json
│   │   └── tsconfig.json
│   │
│   └── mobile/                     # Android application
│       ├── app/
│       ├── src/
│       │   ├── components/
│       │   ├── features/
│       │   │   ├── auth/
│       │   │   ├── commitment/
│       │   │   ├── protection/
│       │   │   ├── dashboard/
│       │   │   └── settings/
│       │   ├── services/
│       │   ├── stores/
│       │   ├── types/
│       │   └── utils/
│       ├── modules/
│       │   └── protection/         # Native Android protection module
│       ├── assets/
│       ├── app.json
│       └── package.json
│
├── packages/
│   │
│   ├── blocking-engine/            # Core Go blocking engine
│   │   ├── cmd/
│   │   │   └── blocklist-tool/     # Optional development utility
│   │   ├── domain/
│   │   │   ├── normalize.go
│   │   │   ├── validate.go
│   │   │   └── normalize_test.go
│   │   ├── classifier/
│   │   │   ├── classifier.go
│   │   │   └── classifier_test.go
│   │   ├── blocklist/
│   │   │   ├── blocklist.go
│   │   │   ├── loader.go
│   │   │   └── blocklist_test.go
│   │   ├── allowlist/
│   │   │   ├── allowlist.go
│   │   │   └── allowlist_test.go
│   │   ├── engine/
│   │   │   ├── engine.go
│   │   │   └── engine_test.go
│   │   ├── testdata/
│   │   ├── benchmarks/
│   │   ├── go.mod
│   │   └── go.sum
│   │
│   └── contracts/                 # Shared API specifications
│       └── openapi.yaml
│
├── services/
│   └── api/                       # Go backend
│       ├── cmd/
│       │   └── api/
│       │       └── main.go
│       ├── internal/
│       │   ├── auth/
│       │   ├── commitment/
│       │   ├── device/
│       │   ├── blocklist/
│       │   ├── protection/
│       │   ├── httpapi/
│       │   ├── database/
│       │   └── config/
│       ├── go.mod
│       └── go.sum
│
├── supabase/
│   ├── migrations/
│   ├── seed.sql
│   └── config.toml
│
├── docs/
│   ├── architecture/
│   ├── api/
│   ├── development/
│   └── decisions/
│
├── scripts/
│
└── .github/
    └── workflows/
```

### Structure Rules

- `packages/blocking-engine` contains reusable Go business logic.
- `apps/extension` contains browser-specific integration and UI.
- `apps/web` contains dashboard functionality only.
- `apps/mobile` contains React Native and Android-specific functionality.
- `services/api` contains backend business operations and API endpoints.
- `supabase` contains database migrations and Supabase configuration.
- `packages/contracts` contains shared API specifications.
- `docs` contains technical documentation and architectural decisions.

This is the target architecture, not a requirement to create every directory immediately.

Create directories only when implementation requires them.

Avoid unnecessary empty folders.

---

## 3. Current Development Phase

### Phase 1 — Core Blocking Engine

Build an independent Go engine that:

- Accepts domain names.
- Normalizes domain names.
- Validates domain inputs.
- Identifies adult domains using a maintained blocklist.
- Supports subdomain matching.
- Supports an allowlist.
- Returns a clear blocking decision.
- Supports blocklist updates.

The engine must work independently without a browser, mobile application, backend, or internet connection.

### Phase 2 — Browser Extension

Integrate the tested engine with browser-specific request-blocking capabilities.

Features:

- Website request blocking.
- Custom blocked page.
- Basic extension popup.
- Blocklist management.
- Protection status.

Do not assume the Go engine can directly run inside a browser extension. Evaluate a suitable integration strategy, including WebAssembly or generating browser-compatible blocking rules.

### Phase 3 — Website Dashboard

Build a minimal dashboard.

Features:

- Registration and login.
- User dashboard.
- Commitment countdown.
- Protection status.
- Basic settings.

### Phase 4 — Backend and Authentication

Build the Go REST API and integrate Supabase.

Features:

- Authentication verification.
- User accounts.
- Commitment management.
- Device registration.
- Blocklist distribution.
- Synchronization.

### Phase 5 — Commitment System

Implement:

- Commitment duration.
- Start and expiry timestamps.
- Protection state.
- Commitment history.
- Progress tracking.
- Milestone reminders.

### Phase 6 — Android Application

Build the React Native application using Expo and Expo Router.

Integrate Android-specific protection through a native Kotlin module using appropriate Android APIs.

### Phase 7 — Reliability and Advanced Protection

Explore:

- Offline synchronization.
- Reboot recovery.
- Protection health monitoring.
- Bypass resistance.
- Advanced content classification.
- Performance and scalability improvements.

---

## 4. Technology Stack

### Core Engine

- Go
- Go standard library where practical
- Go testing package

### Browser Extension

- TypeScript
- Browser Extension APIs
- React if required for the popup interface

### Website

- React / Next.js
- TypeScript

### Android

- React Native
- Expo
- Expo Router
- Zustand
- Kotlin for native Android integration

### Backend

- Go
- REST API
- Supabase PostgreSQL
- Supabase Auth

Keep dependencies minimal and introduce them only when justified.

---

## 5. Core Blocking Engine Requirements

### A. Domain Normalization

The engine must:

- Normalize uppercase and lowercase characters.
- Handle leading and trailing whitespace.
- Handle subdomains.
- Handle internationalized domain names consistently.
- Reject malformed domains.
- Avoid incorrect matching of similar domain names.

### B. Domain Classification

Initially use a locally maintained adult-domain blocklist.

Support:

- Exact domain matching.
- Subdomain matching.
- Efficient lookups.
- Blocklist loading.
- Blocklist updates.
- Duplicate prevention.

Keep classification logic separate from the underlying blocklist storage.

### C. Blocking Decisions

Expected processing flow:

```text
Input domain
     |
Normalize
     |
Validate
     |
Check allowlist
     |
Check blocklist
     |
Return decision
```

Possible results:

- `ALLOW`
- `BLOCK`
- `INVALID`

Allowlist entries take precedence over blocklist entries.

### D. Allowlist

Support:

- Adding domains.
- Removing domains.
- Exact domain matching.
- Defined subdomain matching.
- Allowlist precedence.

### E. Blocklist

Support:

- Adding domains.
- Removing domains.
- Loading local blocklist data.
- Duplicate prevention.
- Efficient domain lookups.
- Future compatibility with backend updates.

Do not implement backend synchronization in Phase 1.

---

## 6. Testing Requirements

Testing is a major priority.

Write unit tests for:

- Domain normalization.
- Domain validation.
- Exact domain matching.
- Subdomain matching.
- Similar but unrelated domain names.
- Allowlist precedence.
- Blocklist updates.
- Empty blocklists.
- Duplicate entries.
- Case-insensitive matching.

Illustrative cases:

| Input | Expected result |
|---|---|
| `adult-example.com` | BLOCK |
| `www.adult-example.com` | BLOCK |
| `example.com` | ALLOW |
| `notadult-example.com` | ALLOW |
| Invalid domain | INVALID |
| Explicitly allowed domain | ALLOW |

Use test fixtures rather than assuming these illustrative domains are actual blocklist entries.

Also include benchmarks for domain lookup performance once the initial implementation is complete.

---

## 7. Architecture Principles

1. Keep the core engine independent of all applications.
2. Separate normalization, validation, classification, and blocking decisions.
3. Avoid unnecessary abstractions.
4. Keep packages small and responsibilities clear.
5. Write idiomatic Go.
6. Prefer testable functions.
7. Avoid global mutable state.
8. Keep data sources separate from business logic.
9. Do not place business logic inside UI components.
10. Do not introduce microservices prematurely.

The engine must be usable by future integrations without depending on their platform-specific code.

---

## 8. Security and Privacy

- Do not collect browsing history.
- Do not store complete URLs unnecessarily.
- Do not transmit browsing activity.
- Validate all domain inputs.
- Avoid executing downloaded blocklist content as code.
- Keep user information separate from blocklist data.
- Do not introduce unnecessary tracking or analytics.

---

## 9. Current Scope Restrictions

For Phase 1, do not implement:

- Authentication.
- User registration.
- Supabase integration.
- Backend API.
- Browser extension.
- Website dashboard.
- Android application.
- VPNService.
- DNS interception.
- AI-based content classification.
- Image or video analysis.
- Anti-uninstall mechanisms.
- Advanced bypass resistance.

These belong to later development phases.

---

## 10. Agent Working Instructions

Before implementation:

1. Read this entire file.
2. Inspect the existing repository.
3. Understand the current code and project structure.
4. Identify missing components.
5. Propose a minimal implementation plan.
6. Explain important technical decisions.

During implementation:

- Work incrementally.
- Write tests alongside functionality.
- Run tests after meaningful changes.
- Explain architectural decisions briefly.
- Avoid unnecessary dependencies.
- Do not rewrite working code without justification.
- Do not implement future phases prematurely.
- Do not create the entire project structure upfront.

Do not begin coding until the user approves the implementation plan.

---

## 11. Definition of Done — Phase 1

The core engine is ready for browser integration when:

- Domain normalization works correctly.
- Domain validation works correctly.
- Adult-domain matching works.
- Subdomain matching works.
- Allowlist precedence works.
- Blocklist updates work.
- Invalid inputs are handled safely.
- Unit tests pass.
- Benchmarks are available.
- The engine works independently without a frontend or internet connection.

**Final objective:** Build a reliable, reusable Go blocking engine that becomes the foundation for browser extension and Android protection.
