/**
 * Database Agent Tool Definitions
 * Defines all database-related MCP tools
 */

export const databaseTools = [
  {
    name: "db_connect",
    description:
      "Connect to a PostgreSQL database and test the connection. Returns connection status and database information.",
    inputSchema: {
      type: "object",
      properties: {
        host: {
          type: "string",
          description: "Database host (e.g., localhost)",
          default: "localhost",
        },
        port: {
          type: "number",
          description: "Database port",
          default: 5432,
        },
        database: {
          type: "string",
          description: "Database name",
        },
        user: {
          type: "string",
          description: "Database user",
        },
        password: {
          type: "string",
          description: "Database password",
        },
      },
      required: ["database", "user", "password"],
    },
  },
  {
    name: "db_query",
    description:
      "Execute a SQL query on the connected database. Returns query results with execution time and row count.",
    inputSchema: {
      type: "object",
      properties: {
        query: {
          type: "string",
          description: "SQL query to execute",
        },
        params: {
          type: "array",
          description: "Query parameters for parameterized queries",
          items: {
            type: "string",
          },
        },
      },
      required: ["query"],
    },
  },
  {
    name: "db_schema",
    description:
      "Analyze and return the database schema including tables, columns, indexes, and foreign keys. Provides a comprehensive overview of the database structure.",
    inputSchema: {
      type: "object",
      properties: {},
    },
  },
  {
    name: "db_migrate",
    description:
      "Execute database migrations including creating tables, altering schema, or rolling back changes. Supports both file-based and inline SQL migrations.",
    inputSchema: {
      type: "object",
      properties: {
        name: {
          type: "string",
          description: "Migration name/description",
        },
        migrationFile: {
          type: "string",
          description: "Path to SQL migration file",
        },
        statements: {
          type: "array",
          description: "Array of SQL statements to execute",
          items: {
            type: "string",
          },
        },
        version: {
          type: "string",
          description: "Migration version identifier",
        },
        rollback: {
          type: "boolean",
          description: "Whether to rollback the last migration",
          default: false,
        },
      },
    },
  },
  {
    name: "db_analyze",
    description:
      "Analyze a SQL query's execution plan and performance. Returns optimization suggestions, estimated costs, and performance issues.",
    inputSchema: {
      type: "object",
      properties: {
        query: {
          type: "string",
          description: "SQL query to analyze",
        },
      },
      required: ["query"],
    },
  },
  {
    name: "db_generate_sql",
    description:
      "Generate SQL queries based on natural language description or structured parameters. Helps create complex queries including JOINs, aggregations, and conditions.",
    inputSchema: {
      type: "object",
      properties: {
        description: {
          type: "string",
          description:
            "Natural language description of the desired query (e.g., 'Get all users with their orders from last month')",
        },
        operation: {
          type: "string",
          description: "SQL operation type",
          enum: ["SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "ALTER"],
        },
        tables: {
          type: "array",
          description: "Tables involved in the query",
          items: {
            type: "string",
          },
        },
        conditions: {
          type: "array",
          description: "WHERE conditions",
          items: {
            type: "string",
          },
        },
        aggregations: {
          type: "array",
          description: "Aggregation functions (COUNT, SUM, AVG, etc.)",
          items: {
            type: "string",
          },
        },
        joins: {
          type: "array",
          description: "JOIN clauses",
          items: {
            type: "object",
            properties: {
              type: {
                type: "string",
                enum: ["INNER", "LEFT", "RIGHT", "FULL"],
              },
              table: {
                type: "string",
              },
              on: {
                type: "string",
              },
            },
          },
        },
      },
    },
  },
];
