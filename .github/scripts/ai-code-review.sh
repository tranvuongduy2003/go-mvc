#!/bin/bash

# AI-Assisted Code Review Script
# This script analyzes code changes and provides AI-powered insights

set -e

BRANCH="${1:-$(git branch --show-current)}"
BASE_BRANCH="${2:-master}"

echo "🤖 AI Code Review Assistant"
echo "=========================="
echo "Branch: $BRANCH"
echo "Base: $BASE_BRANCH"
echo ""

# Get changed files
echo "📝 Analyzing changed files..."
CHANGED_FILES=$(git diff --name-only $BASE_BRANCH...$BRANCH)

if [ -z "$CHANGED_FILES" ]; then
    echo "✅ No changes detected"
    exit 0
fi

echo "Changed files:"
echo "$CHANGED_FILES"
echo ""

# Analyze complexity
echo "📊 Analyzing code complexity..."
for file in $CHANGED_FILES; do
    if [[ $file == *.go ]]; then
        echo "Checking: $file"
        
        # Check cyclomatic complexity
        if command -v gocyclo &> /dev/null; then
            gocyclo -over 10 "$file" || true
        fi
        
        # Check line count
        LINES=$(wc -l < "$file")
        if [ "$LINES" -gt 300 ]; then
            echo "⚠️  Warning: $file has $LINES lines (consider splitting)"
        fi
    fi
done
echo ""

# Check for common issues
echo "🔍 Checking for common issues..."

# Check for TODO comments
echo "TODO comments:"
git diff $BASE_BRANCH...$BRANCH | grep -i "TODO" || echo "None found"
echo ""

# Check for FIXME comments
echo "FIXME comments:"
git diff $BASE_BRANCH...$BRANCH | grep -i "FIXME" || echo "None found"
echo ""

# Check for hardcoded credentials (basic check)
echo "🔒 Security check for hardcoded secrets..."
PATTERNS=(
    "password.*=.*['\"].*['\"]"
    "api[_-]?key.*=.*['\"].*['\"]"
    "secret.*=.*['\"].*['\"]"
    "token.*=.*['\"].*['\"]"
)

for pattern in "${PATTERNS[@]}"; do
    if git diff $BASE_BRANCH...$BRANCH | grep -iE "$pattern"; then
        echo "⚠️  Potential hardcoded secret found! Please review."
    fi
done
echo ""

# Check test coverage
echo "🧪 Checking test coverage..."
if command -v go &> /dev/null; then
    for file in $CHANGED_FILES; do
        if [[ $file == *.go ]] && [[ $file != *_test.go ]]; then
            TEST_FILE="${file%.go}_test.go"
            if [ ! -f "$TEST_FILE" ]; then
                echo "⚠️  Missing test file for: $file"
            fi
        fi
    done
fi
echo ""

# Generate summary
echo "📋 Review Summary"
echo "================"
echo "✅ Analysis complete"
echo ""
echo "💡 Recommendations:"
echo "1. Ensure all new functions have tests"
echo "2. Keep functions under 50 lines when possible"
echo "3. Add documentation for exported functions"
echo "4. Use meaningful variable names"
echo "5. Follow the project's coding standards"
echo ""
echo "📚 Next steps:"
echo "1. Review the AI suggestions above"
echo "2. Run 'make lint' for detailed linting"
echo "3. Run 'make test' to ensure all tests pass"
echo "4. Update documentation if needed"
