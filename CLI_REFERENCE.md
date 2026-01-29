# CLI Command Reference

## Quick Start (Verified)
Follow these steps to initialize the development environment in GitHub Codespaces or a local Linux environment.

1. **Start the Database:**
   `docker compose up -d db`

2. **Configure Environment Credentials:**
   `export DATABASE_URI="postgres://postgres:admin@123@localhost:5432/envm?sslmode=disable"`

3. **Run the Application:**
   `go run internal/cmd/main.go`

4. **Verify System Health:**
   `curl -i http://localhost:5000/health`

---

## Environment Configuration
The application utilizes `godotenv` to manage system configuration.

| Variable | Required | Description | Default (Local Dev) |
| :--- | :--- | :--- | :--- |
| `DATABASE_URI` | Yes | PostgreSQL connection string. | `postgres://postgres:admin@123@localhost:5432/envm?sslmode=disable` |

---

## Service Management (Docker CLI)
The primary orchestration tool for `envm` is Docker Compose.

### `up`
Starts the API server and the database.
* **Command:** `docker compose up [options]`
* **Flag `-d`:** Runs containers in detached mode.
* **Flag `--build`:** Forces a rebuild of the Go binary.

### `build` (Cross-Platform)
To build for specific cloud architectures (e.g., amd64):
`docker build --platform=linux/amd64 -t envm .`

---

## Database Migrations
Database schema changes are managed via `goose`.

**Migration Directory:** `internal/adapters/postgresql/schema/migrations`

### `up`
Applies all pending SQL migrations.
`goose -dir internal/adapters/postgresql/schema/migrations postgres "$DATABASE_URI" up`

### `status`
Checks the current version of the database schema.
`goose -dir internal/adapters/postgresql/schema/migrations postgres "$DATABASE_URI" status`

---

## Development & Troubleshooting

### Building from Source
If building outside of Docker, use the explicit path to the entry point:
`go build -o envm ./internal/cmd/main.go`

### Troubleshooting: Authentication Errors
If you receive `FATAL: password authentication failed for user "postgres"`:
1. Verify `DATABASE_URI` matches the `POSTGRES_PASSWORD` in `compose.yaml`.
2. Default local password is `admin@123`.

### Health Check Specification
**Request:** `curl -i http://localhost:5000/health`
**Expected Response:** `200 OK`