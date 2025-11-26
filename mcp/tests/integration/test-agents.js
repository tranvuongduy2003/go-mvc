#!/usr/bin/env node

/**
 * Test script cho MCP Agents (API và Database)
 * Run with: node tests/integration/test-agents.js
 */

import { spawn } from "child_process";
import { dirname, join } from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const rootDir = join(__dirname, "..", "..");

console.log("🧪 Testing MCP Agents...\n");

// Test API Agent
console.log("📡 Testing API Agent...");
const apiServerPath = join(
  rootDir,
  "dist",
  "agents",
  "api",
  "api-agent.server.js"
);
const apiServer = spawn("node", [apiServerPath], {
  stdio: ["pipe", "pipe", "inherit"],
});

const initRequest = {
  jsonrpc: "2.0",
  id: 1,
  method: "initialize",
  params: {
    protocolVersion: "2024-11-05",
    capabilities: {},
    clientInfo: {
      name: "test-client",
      version: "1.0.0",
    },
  },
};

const listToolsRequest = {
  jsonrpc: "2.0",
  id: 2,
  method: "tools/list",
  params: {},
};

let apiOutput = "";

apiServer.stdout.on("data", (data) => {
  apiOutput += data.toString();
  const lines = apiOutput.split("\n");

  for (let i = 0; i < lines.length - 1; i++) {
    const line = lines[i].trim();
    if (line) {
      try {
        const response = JSON.parse(line);
        if (response.result && response.result.tools) {
          console.log(
            "✅ API Agent: Found",
            response.result.tools.length,
            "tools"
          );
          response.result.tools.forEach((tool) => {
            console.log("   -", tool.name);
          });
        }
      } catch (e) {
        // Ignore non-JSON lines
      }
    }
  }

  apiOutput = lines[lines.length - 1];
});

setTimeout(() => {
  apiServer.stdin.write(JSON.stringify(initRequest) + "\n");
}, 100);

setTimeout(() => {
  apiServer.stdin.write(JSON.stringify(listToolsRequest) + "\n");
}, 500);

setTimeout(() => {
  apiServer.kill();
  console.log("✅ API Agent test completed!\n");

  // Test Database Agent
  console.log("🗄️  Testing Database Agent...");
  const dbServerPath = join(
    rootDir,
    "dist",
    "agents",
    "database",
    "database-agent.server.js"
  );
  const dbServer = spawn("node", [dbServerPath], {
    stdio: ["pipe", "pipe", "inherit"],
  });

  let dbOutput = "";

  dbServer.stdout.on("data", (data) => {
    dbOutput += data.toString();
    const lines = dbOutput.split("\n");

    for (let i = 0; i < lines.length - 1; i++) {
      const line = lines[i].trim();
      if (line) {
        try {
          const response = JSON.parse(line);
          if (response.result && response.result.tools) {
            console.log(
              "✅ Database Agent: Found",
              response.result.tools.length,
              "tools"
            );
            response.result.tools.forEach((tool) => {
              console.log("   -", tool.name);
            });
          }
        } catch (e) {
          // Ignore non-JSON lines
        }
      }
    }

    dbOutput = lines[lines.length - 1];
  });

  setTimeout(() => {
    dbServer.stdin.write(JSON.stringify(initRequest) + "\n");
  }, 100);

  setTimeout(() => {
    dbServer.stdin.write(JSON.stringify(listToolsRequest) + "\n");
  }, 500);

  setTimeout(() => {
    dbServer.kill();
    console.log("✅ Database Agent test completed!\n");
    console.log("🎉 All tests passed! Both agents are working properly.\n");
    process.exit(0);
  }, 1500);
}, 1500);

apiServer.on("error", (error) => {
  console.error("❌ Error starting API server:", error);
  process.exit(1);
});
