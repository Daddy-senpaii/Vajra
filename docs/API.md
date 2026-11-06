# API

Base URL for local development: `http://localhost:8080`.

## Health

### `GET /healthz`

Confirms that the HTTP process is responding. It does not query MongoDB.

```json
{"status":"ok"}
```

### `GET /readyz`

Pings MongoDB and returns `200` when ready or `503` when the database cannot be reached.

```json
{"status":"ready"}
```

Application endpoints will be documented here as features are implemented. Keep request validation and authorization in the backend.
