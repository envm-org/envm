<p align="center">
  <img src="https://github.com/envm-org/envm/blob/dev/assets/logo.png?raw=true" alt="ENVM Logo" width="200"/>

</p>

<p align="center">
 <h1>ENVM - Environment Variable Management</h1>
</p>

<p align="center">

[![License-MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go-Version](https://img.shields.io/badge/Go-1.25.0-blue.svg)](https://golang.org/dl/)
[![Docker-Ready](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)

</p>

**ENVM** is a powerful environment variable management and synchronization tool designed to help developers and teams securely manage, share, and sync environment variables across projects and environments. Built with security as a top priority, ENVM provides encryption at rest, role-based access control, and seamless collaboration features.

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
- [Documentation](#documentation)
- [Contributing](#contributing)

# Features

- **Secure**: AES-256 encryption, audit logging, and RBAC
- **Team Collaboration**: Share sensitive environment variables safely without exposure
- **Multi-Environment Support**: Manage dev, staging, and production environments effortlessly
- **API-First**: RESTful API for programmatic access
- **CLI Tool**: Command-line interface for CI/CD integration

# Installation

```sh
# Clone the repository
git clone https://github.com/envm-org/envm.git
cd envm
# 
# Start the application services
docker compose up -d
```

## Configuration

Create a `.env` file in the root directory:

```bash
DATABASE_URI=postgres://postgres:password@localhost:5432/envm?sslmode=disable
TOKEN_SECRET=your-super-secret-jwt-key
ADDR=:5000
ENCRYPTION_KEY=your-32-byte-encryption-key
```

# Usage

```sh
envm COMMAND
# runs the command
envm --help COMMAND
# outputs help for specific command
```

# Commands

<!-- commands -->
* [`envm env`](#envm-env)
* [`envm init`](#envm-init)
* [`envm load`](#envm-load)
* [`envm login`](#envm-login)
* [`envm logout`](#envm-logout)
* [`envm project`](#envm-project)
* [`envm pull`](#envm-pull)
* [`envm push`](#envm-push)
* [`envm register`](#envm-register)
* [`envm users`](#envm-users)
* [`envm whoami`](#envm-whoami)

## `envm env`

Manage environment variables

```
USAGE
  $ envm env [COMMAND]

DESCRIPTION
  Create, list, or delete environment variables.
```

## `envm init`

Initialize ENVM configuration in the current directory

```
USAGE
  $ envm init

DESCRIPTION
  Creates an .envm file and links to your current project.
```

## `envm load`

Load environment variables into the current shell

```
USAGE
  $ envm load

DESCRIPTION
  Sources environment variables managed by ENVM into your session.
```

## `envm login`

log in with your ENVM account

```
USAGE
  $ envm login

DESCRIPTION
  Authenticates you with the ENVM server.
```

## `envm logout`

log out

```
USAGE
  $ envm logout

DESCRIPTION
  Clears the local authentication token.
```

## `envm project`

Manage projects

```
USAGE
  $ envm project [COMMAND]

DESCRIPTION
  Create, read, or configure your projects.
```

## `envm pull`

Pull environment variables from the server

```
USAGE
  $ envm pull

DESCRIPTION
  Fetches the latest variables for the current project.
```

## `envm push`

Push environment variables to the server

```
USAGE
  $ envm push

DESCRIPTION
  Uploads local changes to the configured environment.
```

## `envm register`

Register a new ENVM account

```
USAGE
  $ envm register

DESCRIPTION
  Creates a new user profile on the ENVM server.
```

## `envm users`

Manage users

```
USAGE
  $ envm users [COMMAND]

DESCRIPTION
  Administer users on the platform.
```

## `envm whoami`

show the username you are logged in as

```
USAGE
  $ envm whoami

DESCRIPTION
  Displays your active session details.
```

# Documentation

- **Documentation**: [envm-docs.vercel.app](https://envm-docs.vercel.app/)
- **CLI Reference**: [CLI Reference](CLI_REFERENCE.md)
- **Issues**: [GitHub Issues](https://github.com/envm-org/envm/issues)
- **Email**: support@envm.dev

# Contributing

We welcome contributions from the community! 

- Read our [Contributing Guidelines](CONTRIBUTING.md) to get started
- Check our [Code of Conduct](CODE_OF_CONDUCT.md)
- Browse [open issues](https://github.com/envm-org/envm/issues) or create a new one

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Security

Security is paramount for ENVM. If you discover a security vulnerability, please follow our [Security Policy](SECURITY.md) for responsible disclosure. **Do not report security vulnerabilities through public GitHub issues.**
