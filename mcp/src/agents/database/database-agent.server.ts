#!/usr/bin/env node

/**
 * Database Agent MCP Server
 * Provides database management tools via Model Context Protocol
 */

import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";

import {
  closeConnection,
  handleDbAnalyze,
  handleDbConnect,
  handleDbGenerateSql,
  handleDbMigrate,
  handleDbQuery,
  handleDbSchema,
} from "./database.handlers.js";
import { databaseTools } from "./database.tools.js";

/**
 * Create and configure the MCP server
 */
const server = new Server(
  {
    name: "database-agent",
    version: "1.0.0",
  },
  {
    capabilities: {
      tools: {},
    },
  }
);

/**
 * List available database tools
 */
server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: databaseTools,
  };
});

/**
 * Handle tool execution requests
 */
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;

  try {
    switch (name) {
      case "db_connect": {
        const {
          host = "localhost",
          port = 5432,
          database,
          user,
          password,
        } = args as any;
        return {
          content: [
            {
              type: "text",
              text: JSON.stringify(
                await handleDbConnect({ host, port, database, user, password }),
                null,
                2
              ),
            },
          ],
        };
      }

      case "db_query": {
        const { query, params } = args as any;
        return {
          content: [
            {
              type: "text",
              text: JSON.stringify(await handleDbQuery(query, params), null, 2),
            },
          ],
        };
      }

      case "db_schema": {
        return {
          content: [
            {
              type: "text",
              text: JSON.stringify(await handleDbSchema(), null, 2),
            },
          ],
        };
      }

      case "db_migrate": {
        const { name, migrationFile, statements, version, rollback } =
          args as any;
        return {
          content: [
            {
              type: "text",
              text: JSON.stringify(
                await handleDbMigrate({
                  name,
                  migrationFile,
                  statements,
                  version,
                  rollback,
                }),
                null,
                2
              ),
            },
          ],
        };
      }

      case "db_analyze": {
        const { query } = args as any;
        return {
          content: [
            {
              type: "text",
              text: JSON.stringify(await handleDbAnalyze(query), null, 2),
            },
          ],
        };
      }

      case "db_generate_sql": {
        const params = args as any;
        return {
          content: [
            {
              type: "text",
              text: JSON.stringify(await handleDbGenerateSql(params), null, 2),
            },
          ],
        };
      }

      default:
        throw new Error(`Unknown tool: ${name}`);
    }
  } catch (error: any) {
    return {
      content: [
        {
          type: "text",
          text: JSON.stringify(
            {
              success: false,
              error: error.message,
            },
            null,
            2
          ),
        },
      ],
      isError: true,
    };
  }
});

/**
 * Start the server
 */
async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);

  console.error("Database Agent MCP Server running on stdio");

  // Handle graceful shutdown
  process.on("SIGINT", async () => {
    await closeConnection();
    await server.close();
    process.exit(0);
  });
}

main().catch((error) => {
  console.error("Fatal error in main():", error);
  process.exit(1);
});
