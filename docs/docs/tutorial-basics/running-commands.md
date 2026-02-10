---
sidebar_position: 5
---

# Running Commands

Learn how to run commands with ENVM environment variables.

## The `run` Command

Execute any command with your `.envm` variables loaded:

```bash
envm run -- <your-command>
```

### Basic Usage

```bash
# Run Node.js app
envm run -- node app.js

# Run npm scripts
envm run -- npm start

# Run Python scripts
envm run -- python main.py

# Run Go applications
envm run -- go run main.go
```

### How It Works

1. ENVM loads variables from `.envm`
2. Decrypts any encrypted values
3. Resolves variable interpolation
4. Injects variables into the command's environment
5. Executes the command

## Environment Selection

### Default Environment

```bash
# Uses the default environment from .envm.config.yaml
envm run -- npm start
```

### Specific Environment

```bash
# Run with development environment
envm run --env=development -- npm start

# Run with production environment
envm run --env=production -- npm start

# Run with staging environment
envm run --env=staging -- npm start
```

### Environment Inheritance

Environment variables cascade:

```
base.envm           # Base configuration
└── development.envm   # Inherits and overrides base
```

```bash
# development inherits from base
envm run --env=development -- npm start
```

## Using with Package Scripts

### npm / yarn / pnpm

```json
{
  "scripts": {
    "start": "envm run -- node dist/index.js",
    "dev": "envm run --env=development -- nodemon src/index.ts",
    "test": "envm run --env=test -- jest"
  }
}
```

### Makefile

```makefile
.PHONY: dev prod test

dev:
	envm run --env=development -- go run cmd/server/main.go

prod:
	envm run --env=production -- ./server

test:
	envm run --env=test -- go test ./...
```

## Advanced Options

### Dry Run

Preview variables without executing:

```bash
envm run --dry-run -- npm start
# Shows all variables that would be set
```

### Verbose Mode

See what's happening:

```bash
envm run --verbose -- npm start
# [envm] Loading .envm
# [envm] Decrypting 3 encrypted values
# [envm] Resolving 2 interpolations
# [envm] Setting 15 environment variables
# [envm] Executing: npm start
```

### Override Variables

Override specific variables for a single run:

```bash
# Override DATABASE_URL for this run only
envm run --set DATABASE_URL="postgres://test:5432/test" -- npm test

# Multiple overrides
envm run --set APP_DEBUG=true --set LOG_LEVEL=debug -- npm start
```

### Subset of Variables

Load only specific variables or groups:

```bash
# Load only database variables
envm run --group=database -- npm start

# Load specific variables
envm run --only="DATABASE_URL,REDIS_URL" -- npm start

# Exclude certain variables
envm run --exclude="DEBUG_*" -- npm start
```

## Shell Integration

### One-liner Export

Export all variables to your current shell:

```bash
# Bash/Zsh
eval $(envm export --format=shell)

# Fish
envm export --format=fish | source

# PowerShell
envm export --format=powershell | Invoke-Expression
```

### Direnv Integration

Use with [direnv](https://direnv.net/) for automatic loading:

```bash
# .envrc
eval $(envm export --format=shell)
```

## CI/CD Integration

### GitHub Actions

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install ENVM
        run: go install github.com/envm-org/envm@latest
      
      - name: Run tests
        env:
          ENVM_MASTER_PASSWORD: ${{ secrets.ENVM_MASTER_PASSWORD }}
        run: envm run --env=test -- npm test
```

### GitLab CI

```yaml
test:
  script:
    - go install github.com/envm-org/envm@latest
    - envm run --env=test -- npm test
  variables:
    ENVM_MASTER_PASSWORD: $ENVM_MASTER_PASSWORD
```

### Docker

```dockerfile
FROM golang:1.21-alpine AS envm
RUN go install github.com/envm-org/envm@latest

FROM node:20-alpine
COPY --from=envm /go/bin/envm /usr/local/bin/envm
COPY . /app
WORKDIR /app

CMD ["envm", "run", "--", "node", "dist/index.js"]
```

## Comparison with dotenv

| Feature | ENVM | dotenv |
|---------|------|--------|
| Load variables | `envm run -- npm start` | `node -r dotenv/config app.js` |
| Encryption | ✓ Built-in (AES-256) | ✗ Not available |
| Multi-environment | ✓ Native support | ⚠️ Manual setup |
| Interpolation | ✓ `${VAR}` syntax | ⚠️ Limited |
| Language agnostic | ✓ Works with any command | ✗ JavaScript only |

## Next Steps

- [Multi-Environment](./multi-environment) - Manage dev, staging, and production
- [Next Steps](./next-steps) - What to learn next

