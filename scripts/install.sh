#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Default to "latest" release if no argument is provided
DEFAULT_VERSION=$(curl -fsSL https://api.github.com/repos/struckchure/go-alchemy/releases/latest | jq -r '.tag_name')

echo "Using version: $VERSION"

# Get the version from the command line argument or use default
VERSION=${1:-$DEFAULT_VERSION}

# Define the base URL for the release artifacts
BASE_URL="https://github.com/struckchure/go-alchemy/releases/download/${VERSION}"

# Define the file names (adjust these as needed)
LINUX_AMD64="go-alchemy_Linux_x86_64.tar.gz"
LINUX_ARM64="go-alchemy_Linux_arm64.tar.gz"
LINUX_I386="go-alchemy_Linux_i386.tar.gz"
MACOS_AMD64="go-alchemy_Darwin_x86_64.tar.gz"
MACOS_ARM64="go-alchemy_Darwin_arm64.tar.gz"

# Determine the OS and architecture
OS=$(uname -s)
ARCH=$(uname -m)

# Set the file to download based on OS and architecture
case "$OS" in
    Linux)
        case "$ARCH" in
            x86_64)
                FILE="$LINUX_AMD64"
                ;;
            aarch64)
                FILE="$LINUX_ARM64"
                ;;
            i386)
                FILE="$LINUX_I386"
                ;;
            *)
                echo "Unsupported architecture: $ARCH"
                exit 1
                ;;
        esac
        ;;
    Darwin)
        case "$ARCH" in
            x86_64)
                FILE="$MACOS_AMD64"
                ;;
            arm64)
                FILE="$MACOS_ARM64"
                ;;
            *)
                echo "Unsupported architecture: $ARCH"
                exit 1
                ;;
        esac
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Define the destination directory
DEST_DIR="$HOME/.go-alchemy/bin"

# Create the destination directory if it does not exist
mkdir -p "$DEST_DIR"

# Download the file
echo "Downloading $FILE ... $BASE_URL/${FILE}"
curl -fsSL "${BASE_URL}/${FILE}" -o "${FILE}"

# Extract the downloaded file to the .storm directory
echo "Extracting $FILE to $DEST_DIR..."
tar -xzf "$FILE" -C "$DEST_DIR"

# Remove the downloaded file
rm "$FILE"

# Add .go-alchemy/bin to PATH if not already present
PATH_ENTRY="export PATH=\"$DEST_DIR:\$PATH\""

# Determine the appropriate profile file based on the operating system
if [ -f /etc/alpine-release ]; then
    PROFILE_FILE="$HOME/.profile"  # Alpine Linux
elif [ "$(uname)" = "Darwin" ]; then
    PROFILE_FILE="$HOME/.zshrc"    # macOS (assuming Zsh)
elif [ "$(uname)" = "Linux" ]; then
    if [ -n "$BASH_VERSION" ]; then
        PROFILE_FILE="$HOME/.bashrc"  # General Linux with Bash
    elif [ -n "$ZSH_VERSION" ]; then
        PROFILE_FILE="$HOME/.zshrc"   # General Linux with Zsh
    else
        PROFILE_FILE="$HOME/.profile" # Fallback
    fi
else
    PROFILE_FILE="$HOME/.profile"     # Default fallback
fi

# Add PATH_ENTRY to the profile file if not already present
if ! grep -qF "$PATH_ENTRY" "$PROFILE_FILE"; then
    echo "$PATH_ENTRY" >> "$PROFILE_FILE"
    echo "Added $DEST_DIR to PATH in $PROFILE_FILE"
else
    echo "$DEST_DIR is already in PATH in $PROFILE_FILE"
fi


if ! grep -Fxq "$PATH_ENTRY" "$PROFILE_FILE"; then
    echo "$PATH_ENTRY" >> "$PROFILE_FILE"
    echo "Updated $PROFILE_FILE to include $DEST_DIR in PATH"
else
    echo "$DEST_DIR is already in PATH in $PROFILE_FILE"
fi