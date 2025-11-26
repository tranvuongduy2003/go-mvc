import { Pool } from "pg";
import {
  ColumnInfo,
  ForeignKeyInfo,
  IndexInfo,
  SchemaAnalysis,
  TableInfo,
} from "../../../shared/index.js";

/**
 * Analyze database schema
 */
export async function analyzeSchema(pool: Pool): Promise<SchemaAnalysis> {
  const result: SchemaAnalysis = { tables: [] };

  // Get all tables
  const tablesQuery = `
    SELECT table_name 
    FROM information_schema.tables 
    WHERE table_schema = 'public' 
    ORDER BY table_name;
  `;
  const tablesResult = await pool.query(tablesQuery);

  for (const row of tablesResult.rows) {
    const tableName = row.table_name;

    // Get columns
    const columnsQuery = `
      SELECT 
        column_name,
        data_type,
        is_nullable,
        column_default
      FROM information_schema.columns
      WHERE table_schema = 'public' AND table_name = $1
      ORDER BY ordinal_position;
    `;
    const columnsResult = await pool.query(columnsQuery, [tableName]);

    // Get indexes
    const indexesQuery = `
      SELECT
        i.relname as index_name,
        array_agg(a.attname ORDER BY array_position(ix.indkey, a.attnum)) as columns,
        ix.indisunique as is_unique
      FROM pg_class t
      JOIN pg_index ix ON t.oid = ix.indrelid
      JOIN pg_class i ON i.oid = ix.indexrelid
      JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = ANY(ix.indkey)
      WHERE t.relname = $1
      GROUP BY i.relname, ix.indisunique;
    `;
    const indexesResult = await pool.query(indexesQuery, [tableName]);

    // Get foreign keys
    const fkQuery = `
      SELECT
        kcu.column_name,
        ccu.table_name || '.' || ccu.column_name as references,
        rc.delete_rule,
        rc.update_rule
      FROM information_schema.table_constraints AS tc
      JOIN information_schema.key_column_usage AS kcu
        ON tc.constraint_name = kcu.constraint_name
      JOIN information_schema.constraint_column_usage AS ccu
        ON ccu.constraint_name = tc.constraint_name
      JOIN information_schema.referential_constraints AS rc
        ON tc.constraint_name = rc.constraint_name
      WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name = $1;
    `;
    const fkResult = await pool.query(fkQuery, [tableName]);

    const tableInfo: TableInfo = {
      name: tableName,
      columns: columnsResult.rows.map(
        (col): ColumnInfo => ({
          name: col.column_name,
          type: col.data_type,
          nullable: col.is_nullable === "YES",
          default: col.column_default,
        })
      ),
      indexes: indexesResult.rows.map(
        (idx): IndexInfo => ({
          name: idx.index_name,
          columns: idx.columns,
          unique: idx.is_unique,
        })
      ),
      foreignKeys: fkResult.rows.map(
        (fk): ForeignKeyInfo => ({
          column: fk.column_name,
          references: fk.references,
          onDelete: fk.delete_rule,
          onUpdate: fk.update_rule,
        })
      ),
    };

    result.tables.push(tableInfo);
  }

  return result;
}
