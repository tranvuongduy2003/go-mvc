# GitHub Actions CI/CD Documentation

> Tài liệu hướng dẫn sử dụng CI/CD Pipeline với GitHub Actions và AI IDE

## 📋 Mục lục

1. [Tổng quan](#tổng-quan)
2. [Workflows](#workflows)
3. [AI Development Scripts](#ai-development-scripts)
4. [Cấu hình](#cấu-hình)
5. [Sử dụng](#sử-dụng)
6. [Best Practices](#best-practices)

## 🎯 Tổng quan

Hệ thống CI/CD được thiết kế với các mục tiêu:

- ✅ **Automation**: Tự động hóa build, test, deploy
- 🔒 **Security**: Quét bảo mật tự động
- 📊 **Quality**: Kiểm tra chất lượng code
- 🤖 **AI-Assisted**: Hỗ trợ AI trong development
- 🚀 **Fast Feedback**: Phản hồi nhanh cho developers

## 📦 Workflows

### 1. CI Pipeline (`ci.yml`)

**Trigger**: Push/PR vào `main`, `master`, `develop`

**Jobs**:
- `lint`: Kiểm tra code quality và formatting
- `test`: Build và run unit tests
- `integration-test`: Run integration tests
- `benchmark`: Performance benchmarks (chỉ chạy trên PR)

**Services**:
- PostgreSQL 15
- Redis 7
- NATS 2.10

**Artifacts**:
- Coverage report
- Binaries
- Benchmark results

**Sử dụng**:
```bash
# Workflow tự động chạy khi push code
git push origin develop

# Kiểm tra status
gh run list --workflow=ci.yml

# Xem chi tiết
gh run view <run-id>
```

### 2. Security & Code Quality (`security.yml`)

**Trigger**: 
- Push/PR vào main branches
- Schedule: Hàng ngày lúc 2 AM UTC

**Jobs**:
- `security-scan`: Gosec, Trivy vulnerability scanning
- `dependency-check`: Kiểm tra dependencies bị lỗi
- `code-analysis`: Static code analysis (staticcheck, gocyclo)
- `codeql`: GitHub CodeQL analysis
- `license-check`: Kiểm tra license compliance
- `secret-scan`: TruffleHog secret scanning
- `docker-scan`: Scan Docker images

**Cấu hình**:
```yaml
# Thêm GitHub Advanced Security (nếu có)
# Settings > Security > Code security and analysis
# Enable: Dependabot, Secret scanning, Code scanning
```

### 3. Docker Build & Push (`docker.yml`)

**Trigger**:
- Push vào main branches
- Tags `v*.*.*`
- PR vào main branches

**Jobs**:
- `build`: Build Docker images cho api, worker, migrate
- `build-dev`: Build development image
- `test-compose`: Test docker-compose configuration
- `build-multiplatform`: Multi-platform build (linux/amd64, linux/arm64)

**Images**:
```bash
# Pull images
docker pull ghcr.io/tranvuongduy2003/go-mvc-api:latest
docker pull ghcr.io/tranvuongduy2003/go-mvc-worker:latest
docker pull ghcr.io/tranvuongduy2003/go-mvc-migrate:latest
```

**Registry**: GitHub Container Registry (ghcr.io)

### 4. Release & Deployment (`release.yml`)

**Trigger**: Tags `v*.*.*`

**Jobs**:
1. `release`: Tạo GitHub Release với changelog
2. `build-binaries`: Build binaries cho multiple platforms
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64)
3. `deploy-staging`: Deploy lên staging (simulation)
4. `deploy-production`: Deploy lên production với approval
5. `update-docs`: Update documentation
6. `verify-deployment`: Post-deployment verification

**Tạo release**:
```bash
# Tạo tag và push
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# GitHub Actions sẽ tự động:
# 1. Tạo release
# 2. Build binaries
# 3. Deploy (với approval)
```

### 5. Dependency Updates (`dependencies.yml`)

**Trigger**: 
- Schedule: Thứ 2 hàng tuần lúc 9 AM UTC
- Manual dispatch

**Jobs**:
- `update-go-deps`: Update Go dependencies
- `update-actions`: Update GitHub Actions versions
- `update-docker`: Update Docker base images
- `security-advisories`: Kiểm tra security advisories
- `comprehensive-check`: Dependency audit report

**Tự động tạo PR** cho updates

### 6. AI Development Assistant (`ai-assistant.yml`)

**Trigger**:
- PR opened/updated
- Comment chứa `/ai` commands

**Jobs**:
- `ai-code-review`: AI-powered code review
- `ai-documentation`: Kiểm tra documentation
- `ai-test-coverage`: Phân tích test coverage
- `ai-performance`: Performance suggestions
- `ai-bot-commands`: Xử lý AI bot commands

**Commands**:
```bash
# Trong PR comment:
/ai review      # Trigger AI code review
/ai optimize    # Analyze workflow
/ai docs        # Generate documentation
/ai help        # Show help
```

## 🤖 AI Development Scripts

### 1. AI Code Review (`ai-code-review.sh`)

**Chức năng**:
- Phân tích code complexity
- Tìm TODO/FIXME comments
- Kiểm tra security issues
- Kiểm tra test coverage
- Đưa ra recommendations

**Sử dụng**:
```bash
# Review current branch
./.github/scripts/ai-code-review.sh

# Review specific branch
./.github/scripts/ai-code-review.sh feature/new-api

# Compare with base branch
./.github/scripts/ai-code-review.sh feature/new-api develop
```

### 2. AI Code Generator (`ai-code-generator.sh`)

**Chức năng**: Generate boilerplate code theo Clean Architecture

**Sử dụng**:
```bash
# Generate full CRUD
./.github/scripts/ai-code-generator.sh Product full

# Generate specific components
./.github/scripts/ai-code-generator.sh Order model
./.github/scripts/ai-code-generator.sh Customer repository
./.github/scripts/ai-code-generator.sh Invoice service
./.github/scripts/ai-code-generator.sh Payment handler
```

**Generated files**:
- Domain model (`internal/domain/{entity}`)
- Repository interface (`internal/domain/repositories`)
- Repository implementation (`internal/infrastructure/persistence`)
- Service (`internal/application/services`)
- HTTP handler (`internal/presentation/http/handlers`)

### 3. AI Workflow Optimizer (`ai-workflow-optimizer.sh`)

**Chức năng**:
- Phân tích git workflow patterns
- Kiểm tra commit conventions
- PR metrics
- Test coverage trends
- Code churn analysis
- Development velocity metrics

**Sử dụng**:
```bash
# Analyze workflow
./.github/scripts/ai-workflow-optimizer.sh

# Output includes:
# - Commit patterns
# - PR statistics
# - Code quality metrics
# - Recommendations
```

## ⚙️ Cấu hình

### 1. GitHub Secrets

Cấu hình tại: `Settings > Secrets and variables > Actions`

**Required secrets**:
```bash
# Không cần secrets cho local development
# Chỉ cần khi deploy lên cloud
```

**Optional secrets** (cho production):
```bash
DEPLOY_SSH_KEY          # SSH key cho deployment
SLACK_WEBHOOK_URL       # Slack notifications
DISCORD_WEBHOOK_URL     # Discord notifications
CODECOV_TOKEN          # Codecov integration
```

### 2. Branch Protection

Cấu hình tại: `Settings > Branches`

**Recommended settings**:
```yaml
Protected branches: main, master, develop

Rules:
- ✅ Require pull request reviews (1+ approvals)
- ✅ Require status checks to pass
  - lint
  - test
  - integration-test
- ✅ Require branches to be up to date
- ✅ Require conversation resolution
- ✅ Include administrators
```

### 3. GitHub Environments

Tạo environments: `Settings > Environments`

**Staging**:
```yaml
Name: staging
Deployment branches: develop, main
Required reviewers: None
Environment secrets: {}
```

**Production**:
```yaml
Name: production
Deployment branches: main, master
Required reviewers: [team-leads]
Environment secrets: {}
Wait timer: 0 minutes
```

### 4. Actions Permissions

Cấu hình tại: `Settings > Actions > General`

```yaml
✅ Allow all actions and reusable workflows
✅ Read and write permissions
✅ Allow GitHub Actions to create and approve pull requests
```

## 🚀 Sử dụng

### Local Development

```bash
# 1. Clone repository
git clone https://github.com/tranvuongduy2003/go-mvc.git
cd go-mvc

# 2. Setup development environment
make setup

# 3. Start development
make dev

# 4. Run tests
make test

# 5. Generate code với AI
./.github/scripts/ai-code-generator.sh Product full

# 6. Review code với AI
./.github/scripts/ai-code-review.sh
```

### Workflow Development

```bash
# 1. Create feature branch
git checkout -b feature/new-feature

# 2. Make changes
# ... code changes ...

# 3. Commit with conventional commits
git commit -m "feat: add new feature"

# 4. Push and create PR
git push origin feature/new-feature
gh pr create

# 5. AI Assistant sẽ tự động:
#    - Review code
#    - Check documentation
#    - Analyze test coverage
#    - Suggest improvements

# 6. Trigger AI commands
# Comment trong PR: /ai review

# 7. Merge sau khi được approve
gh pr merge --squash
```

### Release Process

```bash
# 1. Update version
# Update CHANGELOG.md

# 2. Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 3. GitHub Actions tự động:
#    ✅ Build binaries cho multiple platforms
#    ✅ Create GitHub Release
#    ✅ Generate changelog
#    ✅ Deploy lên staging
#    ✅ Đợi approval cho production
#    ✅ Deploy lên production
#    ✅ Verify deployment
#    ✅ Update documentation

# 4. Monitor deployment
gh run list --workflow=release.yml
gh run watch
```

### Docker Usage

```bash
# 1. Pull latest images
docker pull ghcr.io/tranvuongduy2003/go-mvc-api:latest

# 2. Run with docker-compose
docker-compose up -d

# 3. Check logs
docker-compose logs -f api

# 4. Run migrations
docker-compose exec api ./bin/migrate up

# 5. Stop services
docker-compose down
```

## 📝 Best Practices

### 1. Commit Messages

Sử dụng [Conventional Commits](https://www.conventionalcommits.org/):

```bash
feat: add user authentication
fix: resolve database connection issue
docs: update API documentation
style: format code with gofmt
refactor: restructure user service
perf: optimize database queries
test: add integration tests
build: update dependencies
ci: improve workflow performance
chore: update .gitignore
```

### 2. Pull Requests

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [x] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [x] Code follows style guidelines
- [x] Documentation updated
- [x] Tests added/updated
- [x] No breaking changes
```

### 3. Code Review

```bash
# 1. Review tự động với AI
/ai review

# 2. Kiểm tra suggestions
# - Code complexity
# - Security issues
# - Test coverage
# - Documentation

# 3. Resolve comments

# 4. Request review từ team members
gh pr review --approve
```

### 4. Security

```bash
# 1. Không commit secrets
git secrets --scan

# 2. Sử dụng environment variables
export DATABASE_URL="..."

# 3. Kiểm tra dependencies
make security-check

# 4. Update thường xuyên
make update-deps
```

### 5. Testing

```bash
# Unit tests
go test ./...

# Integration tests
go test -tags=integration ./tests/integration/...

# Coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Benchmarks
go test -bench=. -benchmem ./...
```

## 🔧 Troubleshooting

### Workflow Fails

```bash
# 1. Check workflow logs
gh run view <run-id> --log-failed

# 2. Re-run failed jobs
gh run rerun <run-id> --failed

# 3. Debug locally với act
act -l  # List workflows
act pull_request  # Run PR workflow locally
```

### Docker Issues

```bash
# 1. Clean up
docker-compose down -v
docker system prune -af

# 2. Rebuild
docker-compose build --no-cache

# 3. Check logs
docker-compose logs api
```

### Test Failures

```bash
# 1. Run specific test
go test -v -run TestName ./package

# 2. Debug test
go test -v -run TestName -count=1 ./package

# 3. Check race conditions
go test -race ./...
```

## 📊 Monitoring

### GitHub Actions

```bash
# Workflow status
gh run list

# View specific run
gh run view <run-id>

# Watch live
gh run watch
```

### Metrics

Theo dõi metrics tại:
- GitHub Insights
- Actions usage
- Codecov
- Security alerts

## 🎓 Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Testing](https://golang.org/pkg/testing/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Conventional Commits](https://www.conventionalcommits.org/)

## 📞 Support

- Create issue: `gh issue create`
- Discussion: GitHub Discussions
- Documentation: `/docs`

---

**Note**: Tất cả workflows được thiết kế để chạy local, không cần cloud deployment. Simulation được sử dụng cho deployment steps.
