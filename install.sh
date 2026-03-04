#!/bin/bash
set -e

# Repository configuration
REPO_OWNER="envm-org"
REPO_NAME="envm"

# Text formatting
bold=$(tput bold 2>/dev/null || echo "")
green=$(tput setaf 2 2>/dev/null || echo "")
red=$(tput setaf 1 2>/dev/null || echo "")
reset=$(tput sgr0 2>/dev/null || echo "")

echo "${bold}Welcome to the envm installer!${reset}"

# Detect OS and Architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case "${OS}" in
    Linux*)     OS_NAME="Linux" ;;
    Darwin*)    OS_NAME="Darwin" ;;
    CYGWIN*|MINGW*|MSYS*) OS_NAME="Windows" ;;
    *)          echo "${red}Unsupported OS: ${OS}${reset}"; exit 1 ;;
esac

case "${ARCH}" in
    x86_64|amd64) ARCH_NAME="x86_64" ;;
    i386|i686)    ARCH_NAME="i386" ;;
    arm64|aarch64) ARCH_NAME="arm64" ;;
    *)            echo "${red}Unsupported architecture: ${ARCH}${reset}"; exit 1 ;;
esac

echo "${green}==> Detected OS: ${OS_NAME}, Architecture: ${ARCH_NAME}${reset}"

# Dependency checks
if ! command -v curl &> /dev/null; then
    echo "${red}Error: curl is required to download envm. Please install curl and try again.${reset}"
    exit 1
fi

# Fetch latest release
echo "${green}==> Fetching latest release information...${reset}"
# Using GitHub API to get the latest release tag
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "${red}Error: Failed to fetch the latest release. Please check your internet connection or GitHub API limits.${reset}"
    exit 1
fi

echo "${green}==> Latest release: ${LATEST_RELEASE}${reset}"

# Setup download URLs
BINARY_NAME="${REPO_NAME}-${LATEST_RELEASE}"
DOWNLOAD_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${LATEST_RELEASE}/${BINARY_NAME}"

# Download the binary
TMP_DIR=$(mktemp -d)
echo "${green}==> Downloading ${BINARY_NAME} to temporary directory...${reset}"

cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

cd "$TMP_DIR"
curl -sL "$DOWNLOAD_URL" -o "envm"

if [ ! -s "envm" ]; then
    echo "${red}Error: Failed to download ${DOWNLOAD_URL}${reset}"
    exit 1
fi

# Install the binary
GLOBAL_BIN="$HOME/.local/bin"
mkdir -p "$GLOBAL_BIN"
TARGET="$GLOBAL_BIN/envm"

echo "${green}==> Installing to $TARGET...${reset}"

# Remove existing binary if it exists
if [ -f "$TARGET" ]; then
    echo "${green}==> Overwriting existing installation...${reset}"
    rm -f "$TARGET"
fi

cp "envm" "$TARGET"
chmod +x "$TARGET"

echo "${green}${bold}Installation complete!${bold}${reset}"
echo "You can now run 'envm' from your terminal."

# Check PATH
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo ""
    echo "${bold}WARNING: $HOME/.local/bin is not in your PATH.${reset}"
    echo "Add the following line to your shell configuration file (e.g., ~/.bashrc, ~/.zshrc):"
    echo 'export PATH="$HOME/.local/bin:$PATH"'
    echo "Then restart your terminal or run this command."
fi

# Setup completions
echo ""
echo "${green}==> Setting up shell autocompletions...${reset}"

# Setup Bash
if [ -f "$HOME/.bashrc" ]; then
    if ! grep -q "envm completion bash" "$HOME/.bashrc"; then
        echo -e '\n# envm autocompletion\nif command -v envm >/dev/null 2>&1; then\n  source <(envm completion bash)\nfi' >> "$HOME/.bashrc"
        echo "  - Added bash autocompletion to ~/.bashrc"
    fi
fi

# Setup Zsh
if [ -f "$HOME/.zshrc" ]; then
    if ! grep -q "envm completion zsh" "$HOME/.zshrc"; then
        echo -e '\n# envm autocompletion\nif command -v envm >/dev/null 2>&1; then\n  source <(envm completion zsh)\nfi' >> "$HOME/.zshrc"
        echo "  - Added zsh autocompletion to ~/.zshrc"
    fi
fi

echo "${green}${bold}All set! Restart your terminal to use autocompletions!${bold}${reset}"
