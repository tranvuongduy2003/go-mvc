# GitHub Agent Setup Guide

## Quick Start

### 1. Build the Agent
```bash
cd mcp
npm install
npm run build
```

### 2. Create GitHub Personal Access Token

1. Go to https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Select scopes:
   - ✅ `repo` (Full control of private repositories)
   - ✅ `workflow` (Update GitHub Action workflows)
4. Generate and copy the token

### 3. Configure MCP Client

Add to your MCP configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "github-agent": {
      "command": "node",
      "args": ["/absolute/path/to/go-mvc/mcp/dist/agents/github/github-agent.server.js"],
      "env": {
        "GITHUB_TOKEN": "ghp_your_token_here",
        "GITHUB_OWNER": "tranvuongduy2003",
        "GITHUB_REPO": "go-mvc"
      }
    }
  }
}
```

**Important:** Replace `/absolute/path/to/go-mvc` with your actual project path.

### 4. Restart MCP Client

Restart Claude Desktop or your MCP client to load the new agent.

## Usage Examples

### Get Repository Info
```
Use github_get_repo_info
```

### Create an Issue
```
Use github_create_issue with title: "Fix: Database connection timeout"
and body: "The application fails to connect to database after 30 seconds"
and labels: ["bug", "database"]
```

### List Open Issues
```
Use github_list_issues with state: "open" and limit: 10
```

### Create a Pull Request
```
Use github_create_pr with title: "Feature: Add user authentication"
and body: "Implements JWT-based authentication"
and head: "feature/auth" and base: "master"
```

### Trigger CI Workflow
```
Use github_trigger_workflow with workflow_id: "ci.yml" and ref: "master"
```

### View Workflow Runs
```
Use github_list_workflow_runs with workflow_id: "ci.yml" 
and status: "completed" and limit: 5
```

### Create a New Branch
```
Use github_create_branch with branch: "feature/new-feature" 
and from_branch: "master"
```

### Search for Code
```
Use github_search_code with query: "TODO" and limit: 20
```

## Available Tools

| Tool | Description |
|------|-------------|
| `github_get_repo_info` | Get repository information |
| `github_create_issue` | Create a new issue |
| `github_list_issues` | List issues with filters |
| `github_create_pr` | Create a pull request |
| `github_list_prs` | List pull requests |
| `github_get_workflows` | List GitHub Actions workflows |
| `github_trigger_workflow` | Trigger a workflow manually |
| `github_list_workflow_runs` | List workflow runs with status |
| `github_create_branch` | Create a new branch |
| `github_search_code` | Search code in repository |

## Common Use Cases

### 1. Automated Issue Creation from Code Analysis
When AI discovers bugs or improvements:
```
Use github_create_issue with title: "Performance: Optimize database queries"
and body: "Found N+1 query in UserService.GetAllUsers()"
and labels: ["performance", "enhancement"]
```

### 2. PR Creation After Feature Implementation
After implementing a feature:
```
Use github_create_branch with branch: "feature/new-api" and from_branch: "master"

[Make changes locally]

Use github_create_pr with title: "Add new API endpoints"
and body: "Adds /api/v1/users endpoints with full CRUD"
and head: "feature/new-api" and base: "master"
```

### 3. CI/CD Automation
Trigger tests before deploying:
```
Use github_trigger_workflow with workflow_id: "ci.yml" and ref: "develop"
```

Check results:
```
Use github_list_workflow_runs with workflow_id: "ci.yml" 
and status: "in_progress"
```

### 4. Code Maintenance
Find TODOs and create issues:
```
Use github_search_code with query: "TODO" and limit: 50
```

Then create tracking issues for important TODOs.

### 5. Repository Health Check
```
Use github_get_repo_info
Use github_list_issues with state: "open"
Use github_list_prs with state: "open"
Use github_list_workflow_runs with status: "failure"
```

## Troubleshooting

### Authentication Failed
- Verify your GitHub token has correct scopes (`repo`, `workflow`)
- Token should start with `ghp_`
- Check token hasn't expired

### Agent Not Found
- Ensure you built the project: `npm run build`
- Verify path in MCP config is absolute
- Check file exists: `ls mcp/dist/agents/github/github-agent.server.js`

### Workflow Trigger Failed
- Workflow must have `workflow_dispatch` trigger
- Check workflow file name (use exact filename like `ci.yml`)
- Verify token has `workflow` scope

### Rate Limiting
- GitHub API has rate limits (5000 requests/hour for authenticated)
- Use filters to reduce API calls
- Consider implementing caching if needed

## Security Best Practices

1. **Never commit tokens**: Keep tokens in MCP config only
2. **Use minimal scopes**: Only enable required permissions
3. **Rotate tokens regularly**: Generate new tokens periodically
4. **Use fine-grained tokens**: Consider GitHub's fine-grained PATs
5. **Monitor token usage**: Check GitHub settings for suspicious activity

## Integration with DevOps Workflow

This agent complements the existing CI/CD workflows:

- **`.github/workflows/ci.yml`**: Trigger with `github_trigger_workflow`
- **`.github/workflows/security.yml`**: Monitor with `github_list_workflow_runs`
- **`.github/workflows/ai-assistant.yml`**: Get results and create issues
- **Issues**: Auto-create from AI code analysis
- **PRs**: Auto-create after feature implementation

## Next Steps

1. ✅ Build and configure the agent
2. ✅ Test basic operations (repo info, list issues)
3. ✅ Create a test issue
4. ✅ Trigger a workflow
5. ✅ Integrate into your AI development workflow

## Support

For issues or questions:
- Check MCP logs in your client
- Verify environment variables are set
- Test GitHub token with curl:
  ```bash
  curl -H "Authorization: Bearer YOUR_TOKEN" \
       https://api.github.com/user
  ```

---

**Agent Version**: 1.0.0  
**Requires**: Node.js 18+, GitHub Personal Access Token  
**MCP SDK**: ^1.0.4
