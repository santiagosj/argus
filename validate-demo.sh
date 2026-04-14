#!/bin/bash

# Argus Pre-Demo Validation Script
# Run this before presenting to ensure everything works

set -e

echo "🔍 Argus Pre-Demo Validation"
echo "=============================="
echo ""

PASS="✓"
FAIL="✗"
WARN="⚠"

validation_pass=0
validation_fail=0
validation_warn=0

# Helper functions
check_file() {
    if [ -f "$1" ]; then
        echo "$PASS $1"
        ((validation_pass++))
    else
        echo "$FAIL $1 (NOT FOUND)"
        ((validation_fail++))
    fi
}

check_dir() {
    if [ -d "$1" ]; then
        echo "$PASS $1/"
        ((validation_pass++))
    else
        echo "$FAIL $1/ (NOT FOUND)"
        ((validation_fail++))
    fi
}

check_command() {
    if command -v "$1" &> /dev/null; then
        echo "$PASS $1 installed"
        ((validation_pass++))
    else
        echo "$WARN $1 NOT installed (but might not be required)"
        ((validation_warn++))
    fi
}

# Start validation
echo "1️⃣  Checking Required Files"
echo "----"
check_file "README.md"
check_file "go.mod"
check_file "go.sum"
check_file "config.yaml"
check_file "argus"  # Should exist after build
check_file "Dockerfile"
check_file "Makefile"
check_file "install.sh"
echo ""

echo "2️⃣  Checking Directories"
echo "----"
check_dir "cmd"
check_dir "internal"
check_dir "skills"
check_dir "skills/Identify"
check_dir "skills/Protect"
check_dir "skills/Detect"
check_dir "skills/Respond"
check_dir "skills/Recover"
check_dir "skills/Evolve"
echo ""

echo "3️⃣  Checking Documentation"
echo "----"
check_file "README.md"
check_file "README_DEMO.md"
check_file "TESTING.md"
check_file "DEMO_SCRIPT.md"
check_file "DEVELOPMENT.md"
echo ""

echo "4️⃣  Checking Dependencies"
echo "----"
check_command "go"
check_command "sqlite3"
check_command "ollama"
echo ""

echo "5️⃣  Build Test"
echo "----"
if [ -f "argus" ]; then
    echo "$PASS Argus binary exists"
    ((validation_pass++))
    
    # Test if it runs
    if ./argus help > /dev/null 2>&1; then
        echo "$PASS Argus binary executes"
        ((validation_pass++))
    else
        echo "$FAIL Argus binary does not execute"
        ((validation_fail++))
    fi
else
    echo "$WARN Argus binary not found"
    echo "   Run: make build"
    ((validation_warn++))
fi
echo ""

echo "6️⃣  Configuration Test"
echo "----"
if grep -q "provider:" config.yaml; then
    echo "$PASS config.yaml has AI provider"
    ((validation_pass++))
else
    echo "$FAIL config.yaml incomplete"
    ((validation_fail++))
fi

if grep -q "persistence:" config.yaml; then
    echo "$PASS config.yaml has persistence config"
    ((validation_pass++))
else
    echo "$FAIL config.yaml incomplete"
    ((validation_fail++))
fi
echo ""

echo "7️⃣  Skills Availability"
echo "----"
skill_count=$(find skills -name "*.md" | wc -l)
echo "$PASS Found $skill_count skill files"
((validation_pass++))

if [ $skill_count -lt 10 ]; then
    echo "$WARN Only $skill_count skills (consider adding more for demo)"
    ((validation_warn++))
fi
echo ""

echo "8️⃣  Network Prerequisites"
echo "----"
if lsof -i :8080 &> /dev/null; then
    echo "$WARN Port 8080 is already in use!"
    ((validation_warn++))
else
    echo "$PASS Port 8080 is available"
    ((validation_pass++))
fi

if lsof -i :11434 &> /dev/null; then
    echo "$PASS Ollama is running on port 11434"
    ((validation_pass++))
else
    echo "$WARN Ollama not detected (will use fallback or will fail)"
    ((validation_warn++))
fi
echo ""

# Summary
echo "=============================="
echo "📊 Validation Summary"
echo "=============================="
echo "✓ Passed:  $validation_pass"
echo "✗ Failed:  $validation_fail"
echo "⚠ Warnings: $validation_warn"
echo ""

if [ $validation_fail -eq 0 ]; then
    if [ $validation_warn -eq 0 ]; then
        echo "🎉 All checks passed! Ready for demo."
    else
        echo "⚠️  Some warnings. Review above before presenting."
    fi
else
    echo "❌ Some checks failed. Fix before presenting:"
    echo ""
    echo "Quick fixes:"
    echo "  1. Build: make build"
    echo "  2. Check: grep -r 'TODO\|FIXME' internal/"
    echo "  3. Test: ./argus demo"
fi

exit $validation_fail
