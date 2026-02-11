# Development Guide

Welcome to the `envm` project! This guide will help you set up your environment and run the various components of the application.

## Prerequisites

- **Go**: Version 1.22+
- **Node.js**: Version 20+
- **Docker**: For running the database and services.
- **Make**: For running build automation commands.

## Getting Started

### 1. Database Setup

Start the PostgreSQL database using Docker Compose:

```bash
docker compose up -d db
```

Ensure your `.env` file or environment variables are configured with `DATABASE_URI`.
Default for local development:
`postgres://postgres:admin@123@localhost:5432/envm?sslmode=disable`

### 2. Running the Server

You can run the server directly or via Docker.

**Local (Go):**
```bash
make run
```
Or for hot-reloading:
```bash
make dev
```

**Docker:**
```bash
docker compose up server
```

### 3. Running the CLI

The CLI source code is located in the `cli/` directory.

**Build:**
```bash
make build-cli
```
The binary will be placed in `bin/envm`.

**Run from source:**
```bash
cd cli
go run . [command]
```

**Install Script:**
```bash
./scripts/install-cli.sh
```

### 4. Running Documentation

The documentation is built with Docusaurus and is located in `docs/`.

**Install Dependencies & Start:**
```bash
make run-docs
```
This will start the docs server on http://localhost:3000.

## Useful Commands

Run `make help` to see all available commands.

| Command | Description |
| :--- | :--- |
| `make build-all` | Build Server, CLI, and Docs |
| `make test` | Run all tests |
| `make lint` | Run code linters |
| `make format` | Format Go code |
| `make migration name=foo` | Create a new SQL migration |

## Directory Structure

- `cmd/`: Server entry point.
- `internal/`: Private application code.
- `pkg/`: Public library code.
- `cli/`: Command-line interface tool.
- `docs/`: Project documentation site.
- `scripts/`: Utility scripts.
