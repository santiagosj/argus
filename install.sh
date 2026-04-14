#!/bin/bash

# Argus Quick Installation Script

set -e

echo "🔍 Argus - Cognitive Security Framework Installer"
echo "=================================================="
echo ""

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.25+ from https://go.dev/dl/"
    exit 1
fi

GO_VERSION=$(go version | grep -oP 'go\d+\.\d+')
echo "✓ Go detected: $GO_VERSION"

# Navigate to directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# Download dependencies
echo ""
echo "📦 Downloading dependencies..."
go mod download

# Build
echo "🔨 Building Argus..."
CGO_ENABLED=1 go build -o argus ./cmd/argus

# Verify
if [ -f "argus" ]; then
    echo "✓ Build successful!"
    echo ""
    echo "Quick start:"
    echo "  ./argus demo              # Run demo"
    echo "  ./argus help              # Show help"
    echo "  ./argus status            # Show status"
    echo ""
    echo "Dashboard: http://localhost:8080"
else
    echo "❌ Build failed"
    exit 1
fi
