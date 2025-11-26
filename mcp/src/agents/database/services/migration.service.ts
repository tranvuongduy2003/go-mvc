import fs from "fs/promises";
import path from "path";
import { Pool } from "pg";
import { MigrationParams, MigrationResult } from "../../../shared/index.js";

/**
 * Execute database migration
 */
export async function executeMigration(
  pool: Pool,
  params: MigrationParams
): Promise<MigrationResult> {
  const startTime = Date.now();
  const results: string[] = [];

  try {
    // Create migrations table if not exists
    await pool.query(`
      CREATE TABLE IF NOT EXISTS schema_migrations (
        id SERIAL PRIMARY KEY,
        version VARCHAR(255) NOT NULL UNIQUE,
        name VARCHAR(255) NOT NULL,
        executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
      );
    `);

    // Handle file-based migration
    if (params.migrationFile) {
      const filePath = path.resolve(params.migrationFile);
      const sql = await fs.readFile(filePath, "utf8");

      await pool.query("BEGIN");
      await pool.query(sql);

      // Record migration
      const version = path.basename(params.migrationFile, ".sql");
      await pool.query(
        "INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
        [version, params.name || version]
      );

      await pool.query("COMMIT");
      results.push(`Executed migration from file: ${params.migrationFile}`);
    }

    // Handle SQL statements
    if (params.statements && params.statements.length > 0) {
      await pool.query("BEGIN");

      for (const statement of params.statements) {
        await pool.query(statement);
        results.push(`Executed: ${statement.substring(0, 50)}...`);
      }

      // Record migration
      const version = params.version || Date.now().toString();
      await pool.query(
        "INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
        [version, params.name || "manual_migration"]
      );

      await pool.query("COMMIT");
    }

    // Handle rollback
    if (params.rollback) {
      await pool.query("BEGIN");

      // Get last migration
      const lastMigration = await pool.query(
        "SELECT version, name FROM schema_migrations ORDER BY executed_at DESC LIMIT 1"
      );

      if (lastMigration.rows.length > 0) {
        // Execute rollback statements
        if (params.statements) {
          for (const statement of params.statements) {
            await pool.query(statement);
            results.push(`Rolled back: ${statement.substring(0, 50)}...`);
          }
        }

        // Remove migration record
        await pool.query("DELETE FROM schema_migrations WHERE version = $1", [
          lastMigration.rows[0].version,
        ]);

        await pool.query("COMMIT");
        results.push(`Rolled back migration: ${lastMigration.rows[0].version}`);
      } else {
        results.push("No migrations to rollback");
      }
    }

    const executionTime = Date.now() - startTime;

    return {
      success: true,
      appliedMigrations: results,
      executionTime,
    };
  } catch (error) {
    await pool.query("ROLLBACK");
    throw error;
  }
}

/**
 * Get migration history
 */
export async function getMigrationHistory(pool: Pool): Promise<any[]> {
  const result = await pool.query(
    "SELECT * FROM schema_migrations ORDER BY executed_at DESC"
  );
  return result.rows;
}
