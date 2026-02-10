---
sidebar_position: 3
---

# Managing Variables

Learn how to add, update, delete, and organize environment variables.

## Setting Variables

### Basic Set

```bash
# Set a single variable
envm set DATABASE_URL="postgres://localhost:5432/mydb"

# Set multiple variables
envm set API_KEY="abc123" API_SECRET="xyz789"
```

### Set with Options

```bash
# Set and encrypt the value
envm set API_SECRET="super-secret" --encrypt

# Set in a specific group
envm set REDIS_HOST="localhost" --group=cache

# Set with a description
envm set APP_NAME="My App" --description="Application display name"
```

## Getting Variables

### Get Single Variable

```bash
# Get a variable value
envm get DATABASE_URL
# Output: postgres://localhost:5432/mydb

# Get with metadata
envm get DATABASE_URL --verbose
# Output:
# Key: DATABASE_URL
# Value: postgres://localhost:5432/mydb
# Group: database
# Encrypted: false
# Last Modified: 2026-01-27 10:30:00
```

### Get Multiple Variables

```bash
# Get all variables in a group
envm get --group=database

# Get by pattern
envm get "DATABASE_*"
```

### Output Formats

```bash
# JSON output
envm get DATABASE_URL --output=json
# {"key":"DATABASE_URL","value":"postgres://localhost:5432/mydb"}

# Export format (ready for shell)
envm get --output=export
# export DATABASE_URL="postgres://localhost:5432/mydb"
```

## Listing Variables

### Basic List

```bash
# List all variables
envm list

# Output:
# ┌──────────────────┬─────────────────────────────────────┬──────────┐
# │ KEY              │ VALUE                               │ GROUP    │
# ├──────────────────┼─────────────────────────────────────┼──────────┤
# │ APP_NAME         │ My App                              │ app      │
# │ DATABASE_URL     │ postgres://localhost:5432/mydb      │ database │
# │ API_SECRET       │ ********                            │ secrets  │
# └──────────────────┴─────────────────────────────────────┴──────────┘
```

### Filtering

```bash
# List by group
envm list --group=database

# Search by pattern
envm list --search="API"

# Show only keys (no values)
envm list --keys-only

# Show encrypted values (requires authentication)
envm list --show-secrets
```

### Output Formats

```bash
# JSON format
envm list --output=json

# Dotenv format
envm list --output=dotenv

# YAML format
envm list --output=yaml

# Table format (default)
envm list --output=table
```

## Deleting Variables

```bash
# Delete a variable
envm delete API_KEY

# Delete with confirmation prompt
envm delete DATABASE_URL
# Are you sure you want to delete DATABASE_URL? [y/N]

# Force delete without confirmation
envm delete API_KEY --force

# Delete multiple variables
envm delete API_KEY API_SECRET --force
```

## Organizing with Groups

### Create Groups

Groups help organize related variables:

```bash
# Variables are automatically grouped when set with --group
envm set REDIS_HOST="localhost" --group=cache
envm set REDIS_PORT="6379" --group=cache

# Or define in .envm file:
# [cache]
# REDIS_HOST=localhost
# REDIS_PORT=6379
```

### Managing Groups

```bash
# List all groups
envm group list

# Create an empty group
envm group create payments

# Delete a group (moves variables to default)
envm group delete payments

# Rename a group
envm group rename cache redis
```

## Searching Variables

```bash
# Search by key name
envm search "DATABASE"

# Search by value
envm search --value "localhost"

# Regex search
envm search --regex "^API_.*"

# Search across all environments
envm search "SECRET" --all-envs
```

## Bulk Operations

### Import Variables

```bash
# Import from .env file
envm import .env

# Import from JSON
envm import config.json --format=json

# Import from YAML
envm import config.yaml --format=yaml

# Import with conflict resolution
envm import .env --on-conflict=skip     # Skip existing
envm import .env --on-conflict=overwrite # Overwrite existing
envm import .env --on-conflict=prompt    # Ask for each conflict
```

### Export Variables

```bash
# Export to .env format
envm export > .env

# Export specific group
envm export --group=database > database.env

# Export to JSON
envm export --format=json > config.json

# Export to YAML
envm export --format=yaml > config.yaml
```

## Variable Interpolation

Reference other variables using `${VAR}` syntax:

```bash
# In .envm file:
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=myapp
DATABASE_URL=${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}

# Result: DATABASE_URL = localhost:5432/myapp
```

### Default Values

```bash
# Use default if variable is not set
DATABASE_PORT=${DB_PORT:-5432}

# Nested interpolation
CONNECTION_STRING=${${ENV}_DATABASE_URL}
```

## Next Steps

- [Encryption](./encryption) - Learn to encrypt sensitive values
- [Running Commands](./running-commands) - Execute commands with env vars
