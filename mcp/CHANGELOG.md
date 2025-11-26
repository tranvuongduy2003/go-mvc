# Changelog

All notable changes to the MCP Agents project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-11-26

### 🎉 Major Refactoring - Modular Architecture

#### Added
- **Modular architecture** following best practices
- **Service layer** for Database Agent (schema, migration, query services)
- **Shared utilities module** with type-safe implementations
- **Comprehensive documentation**:
  - ARCHITECTURE.md (technical architecture)
  - QUICK_REFERENCE.md (developer quick start)
  - Updated README.md (user guide)
- **Integration tests** in proper directory structure
- **TypeScript strict type checking** across all modules

#### Changed
- **Restructured project** from flat to modular hierarchy:
  - `src/index.ts` → `src/agents/api/api-agent.server.ts`
  - `src/database-agent.ts` → `src/agents/database/database-agent.server.ts`
- **Separated concerns**:
  - Tool definitions in `*.tools.ts`
  - Business logic in `*.handlers.ts`
  - Server setup in `*-agent.server.ts`
- **Package.json bin entries** updated to new paths
- **Build output structure** now mirrors source structure

#### Improved
- **Code organization**: Clear module boundaries
- **Type safety**: Comprehensive TypeScript interfaces
- **Maintainability**: Easier to extend and modify
- **Testability**: Modular components enable unit testing
- **Documentation**: Complete architecture and usage docs
- **DRY principle**: Shared utilities eliminate duplication

#### Database Agent Tools (Renamed/Updated)
- `db_create_migration` → removed
- `db_analyze_schema` → `db_schema` (simplified)
- `db_suggest_indexes` → removed (integrated into `db_analyze`)
- `db_execute_query` → `db_query` (renamed)
- `db_debug_query` → `db_analyze` (renamed)
- ✨ **New**: `db_connect` - Explicit connection management
- ✨ **New**: `db_migrate` - Execute migrations with tracking

#### Technical Improvements
- Connection pooling with proper resource management
- Migration tracking with `schema_migrations` table
- Better error handling and error messages
- Query analysis with execution plan insights
- Natural language SQL generation

### 📁 New File Structure

```
Before (Flat):
src/
├── index.ts (539 lines)
└── database-agent.ts (700+ lines)

After (Modular):
src/
├── agents/
│   ├── api/
│   │   ├── api.tools.ts
│   │   ├── api.handlers.ts
│   │   ├── api-agent.server.ts
│   │   └── index.ts
│   └── database/
│       ├── database.tools.ts
│       ├── database.handlers.ts
│       ├── database-agent.server.ts
│       ├── services/
│       │   ├── schema.service.ts
│       │   ├── migration.service.ts
│       │   ├── query.service.ts
│       │   └── index.ts
│       └── index.ts
└── shared/
    ├── types/
    │   ├── api.types.ts
    │   ├── database.types.ts
    │   └── index.ts
    ├── utils/
    │   ├── http.utils.ts
    │   ├── validation.utils.ts
    │   ├── database.utils.ts
    │   └── index.ts
    └── index.ts
```

### 🔧 Configuration Changes

**Old:**
```json
{
  "mcp-api": "./dist/index.js",
  "mcp-db": "./dist/database-agent.js"
}
```

**New:**
```json
{
  "mcp-api": "./dist/agents/api/api-agent.server.js",
  "mcp-db": "./dist/agents/database/database-agent.server.js"
}
```

### ✅ Testing

- All 12 tools tested and working
- Integration tests passing
- Build succeeds without errors
- Type checking passes

### 📊 Metrics

- **Files**: 20+ TypeScript files (was 2)
- **Lines of Code**: ~2000 LOC (better organized)
- **Modules**: 3 main modules (agents, shared, tests)
- **Tools**: 12 tools (6 API + 6 Database)
- **Build Time**: <5 seconds
- **Test Time**: <3 seconds

## [0.2.0] - 2024-11-25

### Added
- Database Agent with 6 tools
- PostgreSQL support
- Migration management
- Schema analysis
- Query optimization

## [0.1.0] - 2024-11-24

### Added
- Initial API Testing Agent
- 6 HTTP method tools (GET, POST, PUT, PATCH, DELETE, TEST)
- Basic MCP protocol implementation
- Axios integration for HTTP requests

---

## Migration Guide

### For Users

If you were using the old paths in your MCP configuration:

1. Update your MCP config file paths:
   - Change `dist/index.js` → `dist/agents/api/api-agent.server.js`
   - Change `dist/database-agent.js` → `dist/agents/database/database-agent.server.js`

2. Rebuild the project:
   ```bash
   cd mcp
   npm run build
   ```

3. Restart your MCP client

### For Developers

If you were extending the agents:

1. **Tool Definitions**: Now in `*.tools.ts` files
2. **Handlers**: Now in `*.handlers.ts` files
3. **Services**: Complex logic should go in `services/` directory
4. **Types**: Add shared types to `shared/types/`
5. **Utils**: Add shared utilities to `shared/utils/`

See ARCHITECTURE.md for detailed guidance.

---

## Future Roadmap

### v1.1.0 (Planned)
- [ ] Unit tests for all handlers
- [ ] Unit tests for all services
- [ ] Enhanced error handling
- [ ] Request/response logging
- [ ] Performance metrics

### v1.2.0 (Planned)
- [ ] MySQL/MariaDB support
- [ ] MongoDB support
- [ ] Redis caching support
- [ ] Batch operations

### v2.0.0 (Future)
- [ ] GraphQL agent
- [ ] WebSocket agent
- [ ] Authentication/Authorization middleware
- [ ] Rate limiting
- [ ] Audit logging

---

## Breaking Changes

### v1.0.0

**Database Agent Tools:**
- `db_create_migration` removed (use `db_migrate`)
- `db_analyze_schema` renamed to `db_schema`
- `db_execute_query` renamed to `db_query`
- `db_debug_query` renamed to `db_analyze`
- New tool: `db_connect` (must call before other operations)

**Configuration:**
- Bin paths changed (see Configuration Changes above)

**Type Signatures:**
- Some handler function signatures changed
- Import paths changed due to restructuring

---

## Contributors

- Initial Development & Refactoring: @tranvuongduy2003

## License

MIT License - see LICENSE file for details
