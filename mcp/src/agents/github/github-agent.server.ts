#!/usr/bin/env node

/**
 * GitHub MCP Server
 * Provides GitHub API integration for repository management
 */

import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ErrorCode,
  ListToolsRequestSchema,
  McpError,
} from "@modelcontextprotocol/sdk/types.js";

interface GitHubConfig {
  token: string;
  owner: string;
  repo: string;
}

class GitHubServer {
  private server: Server;
  private config: GitHubConfig;
  private baseUrl = "https://api.github.com";

  constructor() {
    this.server = new Server(
      {
        name: "github-agent",
        version: "1.0.0",
      },
      {
        capabilities: {
          tools: {},
        },
      }
    );

    // Get config from environment or Claude Desktop config
    this.config = {
      token: process.env.GITHUB_TOKEN || "",
      owner: process.env.GITHUB_OWNER || "tranvuongduy2003",
      repo: process.env.GITHUB_REPO || "go-mvc",
    };

    this.setupHandlers();
  }

  private setupHandlers(): void {
    // List available tools
    this.server.setRequestHandler(ListToolsRequestSchema, async () => ({
      tools: [
        {
          name: "github_create_issue",
          description: "Create a new issue in the repository",
          inputSchema: {
            type: "object",
            properties: {
              title: {
                type: "string",
                description: "Issue title",
              },
              body: {
                type: "string",
                description: "Issue description (markdown supported)",
              },
              labels: {
                type: "array",
                items: { type: "string" },
                description: "Issue labels",
              },
              assignees: {
                type: "array",
                items: { type: "string" },
                description: "Assignees (GitHub usernames)",
              },
            },
            required: ["title", "body"],
          },
        },
        {
          name: "github_list_issues",
          description: "List issues in the repository",
          inputSchema: {
            type: "object",
            properties: {
              state: {
                type: "string",
                enum: ["open", "closed", "all"],
                description: "Filter by state",
              },
              labels: {
                type: "string",
                description: "Filter by labels (comma-separated)",
              },
              limit: {
                type: "number",
                description: "Number of issues to return",
                default: 10,
              },
            },
          },
        },
        {
          name: "github_create_pr",
          description: "Create a new pull request",
          inputSchema: {
            type: "object",
            properties: {
              title: {
                type: "string",
                description: "PR title",
              },
              body: {
                type: "string",
                description: "PR description",
              },
              head: {
                type: "string",
                description: "Source branch name",
              },
              base: {
                type: "string",
                description: "Target branch name",
                default: "master",
              },
              draft: {
                type: "boolean",
                description: "Create as draft PR",
                default: false,
              },
            },
            required: ["title", "body", "head"],
          },
        },
        {
          name: "github_list_prs",
          description: "List pull requests in the repository",
          inputSchema: {
            type: "object",
            properties: {
              state: {
                type: "string",
                enum: ["open", "closed", "all"],
                description: "Filter by state",
              },
              limit: {
                type: "number",
                description: "Number of PRs to return",
                default: 10,
              },
            },
          },
        },
        {
          name: "github_get_workflows",
          description: "List GitHub Actions workflows",
          inputSchema: {
            type: "object",
            properties: {},
          },
        },
        {
          name: "github_trigger_workflow",
          description: "Trigger a GitHub Actions workflow",
          inputSchema: {
            type: "object",
            properties: {
              workflow_id: {
                type: "string",
                description: "Workflow file name or ID",
              },
              ref: {
                type: "string",
                description: "Branch or tag to run workflow on",
                default: "master",
              },
              inputs: {
                type: "object",
                description: "Workflow inputs",
              },
            },
            required: ["workflow_id"],
          },
        },
        {
          name: "github_list_workflow_runs",
          description: "List workflow runs",
          inputSchema: {
            type: "object",
            properties: {
              workflow_id: {
                type: "string",
                description: "Filter by workflow file name or ID",
              },
              status: {
                type: "string",
                enum: [
                  "completed",
                  "action_required",
                  "cancelled",
                  "failure",
                  "neutral",
                  "skipped",
                  "stale",
                  "success",
                  "timed_out",
                  "in_progress",
                  "queued",
                  "requested",
                  "waiting",
                ],
                description: "Filter by status",
              },
              limit: {
                type: "number",
                description: "Number of runs to return",
                default: 10,
              },
            },
          },
        },
        {
          name: "github_create_branch",
          description: "Create a new branch",
          inputSchema: {
            type: "object",
            properties: {
              branch: {
                type: "string",
                description: "New branch name",
              },
              from_branch: {
                type: "string",
                description: "Base branch to create from",
                default: "master",
              },
            },
            required: ["branch"],
          },
        },
        {
          name: "github_get_repo_info",
          description: "Get repository information",
          inputSchema: {
            type: "object",
            properties: {},
          },
        },
        {
          name: "github_search_code",
          description: "Search code in the repository",
          inputSchema: {
            type: "object",
            properties: {
              query: {
                type: "string",
                description: "Search query",
              },
              limit: {
                type: "number",
                description: "Number of results to return",
                default: 10,
              },
            },
            required: ["query"],
          },
        },
      ],
    }));

    // Handle tool calls
    this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
      const { name, arguments: args } = request.params;

      try {
        switch (name) {
          case "github_create_issue":
            return await this.createIssue(args);
          case "github_list_issues":
            return await this.listIssues(args);
          case "github_create_pr":
            return await this.createPR(args);
          case "github_list_prs":
            return await this.listPRs(args);
          case "github_get_workflows":
            return await this.getWorkflows();
          case "github_trigger_workflow":
            return await this.triggerWorkflow(args);
          case "github_list_workflow_runs":
            return await this.listWorkflowRuns(args);
          case "github_create_branch":
            return await this.createBranch(args);
          case "github_get_repo_info":
            return await this.getRepoInfo();
          case "github_search_code":
            return await this.searchCode(args);
          default:
            throw new McpError(
              ErrorCode.MethodNotFound,
              `Unknown tool: ${name}`
            );
        }
      } catch (error: any) {
        throw new McpError(
          ErrorCode.InternalError,
          `Error executing ${name}: ${error.message}`
        );
      }
    });
  }

  private async makeRequest(
    method: string,
    endpoint: string,
    data?: any
  ): Promise<any> {
    const url = `${this.baseUrl}${endpoint}`;
    const headers: Record<string, string> = {
      Accept: "application/vnd.github+json",
      "X-GitHub-Api-Version": "2022-11-28",
    };

    if (this.config.token) {
      headers["Authorization"] = `Bearer ${this.config.token}`;
    }

    const response = await fetch(url, {
      method,
      headers,
      body: data ? JSON.stringify(data) : undefined,
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(`GitHub API error: ${response.status} - ${error}`);
    }

    return response.json();
  }

  private async createIssue(args: any) {
    const data = await this.makeRequest(
      "POST",
      `/repos/${this.config.owner}/${this.config.repo}/issues`,
      {
        title: args.title,
        body: args.body,
        labels: args.labels || [],
        assignees: args.assignees || [],
      }
    );

    return {
      content: [
        {
          type: "text",
          text: `✅ Issue created successfully!\n\nTitle: ${data.title}\nNumber: #${data.number}\nURL: ${data.html_url}`,
        },
      ],
    };
  }

  private async listIssues(args: any) {
    const params = new URLSearchParams({
      state: args.state || "open",
      per_page: String(args.limit || 10),
    });

    if (args.labels) {
      params.append("labels", args.labels);
    }

    const data = await this.makeRequest(
      "GET",
      `/repos/${this.config.owner}/${this.config.repo}/issues?${params}`
    );

    const issueList = data
      .map(
        (issue: any) =>
          `#${issue.number} - ${issue.title}\n  State: ${issue.state}\n  URL: ${issue.html_url}`
      )
      .join("\n\n");

    return {
      content: [
        {
          type: "text",
          text: `📋 Issues (${data.length}):\n\n${issueList}`,
        },
      ],
    };
  }

  private async createPR(args: any) {
    const data = await this.makeRequest(
      "POST",
      `/repos/${this.config.owner}/${this.config.repo}/pulls`,
      {
        title: args.title,
        body: args.body,
        head: args.head,
        base: args.base || "master",
        draft: args.draft || false,
      }
    );

    return {
      content: [
        {
          type: "text",
          text: `✅ Pull request created successfully!\n\nTitle: ${data.title}\nNumber: #${data.number}\nURL: ${data.html_url}`,
        },
      ],
    };
  }

  private async listPRs(args: any) {
    const params = new URLSearchParams({
      state: args.state || "open",
      per_page: String(args.limit || 10),
    });

    const data = await this.makeRequest(
      "GET",
      `/repos/${this.config.owner}/${this.config.repo}/pulls?${params}`
    );

    const prList = data
      .map(
        (pr: any) =>
          `#${pr.number} - ${pr.title}\n  State: ${pr.state}\n  From: ${pr.head.ref} → ${pr.base.ref}\n  URL: ${pr.html_url}`
      )
      .join("\n\n");

    return {
      content: [
        {
          type: "text",
          text: `🔀 Pull Requests (${data.length}):\n\n${prList}`,
        },
      ],
    };
  }

  private async getWorkflows() {
    const data = await this.makeRequest(
      "GET",
      `/repos/${this.config.owner}/${this.config.repo}/actions/workflows`
    );

    const workflowList = data.workflows
      .map(
        (wf: any) =>
          `${wf.name}\n  ID: ${wf.id}\n  File: ${wf.path}\n  State: ${wf.state}`
      )
      .join("\n\n");

    return {
      content: [
        {
          type: "text",
          text: `⚙️ Workflows (${data.total_count}):\n\n${workflowList}`,
        },
      ],
    };
  }

  private async triggerWorkflow(args: any) {
    await this.makeRequest(
      "POST",
      `/repos/${this.config.owner}/${this.config.repo}/actions/workflows/${args.workflow_id}/dispatches`,
      {
        ref: args.ref || "master",
        inputs: args.inputs || {},
      }
    );

    return {
      content: [
        {
          type: "text",
          text: `✅ Workflow triggered successfully!\n\nWorkflow: ${
            args.workflow_id
          }\nBranch: ${args.ref || "master"}`,
        },
      ],
    };
  }

  private async listWorkflowRuns(args: any) {
    const params = new URLSearchParams({
      per_page: String(args.limit || 10),
    });

    if (args.status) {
      params.append("status", args.status);
    }

    let endpoint = `/repos/${this.config.owner}/${this.config.repo}/actions/runs?${params}`;

    if (args.workflow_id) {
      endpoint = `/repos/${this.config.owner}/${this.config.repo}/actions/workflows/${args.workflow_id}/runs?${params}`;
    }

    const data = await this.makeRequest("GET", endpoint);

    const runList = data.workflow_runs
      .map(
        (run: any) =>
          `Run #${run.run_number} - ${run.name}\n  Status: ${
            run.status
          }\n  Conclusion: ${run.conclusion || "N/A"}\n  Branch: ${
            run.head_branch
          }\n  URL: ${run.html_url}`
      )
      .join("\n\n");

    return {
      content: [
        {
          type: "text",
          text: `🏃 Workflow Runs (${data.total_count}):\n\n${runList}`,
        },
      ],
    };
  }

  private async createBranch(args: any) {
    // Get the SHA of the base branch
    const refData = await this.makeRequest(
      "GET",
      `/repos/${this.config.owner}/${this.config.repo}/git/ref/heads/${
        args.from_branch || "master"
      }`
    );

    // Create new branch
    await this.makeRequest(
      "POST",
      `/repos/${this.config.owner}/${this.config.repo}/git/refs`,
      {
        ref: `refs/heads/${args.branch}`,
        sha: refData.object.sha,
      }
    );

    return {
      content: [
        {
          type: "text",
          text: `✅ Branch created successfully!\n\nBranch: ${
            args.branch
          }\nFrom: ${args.from_branch || "master"}`,
        },
      ],
    };
  }

  private async getRepoInfo() {
    const data = await this.makeRequest(
      "GET",
      `/repos/${this.config.owner}/${this.config.repo}`
    );

    return {
      content: [
        {
          type: "text",
          text: `📦 Repository: ${data.full_name}\n\nDescription: ${
            data.description || "N/A"
          }\nLanguage: ${data.language}\nStars: ${
            data.stargazers_count
          }\nForks: ${data.forks_count}\nOpen Issues: ${
            data.open_issues_count
          }\nDefault Branch: ${data.default_branch}\nURL: ${data.html_url}`,
        },
      ],
    };
  }

  private async searchCode(args: any) {
    const params = new URLSearchParams({
      q: `${args.query}+repo:${this.config.owner}/${this.config.repo}`,
      per_page: String(args.limit || 10),
    });

    const data = await this.makeRequest("GET", `/search/code?${params}`);

    const resultList = data.items
      .map(
        (item: any) =>
          `${item.name}\n  Path: ${item.path}\n  URL: ${item.html_url}`
      )
      .join("\n\n");

    return {
      content: [
        {
          type: "text",
          text: `🔍 Search Results (${data.total_count}):\n\n${resultList}`,
        },
      ],
    };
  }

  async run(): Promise<void> {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
    console.error("GitHub MCP server running on stdio");
  }
}

const server = new GitHubServer();
server.run().catch(console.error);
