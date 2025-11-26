/**
 * Database Configuration
 */
export interface DatabaseConfig {
  host: string;
  port: number;
  database: string;
  user: string;
  password: string;
}

/**
 * Migration Parameters
 */
export interface MigrationParams {
  name?: string;
  migrationFile?: string;
  statements?: string[];
  version?: string;
  rollback?: boolean;
  up?: string;
  down?: string;
}

/**
 * Migration result
 */
export interface MigrationResult {
  success: boolean;
  appliedMigrations: string[];
  executionTime: number;
}

/**
 * Database Column information
 */
export interface ColumnInfo {
  name: string;
  type: string;
  nullable: boolean;
  default?: string;
}

/**
 * Database Index information
 */
export interface IndexInfo {
  name: string;
  columns: string[];
  unique: boolean;
}

/**
 * Foreign Key information
 */
export interface ForeignKeyInfo {
  column: string;
  references: string;
  onDelete?: string;
  onUpdate?: string;
}

/**
 * Table information
 */
export interface TableInfo {
  name: string;
  columns: ColumnInfo[];
  indexes: IndexInfo[];
  foreignKeys: ForeignKeyInfo[];
}

/**
 * Schema Analysis result
 */
export interface SchemaAnalysis {
  tables: TableInfo[];
}

/**
 * Query Analysis result
 */
export interface QueryAnalysis {
  executionPlan: any;
  estimatedCost?: number;
  actualTime: number;
  rowsReturned: number;
  suggestions: string[];
  query?: string;
  suggestedIndexes?: string[];
  performanceIssues?: string[];
  optimizedQuery?: string;
}

/**
 * SQL Generation Parameters
 */
export interface SqlGenerationParams {
  operation: string;
  tables: string[];
  conditions?: string[];
  aggregations?: string[];
  joins?: JoinClause[];
}

/**
 * JOIN clause information
 */
export interface JoinClause {
  type: string;
  table: string;
  on: string;
}
