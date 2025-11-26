# 🤖 AI-Assisted DevOps Guide

> Hướng dẫn sử dụng AI trong quy trình DevOps

## 📋 Mục lục

1. [Giới thiệu](#giới-thiệu)
2. [AI Code Review](#ai-code-review)
3. [AI Code Generation](#ai-code-generation)
4. [AI Workflow Optimization](#ai-workflow-optimization)
5. [AI Bot Commands](#ai-bot-commands)
6. [Integration với IDE](#integration-với-ide)

## 🎯 Giới thiệu

AI DevOps Assistant giúp:
- 🔍 Review code tự động
- 🏗️ Generate boilerplate code
- 📊 Phân tích workflow
- 🤝 Tương tác qua comments
- 💡 Đưa ra suggestions

## 🔍 AI Code Review

### Tính năng

AI Code Review tự động phân tích:
- **Code Complexity**: Cyclomatic complexity, function length
- **Security**: Hardcoded secrets, SQL injection risks
- **Quality**: Code smells, anti-patterns
- **Testing**: Missing tests, coverage
- **Documentation**: Missing comments

### Sử dụng

#### Tự động (trong PR)

Khi tạo PR, AI sẽ tự động review:

```bash
git checkout -b feature/new-api
# Make changes
git commit -m "feat: add new API endpoint"
git push origin feature/new-api
gh pr create

# AI Code Review sẽ tự động chạy và comment
```

#### Manual Review

```bash
# Review current branch
chmod +x .github/scripts/ai-code-review.sh
./.github/scripts/ai-code-review.sh

# Review specific branch vs master
./.github/scripts/ai-code-review.sh feature/new-api master

# Output example:
🤖 AI Code Review Assistant
==========================
Branch: feature/new-api
Base: master

📝 Analyzing changed files...
Changed files:
internal/application/services/product_service.go
internal/presentation/http/handlers/product_handler.go

📊 Analyzing code complexity...
Checking: internal/application/services/product_service.go
✅ Complexity within limits

🔍 Checking for common issues...
TODO comments:
// TODO: Add caching for product list
```

#### Trigger qua PR Comment

```bash
# Trong PR, comment:
/ai review

# AI sẽ chạy review và reply với kết quả
```

### Configuration

Customize trong `.github/scripts/ai-code-review.sh`:

```bash
# Thay đổi complexity threshold
COMPLEXITY_THRESHOLD=15  # Default: 10

# Thay đổi max file size
MAX_FILE_SIZE=500  # Default: 300 lines

# Thêm custom checks
check_custom_pattern() {
    git diff $BASE_BRANCH...$BRANCH | grep "YOUR_PATTERN"
}
```

## 🏗️ AI Code Generation

### Tính năng

Generate code theo Clean Architecture:
- Domain models
- Repository interfaces & implementations
- Services
- HTTP handlers
- DTOs & Commands/Queries

### Sử dụng

#### Full CRUD Generation

```bash
# Generate toàn bộ CRUD cho entity
./.github/scripts/ai-code-generator.sh Product full

# Generated structure:
# internal/domain/product/product.go
# internal/domain/repositories/product_repository.go
# internal/infrastructure/persistence/repositories/product_repository.go
# internal/application/services/product_service.go
# internal/presentation/http/handlers/product_handler.go
```

#### Partial Generation

```bash
# Chỉ generate model
./.github/scripts/ai-code-generator.sh Order model

# Chỉ generate repository
./.github/scripts/ai-code-generator.sh Customer repository

# Chỉ generate service
./.github/scripts/ai-code-generator.sh Invoice service

# Chỉ generate handler
./.github/scripts/ai-code-generator.sh Payment handler
```

### Example: Generate Product Entity

```bash
./.github/scripts/ai-code-generator.sh Product full
```

**Generated Model** (`internal/domain/product/product.go`):
```go
package product

import (
    "time"
    "github.com/google/uuid"
)

type Product struct {
    ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
    Name      string     `gorm:"type:varchar(255);not null" json:"name"`
    CreatedAt time.Time  `gorm:"not null" json:"created_at"`
    UpdatedAt time.Time  `gorm:"not null" json:"updated_at"`
}
```

**Generated Repository Interface**:
```go
package repositories

type ProductRepository interface {
    Create(ctx context.Context, entity *product.Product) error
    GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error)
    List(ctx context.Context, params pagination.Params) ([]product.Product, *pagination.Metadata, error)
    Update(ctx context.Context, entity *product.Product) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### Customization

Modify template trong `.github/scripts/ai-code-generator.sh`:

```bash
# Custom model template
generate_model() {
    cat > "$MODEL_FILE" << EOF
package ${ENTITY_LOWER}

// Add your custom fields here
type ${ENTITY_NAME} struct {
    ID        uuid.UUID
    Name      string
    // Custom fields
    Status    string
    Priority  int
}
EOF
}
```

## 📊 AI Workflow Optimization

### Tính năng

Phân tích và đề xuất cải thiện workflow:
- Commit patterns analysis
- PR metrics & velocity
- Branch management
- Test coverage trends
- Build performance
- Code churn analysis

### Sử dụng

```bash
# Run workflow analysis
./.github/scripts/ai-workflow-optimizer.sh

# Output example:
🤖 AI Development Workflow Optimizer
====================================

📊 Analyzing commit patterns...

👥 Commits by author (last 30 days):
     42 Tran Vuong Duy
     15 John Doe
     8 Jane Smith

📝 Commit message analysis:
✅ Conventional commits: 58 / 65 (89.23%)

📁 Most frequently changed files:
     12 internal/application/services/user_service.go
     8 internal/presentation/http/handlers/user_handler.go
     5 internal/domain/user/user.go

💡 AI-Powered Recommendations
=============================
1. Enable branch protection rules
2. Set up automated code review
3. Configure automatic PR labeling
```

### Trigger qua Comment

```bash
# Trong PR, comment:
/ai optimize

# AI sẽ analyze workflow và reply
```

### Metrics Tracked

- **Commit velocity**: Commits per day
- **Code churn**: Lines added/removed
- **PR lifetime**: Average time to merge
- **Test coverage**: Coverage percentage
- **Build times**: CI/CD duration

## 🤝 AI Bot Commands

### Available Commands

Sử dụng trong PR/Issue comments:

```bash
/ai review          # Trigger full code review
/ai optimize        # Analyze workflow and suggest improvements
/ai docs            # Generate/update documentation
/ai help            # Show available commands
```

### Examples

#### Code Review Command

```bash
# Comment in PR:
/ai review

# Bot response:
🤖 AI Code Review
==================
✅ Code quality: Good
⚠️  Suggestions:
1. Add tests for UserService.CreateUser
2. Consider extracting complex logic in OrderHandler
3. Missing documentation for exported functions
```

#### Workflow Optimization Command

```bash
# Comment in PR:
/ai optimize

# Bot response:
🤖 AI Workflow Analysis
======================
📊 Metrics:
- Commits/day: 3.5
- Test coverage: 85%
- PR merge time: 2.3 days

💡 Suggestions:
1. Consider splitting large PRs
2. Add more integration tests
3. Update stale branches
```

### Custom Commands

Add custom commands trong `.github/workflows/ai-assistant.yml`:

```yaml
- name: Parse command
  id: parse
  run: |
    COMMENT="${{ github.event.comment.body }}"
    
    if [[ "$COMMENT" == *"/ai custom"* ]]; then
      echo "command=custom" >> $GITHUB_OUTPUT
    fi

- name: Execute custom command
  if: steps.parse.outputs.command == 'custom'
  run: |
    # Your custom logic here
    echo "Executing custom command..."
```

## 🔌 Integration với IDE

### VS Code

#### Setup GitHub Copilot

```bash
# Install extension
code --install-extension GitHub.copilot

# Configure
# Ctrl+Shift+P > GitHub Copilot: Sign In
```

#### Sử dụng với Project

```json
// .vscode/settings.json
{
  "github.copilot.enable": {
    "*": true,
    "yaml": true,
    "markdown": true,
    "go": true
  },
  "github.copilot.advanced": {
    "debug.overrideEngine": "gpt-4"
  }
}
```

### GoLand / IntelliJ IDEA

#### Setup

```bash
# Install AI Assistant Plugin
# Settings > Plugins > Search "AI Assistant"
```

#### Usage

```go
// Type comment describing what you want
// AI will suggest implementation

// Example:
// Create a function to calculate fibonacci
// Press Tab to accept suggestion
```

### AI-Assisted Development Workflow

```bash
# 1. Use AI for code generation
./.github/scripts/ai-code-generator.sh Product full

# 2. Use Copilot for implementation
# - Open generated files
# - Add business logic with AI assistance

# 3. Use AI for testing
# - Generate test cases with Copilot
# - Run tests: make test

# 4. Use AI for review
./.github/scripts/ai-code-review.sh

# 5. Create PR
gh pr create

# 6. AI will automatically review PR
# 7. Use /ai commands for additional help

# 8. Merge after approval
gh pr merge
```

## 💡 Best Practices

### 0. Documentation Standards

```bash
# ❌ DON'T: Create summary files after implementation
# - IMPLEMENTATION_COMPLETE.md
# - TASK_SUMMARY.md
# - FEATURE_CHECKLIST.md
# - ARCHITECTURE_VISUALIZATION.md

# ✅ DO: Update existing documentation
# - Update README.md for major features
# - Add to CHANGELOG.md for version tracking
# - Update existing docs in docs/ directory
# - Only create new .md files if explicitly requested

# Rule: Focus on code quality, not documentation overhead
```

### 1. Code Generation

```bash
# ✅ DO: Use for boilerplate
./.github/scripts/ai-code-generator.sh Product full

# ✅ DO: Review and customize generated code
# - Add validation rules
# - Implement business logic
# - Add proper error handling

# ❌ DON'T: Use generated code as-is in production
# ❌ DON'T: Generate code without understanding
# ❌ DON'T: Create summary markdown files after generation
```

### 2. AI Review

```bash
# ✅ DO: Use AI review as first pass
/ai review

# ✅ DO: Combine with human review
gh pr review

# ✅ DO: Address AI suggestions
# - Fix high-priority issues
# - Document why you ignore suggestions

# ❌ DON'T: Trust AI blindly
# ❌ DON'T: Skip human review
```

### 3. Workflow Optimization

```bash
# ✅ DO: Run regularly
./.github/scripts/ai-workflow-optimizer.sh

# ✅ DO: Track metrics over time
# - Save reports
# - Compare trends
# - Set improvement goals

# ✅ DO: Act on recommendations
# - Implement suggested changes
# - Measure impact

# ❌ DON'T: Ignore metrics
# ❌ DON'T: Optimize prematurely
```

## 🎓 Advanced Usage

### Custom AI Scripts

Create custom AI scripts:

```bash
#!/bin/bash
# .github/scripts/ai-custom-analyzer.sh

# Your custom analysis logic
analyze_architecture() {
    echo "Analyzing architecture..."
    # Check layer dependencies
    # Verify design patterns
    # Suggest improvements
}

analyze_architecture
```

### Integration với External AI

```yaml
# .github/workflows/external-ai.yml
- name: Use OpenAI API
  env:
    OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
  run: |
    # Call OpenAI API for advanced analysis
    curl -X POST https://api.openai.com/v1/chat/completions \
      -H "Authorization: Bearer $OPENAI_API_KEY" \
      -d '{"model": "gpt-4", "messages": [...]}'
```

### AI Metrics Dashboard

Create dashboard để theo dõi AI metrics:

```bash
# metrics.sh
#!/bin/bash

echo "## AI Metrics Dashboard"
echo "======================"

# AI Review acceptance rate
# AI suggestion implementation rate
# Code generation usage
# Time saved by AI tools
```

## 🔧 Troubleshooting

### AI Scripts Not Running

```bash
# Make scripts executable
chmod +x .github/scripts/*.sh

# Test locally
./.github/scripts/ai-code-review.sh

# Check workflow permissions
# Settings > Actions > General > Workflow permissions
```

### Bot Not Responding

```bash
# Check workflow runs
gh run list --workflow=ai-assistant.yml

# View logs
gh run view <run-id> --log

# Verify comment format
# Must include: /ai <command>
```

### Poor AI Suggestions

```bash
# Provide more context
# - Better commit messages
# - Clear PR descriptions
# - Proper documentation

# Adjust thresholds in scripts
# - Complexity limits
# - File size limits
# - Coverage requirements
```

## 📚 Resources

- [GitHub Copilot Docs](https://docs.github.com/en/copilot)
- [OpenAI Best Practices](https://platform.openai.com/docs/guides/best-practices)
- [AI in DevOps](https://docs.microsoft.com/en-us/azure/devops/ai/)

## 🎉 Conclusion

AI DevOps Assistant giúp:
- ⚡ Tăng tốc development
- 🎯 Cải thiện code quality
- 📊 Insights về workflow
- 🤖 Automation tasks

Kết hợp AI với human expertise để đạt hiệu quả tốt nhất!

---

**Next Steps**:
1. Setup AI tools trong IDE
2. Try code generation
3. Enable AI reviews
4. Monitor metrics
5. Iterate and improve
