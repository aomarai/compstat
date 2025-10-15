#!/bin/bash
# Install compression codecs for benchmarking

set -e

echo "Installing compression codecs..."

# Detect OS
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Detect Linux distribution
    if [ -f /etc/debian_version ]; then
        # Debian/Ubuntu
        echo "Detected Debian/Ubuntu"
        sudo apt-get update
        sudo apt-get install -y zstd xz-utils pigz lz4 pbzip2 brotli
    elif [ -f /etc/redhat-release ]; then
        # RHEL/CentOS/Fedora
        echo "Detected RHEL/CentOS/Fedora"
        sudo yum install -y zstd xz pigz lz4 pbzip2 brotli
    elif [ -f /etc/arch-release ]; then
        # Arch Linux
        echo "Detected Arch Linux"
        sudo pacman -S --noconfirm zstd xz pigz lz4 pbzip2 brotli
    else
        echo "Unknown Linux distribution. Please install manually:"
        echo "  - zstd"
        echo "  - xz-utils (or xz)"
        echo "  - pigz"
        echo "  - lz4"
        echo "  - pbzip2"
        echo "  - brotli"
        exit 1
    fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    echo "Detected macOS"
    if ! command -v brew &> /dev/null; then
        echo "Homebrew not found. Please install from https://brew.sh"
        exit 1
    fi
    brew install zstd xz pigz lz4 pbzip2 brotli
else
    echo "Unsupported OS: $OSTYPE"
    echo "Please install the following tools manually:"
    echo "  - zstd"
    echo "  - xz"
    echo "  - pigz"
    echo "  - lz4"
    echo "  - pbzip2"
    echo "  - brotli"
    exit 1
fi

echo ""
echo "✓ Installation complete!"
echo ""
echo "Installed versions:"
echo "===================="
zstd --version 2>/dev/null || echo "zstd: not found"
xz --version 2>/dev/null | head -1 || echo "xz: not found"
pigz --version 2>&1 | head -1 || echo "pigz: not found"
lz4 --version 2>/dev/null || echo "lz4: not found"
pbzip2 -V 2>&1 | head -1 || echo "pbzip2: not found"
brotli --version 2>/dev/null || echo "brotli: not found"