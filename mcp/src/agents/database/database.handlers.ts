import { Pool } from "pg";
import {
  DatabaseConfig,
  MigrationParams,
  SqlGenerationParams,
  createDatabasePool,
} from "../../shared/index.js";
import {
  analyzeQuery,
  analyzeSchema,
  executeMigration,
  executeQuery,
  getMigrationHistory,
} from "./services/index.js";

// Global connection pool
let connectionPool: Pool | null = null;

/**
 * Handle database connection
 */
export async function handleDbConnect(config: DatabaseConfig): Promise<any> {
  try {
    // Close existing connection if any
    if (connectionPool) {
      await connectionPool.end();
    }

    // Create new connection pool
    connectionPool = createDatabasePool(config);

    // Test connection
    const result = await connectionPool.query(
      "SELECT version(), current_database(), current_user"
    );
    const info = result.rows[0];

    return {
      success: true,
      message: "Connected to database successfully",
      database: info.current_database,
      user: info.current_user,
      version: info.version,
    };
  } catch (error: any) {
    return {
      success: false,
      error: error.message,
    };
  }
}

/**
 * Handle database query execution
 */
export async function handleDbQuery(
  query: string,
  params?: any[]
): Promise<any> {
  if (!connectionPool) {
    throw new Error("No active database connection. Use db_connect first.");
  }

  try {
    const result = await executeQuery(connectionPool, query, params);
    return {
      success: true,
      ...result,
    };
  } catch (error: any) {
    return {
      success: false,
      error: error.message,
    };
  }
}

/**
 * Handle schema analysis
 */
export async function handleDbSchema(): Promise<any> {
  if (!connectionPool) {
    throw new Error("No active database connection. Use db_connect first.");
  }

  try {
    const schema = await analyzeSchema(connectionPool);
    return {
      success: true,
      schema,
    };
  } catch (error: any) {
    return {
      success: false,
      error: error.message,
    };
  }
}

/**
 * Handle database migration
 */
export async function handleDbMigrate(params: MigrationParams): Promise<any> {
  if (!connectionPool) {
    throw new Error("No active database connection. Use db_connect first.");
  }

  try {
    const result = await executeMigration(connectionPool, params);
    const history = await getMigrationHistory(connectionPool);

    return {
      ...result,
      migrationHistory: history,
    };
  } catch (error: any) {
    return {
      success: false,
      error: error.message,
    };
  }
}

/**
 * Handle query analysis
 */
export async function handleDbAnalyze(query: string): Promise<any> {
  if (!connectionPool) {
    throw new Error("No active database connection. Use db_connect first.");
  }

  try {
    const analysis = await analyzeQuery(connectionPool, query);
    return {
      success: true,
      analysis,
    };
  } catch (error: any) {
    return {
      success: false,
      error: error.message,
    };
  }
}

/**
 * Handle SQL generation
 */
export async function handleDbGenerateSql(
  params: SqlGenerationParams & { description?: string }
): Promise<any> {
  try {
    let sql = "";

    // Handle natural language description
    if (params.description) {
      // Simple NLP-based SQL generation
      const desc = params.description.toLowerCase();

      if (
        desc.includes("get") ||
        desc.includes("select") ||
        desc.includes("show")
      ) {
        sql = generateSelectQuery(params);
      } else if (desc.includes("create") || desc.includes("insert")) {
        sql = generateInsertQuery(params);
      } else if (desc.includes("update") || desc.includes("modify")) {
        sql = generateUpdateQuery(params);
      } else if (desc.includes("delete") || desc.includes("remove")) {
        sql = generateDeleteQuery(params);
      } else {
        sql = generateSelectQuery(params); // Default to SELECT
      }
    } else {
      // Generate based on operation type
      switch (params.operation?.toUpperCase()) {
        case "SELECT":
          sql = generateSelectQuery(params);
          break;
        case "INSERT":
          sql = generateInsertQuery(params);
          break;
        case "UPDATE":
          sql = generateUpdateQuery(params);
          break;
        case "DELETE":
          sql = generateDeleteQuery(params);
          break;
        default:
          throw new Error("Invalid operation type");
      }
    }

    return {
      success: true,
      sql,
      explanation: generateExplanation(params),
    };
  } catch (error: any) {
    return {
      success: false,
      error: error.message,
    };
  }
}

/**
 * Generate SELECT query
 */
function generateSelectQuery(params: SqlGenerationParams): string {
  const { tables, conditions, aggregations, joins } = params;

  let sql = "SELECT ";

  // Handle aggregations
  if (aggregations && aggregations.length > 0) {
    sql += aggregations.join(", ");
  } else {
    sql += "*";
  }

  // FROM clause
  sql += `\nFROM ${tables[0]}`;

  // Handle JOINs
  if (joins && joins.length > 0) {
    for (const join of joins) {
      sql += `\n${join.type} JOIN ${join.table} ON ${join.on}`;
    }
  }

  // WHERE clause
  if (conditions && conditions.length > 0) {
    sql += `\nWHERE ${conditions.join(" AND ")}`;
  }

  sql += ";";

  return sql;
}

/**
 * Generate INSERT query
 */
function generateInsertQuery(params: SqlGenerationParams): string {
  const { tables } = params;
  return `INSERT INTO ${tables[0]} (column1, column2, ...)\nVALUES ($1, $2, ...);`;
}

/**
 * Generate UPDATE query
 */
function generateUpdateQuery(params: SqlGenerationParams): string {
  const { tables, conditions } = params;
  let sql = `UPDATE ${tables[0]}\nSET column1 = $1, column2 = $2`;

  if (conditions && conditions.length > 0) {
    sql += `\nWHERE ${conditions.join(" AND ")}`;
  }

  sql += ";";
  return sql;
}

/**
 * Generate DELETE query
 */
function generateDeleteQuery(params: SqlGenerationParams): string {
  const { tables, conditions } = params;
  let sql = `DELETE FROM ${tables[0]}`;

  if (conditions && conditions.length > 0) {
    sql += `\nWHERE ${conditions.join(" AND ")}`;
  }

  sql += ";";
  return sql;
}

/**
 * Generate query explanation
 */
function generateExplanation(params: SqlGenerationParams): string {
  const { operation, tables, conditions, joins, aggregations } = params;

  let explanation = `This query performs a ${
    operation || "SELECT"
  } operation on ${tables.join(", ")}.`;

  if (joins && joins.length > 0) {
    explanation += ` It joins ${joins.length} table(s).`;
  }

  if (conditions && conditions.length > 0) {
    explanation += ` It filters results based on ${conditions.length} condition(s).`;
  }

  if (aggregations && aggregations.length > 0) {
    explanation += ` It uses ${aggregations.length} aggregation(s).`;
  }

  return explanation;
}

/**
 * Close database connection
 */
export async function closeConnection(): Promise<void> {
  if (connectionPool) {
    await connectionPool.end();
    connectionPool = null;
  }
}
