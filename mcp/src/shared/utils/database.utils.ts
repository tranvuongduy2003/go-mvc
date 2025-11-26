import { Pool } from "pg";
import { DatabaseConfig } from "../types/index.js";

/**
 * Create a PostgreSQL connection pool
 * @param config - Database configuration
 * @returns PostgreSQL Pool instance
 */
export function createDatabasePool(config: DatabaseConfig): Pool {
  return new Pool({
    host: config.host,
    port: config.port,
    database: config.database,
    user: config.user,
    password: config.password,
    max: 10,
    idleTimeoutMillis: 30000,
    connectionTimeoutMillis: 10000,
  });
}
