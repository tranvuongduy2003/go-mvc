# MCP Agents - Architecture Documentation

## 📐 Architecture Overview

Dự án MCP Agents được thiết kế theo kiến trúc modular, tuân thủ các best practices về separation of concerns, maintainability và scalability.

## 🏗️ Layered Architecture

```
┌─────────────────────────────────────────┐
│         MCP Protocol Layer              │
│    (api-agent.server.ts,                │
│     database-agent.server.ts)           │
└─────────────────────────────────────────┘
                    ▼
┌─────────────────────────────────────────┐
│         Tool Definition Layer           │
│    (api.tools.ts, database.tools.ts)    │
└─────────────────────────────────────────┘
                    ▼
┌─────────────────────────────────────────┐
│         Business Logic Layer            │
│    (*.handlers.ts, services/)           │
└─────────────────────────────────────────┘
                    ▼
┌─────────────────────────────────────────┐
│         Utility & Type Layer            │
│    (shared/types/, shared/utils/)       │
└─────────────────────────────────────────┘
```

## 📁 Module Structure

### 1. **Agents Module** (`src/agents/`)

Chứa các independent agent modules. Mỗi agent có:

#### API Agent (`src/agents/api/`)
```
api/
├── api.tools.ts           # Tool definitions (MCP protocol)
├── api.handlers.ts        # Request handlers & business logic
├── api-agent.server.ts    # Server entry point
└── index.ts               # Module exports
```

#### Database Agent (`src/agents/database/`)
```
database/
├── database.tools.ts           # Tool definitions
├── database.handlers.ts        # Request handlers
├── database-agent.server.ts    # Server entry point
├── services/                   # Service layer
│   ├── schema.service.ts       # Schema analysis logic
│   ├── migration.service.ts    # Migration management
│   ├── query.service.ts        # Query execution & analysis
│   └── index.ts
└── index.ts
```

**Design Principles:**
- **Single Responsibility**: Mỗi file có một mục đích rõ ràng
- **Service Layer**: Complex business logic được tách ra services
- **Clear Dependencies**: Handlers depend on services, services depend on shared utils

### 2. **Shared Module** (`src/shared/`)

Chứa code được share giữa các agents.

```
shared/
├── types/              # TypeScript interfaces & types
│   ├── api.types.ts    # API-related types
│   ├── database.types.ts
│   └── index.ts
├── utils/              # Utility functions
│   ├── http.utils.ts   # HTTP client wrapper
│   ├── validation.utils.ts
│   ├── database.utils.ts
│   └── index.ts
└── index.ts            # Central export point
```

**Benefits:**
- **DRY Principle**: Không duplicate code
- **Type Safety**: Centralized type definitions
- **Reusability**: Utilities có thể reuse ở nhiều nơi
- **Testability**: Dễ dàng test isolated units

### 3. **Tests Module** (`tests/`)

```
tests/
├── integration/
│   └── test-agents.js     # Integration tests
└── unit/                  # (future) Unit tests
```

## 🔄 Data Flow

### API Agent Request Flow

```
User Request
    ↓
MCP Server (api-agent.server.ts)
    ↓
Tool Router (switch/case)
    ↓
Handler Function (api.handlers.ts)
    ↓
HTTP Utility (http.utils.ts)
    ↓
External API
    ↓
Response ← Response ← Response ← Response
```

### Database Agent Request Flow

```
User Request
    ↓
MCP Server (database-agent.server.ts)
    ↓
Tool Router
    ↓
Handler Function (database.handlers.ts)
    ↓
Service Layer (services/*.service.ts)
    ↓
Database Utility (database.utils.ts)
    ↓
PostgreSQL
    ↓
Response ← Response ← Response ← Response
```

## 🎯 Design Patterns

### 1. **Module Pattern**
- Mỗi agent là một self-contained module
- Clear exports through `index.ts`
- No circular dependencies

### 2. **Service Layer Pattern**
- Complex business logic tách ra services
- Handlers remain thin, focus on routing
- Services can be reused and tested independently

### 3. **Factory Pattern**
- Database connection pooling (`createDatabasePool`)
- Reusable configuration objects

### 4. **Strategy Pattern**
- Different handlers for different tools
- Extensible tool registration

## 🔧 Extension Points

### Adding a New Tool to Existing Agent

1. **Define Tool** in `*.tools.ts`:
```typescript
export const myTools = [
  {
    name: "my_new_tool",
    description: "...",
    inputSchema: {...}
  }
];
```

2. **Create Handler** in `*.handlers.ts`:
```typescript
export async function handleMyNewTool(params: any): Promise<any> {
  // Business logic
}
```

3. **Register in Server** in `*-agent.server.ts`:
```typescript
case "my_new_tool": {
  return await handleMyNewTool(args);
}
```

### Adding a New Agent

1. Create new directory: `src/agents/my-agent/`
2. Create files following pattern:
   - `my-agent.tools.ts`
   - `my-agent.handlers.ts`
   - `my-agent.server.ts`
   - `index.ts`
3. Add shared types/utils if needed
4. Update `package.json` bin entry
5. Add tests

## 📊 Type System

### Type Hierarchy

```
Base Types (shared/types/)
    ├── API Types
    │   ├── ApiRequestParams
    │   ├── ApiResponse
    │   └── TestResult
    └── Database Types
        ├── DatabaseConfig
        ├── SchemaAnalysis
        └── QueryAnalysis
```

### Type Usage

- **Input Validation**: Types ensure correct tool parameters
- **Handler Contracts**: Define clear interfaces between layers
- **Documentation**: Types serve as inline documentation
- **IDE Support**: Full IntelliSense and autocomplete

## 🧪 Testing Strategy

### Current Tests
- **Integration Tests**: Verify MCP protocol communication
- **Tool Discovery**: Ensure all tools are registered correctly

### Future Testing (Recommended)

1. **Unit Tests**
   - Test individual handlers
   - Test service functions
   - Test utility functions

2. **Integration Tests**
   - Test complete request flows
   - Test error handling
   - Test edge cases

3. **E2E Tests**
   - Test with real PostgreSQL
   - Test with real APIs

## 🚀 Performance Considerations

### Connection Pooling
- Database connections use pooling (`pg.Pool`)
- Max 10 connections, 30s idle timeout
- Automatic connection cleanup

### Error Handling
- Try-catch blocks in all handlers
- Graceful degradation
- Meaningful error messages

### Resource Management
- Connection cleanup on shutdown
- Proper event handler cleanup
- No memory leaks

## 📈 Scalability

### Horizontal Scaling
- Stateless servers (except DB connection)
- Can run multiple instances
- Load balancing ready

### Vertical Scaling
- Modular structure allows code splitting
- Services can be extracted to microservices
- Database connection pool can be tuned

## 🔒 Security Considerations

### Current Implementation
- No credentials stored in code
- Parameters passed at runtime
- Connection credentials user-provided

### Recommendations
- Add credential validation
- Implement rate limiting
- Add SQL injection prevention in generated queries
- Audit logging for sensitive operations

## 📝 Maintenance Guidelines

### Code Style
- TypeScript strict mode
- ESLint rules (future)
- Consistent naming conventions
- JSDoc comments for public APIs

### Documentation
- README for user-facing docs
- ARCHITECTURE.md (this file) for developers
- Inline comments for complex logic
- Type definitions as documentation

### Version Control
- Semantic versioning
- Changelog maintenance
- Breaking changes clearly documented

## 🎓 Learning Resources

### Understanding MCP
- [Model Context Protocol Spec](https://modelcontextprotocol.io/)
- [MCP SDK Documentation](https://github.com/modelcontextprotocol/typescript-sdk)

### Best Practices References
- Clean Architecture principles
- SOLID principles
- TypeScript best practices
- Node.js design patterns
