#!/bin/bash
set -e

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Build the CLI
echo "Building envm..."
cd "$PROJECT_ROOT"
make build-cli

# Install to local bin to avoid sudo
GLOBAL_BIN="$HOME/.local/bin"
mkdir -p "$GLOBAL_BIN"
TARGET="$GLOBAL_BIN/envm"

echo "Installing to $TARGET..."
cp "$PROJECT_ROOT/cli/bin/cli" "$TARGET"

echo "Installation complete!"
echo "You can now run 'envm'"

# Check if ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo "WARNING: $HOME/.local/bin is not in your PATH."
    echo "Add the following line to your shell configuration file (e.g., ~/.bashrc, ~/.zshrc):"
    echo 'export PATH="$HOME/.local/bin:$PATH"'
fi

# Setup completions
echo ""
echo "Setting up shell autocompletions..."

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

echo "All set! Restart your terminal to use envm!"
