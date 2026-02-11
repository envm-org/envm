#!/bin/bash
set -e

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Build the CLI
echo "Building envm..."
cd "$PROJECT_ROOT/cli"
go build -o ../bin/envm .
cd ..

# Install to local bin to avoid sudo
GLOBAL_BIN="$HOME/.local/bin"
mkdir -p "$GLOBAL_BIN"
TARGET="$GLOBAL_BIN/envm"

echo "Installing to $TARGET..."
cp "$PROJECT_ROOT/bin/envm" "$TARGET"

echo "Installation complete!"
echo "You can now run 'envm'"

# Check if ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo "WARNING: $HOME/.local/bin is not in your PATH."
    echo "Add the following line to your shell configuration file (e.g., ~/.bashrc, ~/.zshrc):"
    echo 'export PATH="$HOME/.local/bin:$PATH"'
fi
