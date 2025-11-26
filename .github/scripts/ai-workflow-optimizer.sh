#!/bin/bash

# AI Development Workflow Optimizer
# Analyzes git workflow and suggests optimizations

set -e

echo "🤖 AI Development Workflow Optimizer"
echo "===================================="
echo ""

# Analyze commit patterns
echo "📊 Analyzing commit patterns..."
echo ""

# Get commits from last 30 days
COMMITS=$(git log --since="30 days ago" --pretty=format:"%h|%an|%s|%ad" --date=short)

if [ -z "$COMMITS" ]; then
    echo "No commits found in the last 30 days"
    exit 0
fi

# Count commits per author
echo "👥 Commits by author (last 30 days):"
git log --since="30 days ago" --pretty=format:"%an" | sort | uniq -c | sort -rn
echo ""

# Analyze commit messages
echo "📝 Commit message analysis:"
CONVENTIONAL_COMMITS=$(git log --since="30 days ago" --pretty=format:"%s" | grep -E "^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert):" | wc -l)
TOTAL_COMMITS=$(git log --since="30 days ago" --oneline | wc -l)

if [ "$TOTAL_COMMITS" -gt 0 ]; then
    PERCENTAGE=$(echo "scale=2; $CONVENTIONAL_COMMITS * 100 / $TOTAL_COMMITS" | bc)
    echo "✅ Conventional commits: $CONVENTIONAL_COMMITS / $TOTAL_COMMITS ($PERCENTAGE%)"
    
    if (( $(echo "$PERCENTAGE < 80" | bc -l) )); then
        echo "💡 Tip: Consider using conventional commit format for better changelog generation"
    fi
else
    echo "No commits to analyze"
fi
echo ""

# Analyze file change patterns
echo "📁 Most frequently changed files:"
git log --since="30 days ago" --name-only --pretty=format: | sort | uniq -c | sort -rn | head -10
echo ""

# Analyze branch patterns
echo "🌿 Branch analysis:"
BRANCH_COUNT=$(git branch -a | wc -l)
echo "Total branches: $BRANCH_COUNT"

STALE_BRANCHES=$(git branch -a --merged master | grep -v master | wc -l)
echo "Merged branches (can be cleaned): $STALE_BRANCHES"
echo ""

# Analyze PR patterns (if GitHub CLI is available)
if command -v gh &> /dev/null; then
    echo "📋 Pull Request metrics:"
    
    # Get PR stats from last 30 days
    OPEN_PRS=$(gh pr list --state open --json number | jq '. | length')
    CLOSED_PRS=$(gh pr list --state closed --limit 30 --json number | jq '. | length')
    
    echo "Open PRs: $OPEN_PRS"
    echo "Recently closed PRs: $CLOSED_PRS"
    
    # Average time to merge
    echo ""
    echo "⏱️  Average PR lifetime:"
    gh pr list --state closed --limit 10 --json createdAt,closedAt | \
        jq -r '.[] | (.closedAt | fromdateiso8601) - (.createdAt | fromdateiso8601)' | \
        awk '{sum+=$1; count++} END {if(count>0) print sum/count/86400 " days"}' || echo "N/A"
fi
echo ""

# Analyze test coverage trends
echo "🧪 Testing analysis:"
if [ -f "coverage.out" ]; then
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    echo "Current test coverage: $COVERAGE"
    
    if (( $(echo "$COVERAGE" | sed 's/%//' | awk '{print ($1 < 70)}') )); then
        echo "⚠️  Warning: Coverage below 70%. Consider adding more tests."
    fi
else
    echo "No coverage report found. Run 'make test' to generate one."
fi
echo ""

# Analyze build times
echo "⚡ Build performance:"
if [ -f ".github/workflows/ci.yml" ]; then
    echo "CI/CD configuration found"
    
    if command -v gh &> /dev/null; then
        echo "Recent workflow runs:"
        gh run list --limit 5 --json conclusion,startedAt,updatedAt,name | \
            jq -r '.[] | "\(.name): \(.conclusion) (Duration: \((.updatedAt | fromdateiso8601) - (.startedAt | fromdateiso8601))/60 | floor) min)"'
    fi
fi
echo ""

# Generate recommendations
echo "💡 AI-Powered Recommendations"
echo "============================="
echo ""

# Check for large files
LARGE_FILES=$(git ls-files | xargs ls -l 2>/dev/null | awk '$5 > 1000000 {print $9, $5/1024/1024 "MB"}')
if [ -n "$LARGE_FILES" ]; then
    echo "⚠️  Large files detected (>1MB):"
    echo "$LARGE_FILES"
    echo "Consider using Git LFS for large files"
    echo ""
fi

# Check for missing documentation
if [ ! -f "CONTRIBUTING.md" ]; then
    echo "📚 Missing CONTRIBUTING.md - consider adding contribution guidelines"
fi

if [ ! -f ".github/PULL_REQUEST_TEMPLATE.md" ]; then
    echo "📋 Missing PR template - consider adding one for consistency"
fi

if [ ! -f ".github/issue-template" ]; then
    echo "🐛 Missing issue templates - consider adding them"
fi
echo ""

# Workflow optimization suggestions
echo "🚀 Workflow Optimizations:"
echo "1. Enable branch protection rules"
echo "2. Set up automated code review"
echo "3. Configure automatic PR labeling"
echo "4. Enable dependabot for dependency updates"
echo "5. Set up automated release notes"
echo ""

# Development velocity metrics
echo "📈 Development Velocity Metrics:"
echo "================================"

# Calculate velocity
DAYS=30
COMMITS_PER_DAY=$(echo "scale=2; $TOTAL_COMMITS / $DAYS" | bc)
echo "Commits per day: $COMMITS_PER_DAY"

# Files changed per commit
AVG_FILES=$(git log --since="30 days ago" --name-only --pretty=format: | grep -v '^$' | wc -l)
if [ "$TOTAL_COMMITS" -gt 0 ]; then
    FILES_PER_COMMIT=$(echo "scale=2; $AVG_FILES / $TOTAL_COMMITS" | bc)
    echo "Average files per commit: $FILES_PER_COMMIT"
fi
echo ""

# Code churn analysis
echo "🔄 Code Churn (last 30 days):"
git log --since="30 days ago" --numstat --pretty=format:'%h' | \
    awk 'NF==3 {plus+=$1; minus+=$2} END {printf "Lines added: +%d\nLines removed: -%d\nNet change: %d\n", plus, minus, plus-minus}'
echo ""

# Summary
echo "📊 Summary Report"
echo "================="
echo "✅ Analysis complete!"
echo ""
echo "Next actions:"
echo "1. Review the metrics above"
echo "2. Address any warnings or recommendations"
echo "3. Share this report with your team"
echo "4. Schedule regular workflow reviews"
echo ""
echo "For more insights, consider:"
echo "- Setting up GitHub Insights"
echo "- Using code quality tools"
echo "- Implementing automated metrics collection"
