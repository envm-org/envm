---
sidebar_position: 1
---

# Installation

Learn how to install ENVM on your system.

## Prerequisites

Before installing ENVM, ensure you have one of the following:

- **Go 1.21+** (for installing from source)
- **Node.js 18+** (for npm installation)
- **Homebrew** (for macOS/Linux)

## Installation Methods

### Using Go (Recommended)

If you have Go installed, this is the fastest way to get ENVM:

```bash
go install github.com/envm-org/envm@latest
```

Verify the installation:

```bash
envm --version
```

### Using npm

For Node.js projects, you can install ENVM via npm:

```bash
# Global installation
npm install -g envm

# Or as a dev dependency
npm install --save-dev envm
```

The npm package is a thin wrapper that downloads and manages the Go binary automatically.

### Using Homebrew (macOS/Linux)

```bash
# Add the tap
brew tap envm-org/tap

# Install ENVM
brew install envm
```

### Using Docker

For containerized environments:

```bash
docker pull envm/envm:latest

# Run ENVM in a container
docker run --rm -v $(pwd):/app envm/envm list
```

### Manual Installation

Download the pre-built binary for your platform from the [releases page](https://github.com/envm-org/envm/releases):

1. Download the appropriate archive for your OS and architecture
2. Extract the binary
3. Move it to a directory in your PATH

```bash
# Example for Linux amd64
curl -LO https://github.com/envm-org/envm/releases/latest/download/envm_linux_amd64.tar.gz
tar -xzf envm_linux_amd64.tar.gz
sudo mv envm /usr/local/bin/
```

## Shell Completion

ENVM supports shell completion for faster command entry:

```bash
# Bash
envm completion bash > /etc/bash_completion.d/envm

# Zsh
envm completion zsh > "${fpath[1]}/_envm"

# Fish
envm completion fish > ~/.config/fish/completions/envm.fish

# PowerShell
envm completion powershell > envm.ps1
```

## Verify Installation

After installation, verify everything is working:

```bash
# Check version
envm --version

# Show help
envm --help

# Run doctor to check system compatibility
envm doctor
```

## Next Steps

Now that ENVM is installed, learn how to [initialize your first project](./initialize-project).
