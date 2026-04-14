#!/bin/bash

# QUICK START - Argus Demo Setup (2 minutes)
# Run this script once to validate everything is ready

set -e

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║          🔍 ARGUS - Quick Setup & Validation              ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

step=1

# Helper functions
check() {
    local msg=$1
    echo -n "${BLUE}[$step]${NC} $msg ... "
    ((step++))
}

success() {
    echo -e "${GREEN}✓${NC}"
}

error() {
    echo -e "${RED}✗ $1${NC}"
    exit 1
}

warn() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# 1. Check prerequisites
echo "${BLUE}Prerequisites${NC}"
check "Go 1.25+"
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | grep -oP 'go\d+\.\d+')
    success
    echo "     Found: $GO_VERSION"
else
    error "Go not installed. Install from https://go.dev/dl/"
fi

check "Git"
if command -v git &> /dev/null; then
    success
else
    warn "Git not found (optional)"
fi

# 2. Build
echo ""
echo "${BLUE}Building${NC}"
check "Building Argus binary"
if make build > /tmp/build.log 2>&1; then
    success
    [ -f argus ] && echo "     Binary: $(ls -lh argus | awk '{print $5}')"
else
    error "Build failed. Check: cat /tmp/build.log"
fi

# 3. Verify binary
echo ""
echo "${BLUE}Verification${NC}"
check "Binary executable"
if [ -f argus ] && [ -x argus ]; then
    success
else
    error "Binary not executable"
fi

check "CLI works"
if ./argus help > /dev/null 2>&1; then
    success
else
    error "CLI not working"
fi

check "Config file"
if [ -f config.yaml ]; then
    success
else
    error "config.yaml missing"
fi

check "Skills loaded"
SKILL_COUNT=$(find skills -name "*.md" | wc -l)
if [ $SKILL_COUNT -gt 5 ]; then
    success
    echo "     Skills: $SKILL_COUNT found"
else
    warn "$SKILL_COUNT skills (recommend 10+)"
fi

# 4. Database
echo ""
echo "${BLUE}Database${NC}"
check "SQLite available"
if command -v sqlite3 &> /dev/null; then
    success
else
    warn "SQLite not in PATH (will use bundled)"
fi

# 5. Final checks
echo ""
echo "${BLUE}System Status${NC}"
check "Port 8080 free"
if ! lsof -i :8080 &> /dev/null 2>&1; then
    success
else
    warn "Port 8080 in use (stop it before running demo)"
fi

check "Disk space"
if [ $(df . | tail -1 | awk '{print $4}') -gt 102400 ]; then
    success
    echo "     Available: $(df -h . | tail -1 | awk '{print $4}')"
else
    error "Low disk space"
fi

# 6. Done
echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║                                                            ║"
echo -e "║  ${GREEN}✓ SETUP COMPLETE - Ready for Demo${NC}                        ║"
echo "║                                                            ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "Next steps:"
echo ""
echo "  1. ${BLUE}Start demo:${NC}"
echo "     ./argus demo"
echo ""
echo "  2. ${BLUE}Open dashboard:${NC}"
echo "     http://localhost:8080"
echo ""
echo "  3. ${BLUE}Validate everything:${NC}"
echo "     ./validate-demo.sh"
echo ""
echo "  4. ${BLUE}Read guide:${NC}"
echo "     cat DEMO_SCRIPT.md"
echo ""
echo "Press Enter to run demo, or Ctrl+C to exit..."
read -r

echo ""
echo "🚀 Starting Argus Demo..."
echo ""

./argus demo
