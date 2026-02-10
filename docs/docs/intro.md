---
sidebar_position: 1
---

# Introduction to ENVM

**ENVM** is a secure environment variable management and sync tool designed for modern development teams. It helps you manage, encrypt, and share environment variables without exposing sensitive data.

## The Problem

In collaborative environments, sharing environment variables is often done through insecure channels:

- Copying `.env` files via Slack or email
- Storing secrets in shared documents
- Manually syncing variables across team members
- No visibility into who has access to what

**ENVM solves this** by providing a secure, centralized way to manage environment variables with fine-grained access control.

## Key Features

### 🔐 Security First

- **AES-256 encryption** at rest and in transit
- **Role-based access control (RBAC)** with granular permissions
- **Audit logging** for all operations
- **Secret scanning** to detect accidentally exposed credentials

### 👥 Team Collaboration

- Share secrets without exposing raw values
- Temporary access with time-limited sharing
- Change request workflows with approval processes
- Real-time notifications for changes

### ⌨️ Powerful CLI

Built with Go and Cobra for maximum performance:

```bash
# Initialize a new .envm file
envm init

# Set an environment variable
envm set DATABASE_URL="postgres://localhost:5432/mydb"

# Get a variable
envm get DATABASE_URL

# List all variables
envm list

# Sync with remote server
envm sync

# Run a command with env vars loaded
envm run -- npm start
```

### 🌍 Multi-Environment Support

Manage multiple environments with inheritance:

```
base.envm          # Shared defaults
├── dev.envm       # Development overrides
├── staging.envm   # Staging overrides  
└── prod.envm      # Production overrides
```

### 📜 Version Control

- Track every change with full history
- Rollback to any previous version
- Diff tool to compare environments
- Dry-run mode to preview changes

### 🔄 Integration & Compatibility

- Import from `.env`, JSON, YAML files
- Export to popular formats
- NPM package for Node.js projects
- CI/CD pipeline integration

## Supported Languages

ENVM provides native support for:

| Language   | Status    |
|------------|-----------|
| JavaScript | ✅ MVP    |
| TypeScript | ✅ MVP    |
| Python     | ✅ MVP    |
| Go         | ✅ MVP    |
| Java       | 🔜 Next   |
| C#         | 🔜 Next   |
| Ruby       | 🔜 Next   |

## The `.envm` File Format

ENVM uses its own file format (`.envm`) that supports:

```bash
# Comments for documentation
@meta version=1.0.0
@meta owner=backend-team

# Group variables by category
[database]
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=myapp

# Encrypted values (marked with ENC:)
[secrets]
DATABASE_PASSWORD=ENC:base64encodedciphertext

# Variable interpolation
[connection]
DATABASE_URL=${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}
```

## Tech Stack

ENVM is built with a **Go-first architecture**:

- **Core**: Go for all business logic, CLI, and backend API
- **CLI**: Cobra framework with shell completion
- **Database**: SQLite (local), PostgreSQL (production)
- **Cache**: Redis for performance
- **API**: gRPC with REST gateway
- **Events**: NATS for real-time sync
- **Frontend**: React web dashboard

## Quick Start

### Installation

```bash
# Using Go
go install github.com/envm-org/envm@latest

# Using npm (wrapper around Go binary)
npm install -g envm

# Using Homebrew (macOS/Linux)
brew install envm-org/tap/envm
```

### Initialize Your Project

```bash
# Create a new .envm file
envm init

# Import existing .env file
envm import .env --format=dotenv
```

### Basic Usage

```bash
# Set a variable
envm set API_KEY="your-api-key" --encrypt

# Run your application with variables loaded
envm run -- npm start
```

## Next Steps

- [Installation Guide](/docs/tutorial-basics/installation) - Install ENVM on your system
- [Initialize a Project](/docs/tutorial-basics/initialize-project) - Set up your first project
- [Managing Variables](/docs/tutorial-basics/manage-variables) - Add and organize environment variables
- [Encryption](/docs/tutorial-basics/encryption) - Protect sensitive values

