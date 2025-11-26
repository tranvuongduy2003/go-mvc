import { Pool } from "pg";
import { QueryAnalysis } from "../../../shared/index.js";

/**
 * Analyze query execution plan
 */
export async function analyzeQuery(
  pool: Pool,
  query: string
): Promise<QueryAnalysis> {
  const startTime = Date.now();

  // Get execution plan
  const explainResult = await pool.query(`EXPLAIN (FORMAT JSON) ${query}`);
  const plan = explainResult.rows[0]["QUERY PLAN"][0];

  // Execute query to get actual timing
  const executeStart = Date.now();
  const result = await pool.query(query);
  const executionTime = Date.now() - executeStart;

  // Parse plan for suggestions
  const suggestions: string[] = [];

  // Check for sequential scans
  if (JSON.stringify(plan).includes("Seq Scan")) {
    suggestions.push(
      "Consider adding indexes to avoid sequential scans on large tables"
    );
  }

  // Check for nested loops on large datasets
  if (JSON.stringify(plan).includes("Nested Loop")) {
    suggestions.push(
      "Nested loops detected - consider using hash joins for large datasets"
    );
  }

  // Check for sorts
  if (JSON.stringify(plan).includes("Sort")) {
    suggestions.push(
      "Sort operation detected - consider adding indexes on ORDER BY columns"
    );
  }

  const totalTime = Date.now() - startTime;

  return {
    executionPlan: plan,
    estimatedCost: plan["Plan"]["Total Cost"],
    actualTime: executionTime,
    rowsReturned: result.rows.length,
    suggestions,
  };
}

/**
 * Execute query with parameters
 */
export async function executeQuery(
  pool: Pool,
  query: string,
  params?: any[]
): Promise<any> {
  const startTime = Date.now();
  const result = await pool.query(query, params);
  const executionTime = Date.now() - startTime;

  return {
    rows: result.rows,
    rowCount: result.rowCount,
    executionTime,
  };
}
