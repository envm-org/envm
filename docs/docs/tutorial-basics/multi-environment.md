---
sidebar_position: 6
---

# Multi-Environment Setup

Learn how to manage development, staging, and production environments.

## Environment Hierarchy

ENVM supports environment inheritance, allowing you to define base configurations and override specific values per environment:

```
base.envm              # Shared defaults
├── development.envm   # Dev-specific overrides
├── staging.envm       # Staging overrides
└── production.envm    # Production overrides
```

## Creating Environments

### Initialize Multiple Environments

```bash
# Create environment files
envm env create development
envm env create staging
envm env create production

# List all environments
envm env list
# Output:
# - development (default)
# - staging
# - production
```

### Environment Files

After creation, you'll have:

```
your-project/
├── .envm                    # Base configuration
├── .envm.development        # Development overrides
├── .envm.staging            # Staging overrides
└── .envm.production         # Production overrides
```

## Configuring Environments

### Base Configuration (`.envm`)

```bash
# .envm - Shared configuration
@meta version=1.0.0

[app]
APP_NAME=My App
LOG_LEVEL=info

[database]
DATABASE_POOL_SIZE=10

[cache]
CACHE_TTL=3600
```

### Development Overrides

```bash
# .envm.development
@inherit base

[app]
APP_DEBUG=true
LOG_LEVEL=debug

[database]
DATABASE_URL=postgres://localhost:5432/myapp_dev

[cache]
REDIS_URL=redis://localhost:6379
```

### Production Overrides

```bash
# .envm.production
@inherit base

[app]
APP_DEBUG=false
LOG_LEVEL=warn

[database]
DATABASE_URL=ENC:v1:aes256gcm:... # Encrypted production URL
DATABASE_POOL_SIZE=50

[cache]
REDIS_URL=ENC:v1:aes256gcm:... # Encrypted Redis URL
```

## Switching Environments

### Default Environment

Set the default environment:

```bash
# Set default environment
envm env default development

# Check current default
envm env default
# Output: development
```

### Run with Specific Environment

```bash
# Run with development environment
envm run --env=development -- npm start

# Run with production environment
envm run --env=production -- npm start
```

### Environment Variable Override

```bash
# Set environment via shell variable
export ENVM_ENV=staging
envm run -- npm start  # Uses staging
```

## Comparing Environments

### Diff Command

```bash
# Compare development vs production
envm diff development production

# Output:
# ┌──────────────────┬─────────────────────┬─────────────────────┐
# │ KEY              │ development         │ production          │
# ├──────────────────┼─────────────────────┼─────────────────────┤
# │ APP_DEBUG        │ true                │ false               │
# │ LOG_LEVEL        │ debug               │ warn                │
# │ DATABASE_URL     │ localhost:5432/dev  │ ********            │
# │ DATABASE_POOL... │ 10                  │ 50                  │
# └──────────────────┴─────────────────────┴─────────────────────┘
```

### Validate Environment

```bash
# Check all required variables are set
envm validate --env=production

# Output:
# ✓ All required variables are set
# ✓ All types are valid
# ✓ No undefined interpolations
```

## Best Practices

### 1. Never Commit Production Secrets

```bash
# .gitignore
.envm.production   # If using unencrypted secrets
```

Or use encryption and commit the encrypted file.

### 2. Use Encryption for Sensitive Values

```bash
# Only encrypt in production
envm set DATABASE_PASSWORD="prod-secret" --env=production --encrypt
```

### 3. Validate Before Deploy

```bash
# Add to CI/CD pipeline
envm validate --env=production --strict
```

### 4. Document Required Variables

```yaml
# .envm.config.yaml
validation:
  required:
    - DATABASE_URL
    - REDIS_URL
    - JWT_SECRET
```

## Next Steps

- [Next Steps](./next-steps) - What to learn next
- [Encryption](./encryption) - Secure production secrets

