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
for dep in curl tar; do
    if ! command -v "$dep" &> /dev/null; then
        echo "${red}Error: $dep is required but not installed. Please install it and try again.${reset}"
        exit 1
    fi
done

# Determine release version
if [ -n "$VERSION" ]; then
    echo "${green}==> Using specified version: ${VERSION}${reset}"
    LATEST_RELEASE="$VERSION"
else
    echo "${green}==> Fetching latest release information...${reset}"
    LATEST_RELEASE=$(curl -sf "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" \
        | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
fi

if [ -z "$LATEST_RELEASE" ]; then
    echo "${red}Error: Failed to fetch the latest release. Check your internet connection or GitHub API limits.${reset}"
    exit 1
fi

echo "${green}==> Latest release: ${LATEST_RELEASE}${reset}"

# Build archive name matching GoReleaser format:
#   envm_Linux_x86_64.tar.gz
#   envm_Darwin_arm64.tar.gz
#   envm_Windows_x86_64.zip
if [ "$OS_NAME" = "Windows" ]; then
    ARCHIVE_EXT="zip"
else
    ARCHIVE_EXT="tar.gz"
fi

ARCHIVE_NAME="${REPO_NAME}_${OS_NAME}_${ARCH_NAME}.${ARCHIVE_EXT}"
DOWNLOAD_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${LATEST_RELEASE}/${ARCHIVE_NAME}"

echo "${green}==> Downloading ${ARCHIVE_NAME}...${reset}"

# Download to a temp directory
TMP_DIR=$(mktemp -d)
cleanup() { rm -rf "$TMP_DIR"; }
trap cleanup EXIT

cd "$TMP_DIR"

if ! curl -sfL "$DOWNLOAD_URL" -o "$ARCHIVE_NAME"; then
    echo "${red}Error: Failed to download ${DOWNLOAD_URL}${reset}"
    echo "Check that the release exists: https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/tag/${LATEST_RELEASE}"
    exit 1
fi

# Extract archive
echo "${green}==> Extracting...${reset}"
if [ "$ARCHIVE_EXT" = "zip" ]; then
    unzip -q "$ARCHIVE_NAME"
else
    tar -xzf "$ARCHIVE_NAME"
fi

# Install the binary
GLOBAL_BIN="$HOME/.local/bin"
mkdir -p "$GLOBAL_BIN"
TARGET="$GLOBAL_BIN/envm"

echo "${green}==> Installing to $TARGET...${reset}"

if [ -f "$TARGET" ]; then
    echo "${green}==> Overwriting existing installation...${reset}"
    rm -f "$TARGET"
fi

cp "envm" "$TARGET"
chmod +x "$TARGET"

echo "${green}${bold}Installation complete!${reset}"
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
        echo -e '\n# envm autocompletion\nif command -v envm > /dev/null 2>&1; then\n  source <(envm completion bash)\nfi' >> "$HOME/.bashrc"
        echo "  - Added bash autocompletion to ~/.bashrc"
    fi
fi

# Setup Zsh
if [ -f "$HOME/.zshrc" ]; then
    if ! grep -q "envm completion zsh" "$HOME/.zshrc"; then
        echo -e '\n# envm autocompletion\nif command -v envm > /dev/null 2>&1; then\n  source <(envm completion zsh)\nfi' >> "$HOME/.zshrc"
        echo "  - Added zsh autocompletion to ~/.zshrc"
    fi
fi

echo "${green}${bold}All set! Restart your terminal to use envm!${reset}"
