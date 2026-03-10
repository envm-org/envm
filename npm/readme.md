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
