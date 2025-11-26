# MCP Agents - API Testing & Database Management

Complete MCP (Model Context Protocol) server với 2 AI agents mạnh mẽ cho REST API testing và Database management. Được thiết kế theo best practices với kiến trúc modular, scalable và maintainable.

## 🎯 Tính năng

### 🚀 API Testing Agent
- **HTTP Methods đầy đủ**: GET, POST, PUT, PATCH, DELETE
- **Cấu hình linh hoạt**: Headers, query params, request body
- **Response chi tiết**: Status, headers, body, response time
- **API Testing**: Validation status code, response time, JSON schema
- **Dễ dàng mở rộng**: Kiến trúc modular

### 🗄️ Database Agent
- **Connection Management**: Connect và test PostgreSQL databases
- **Query Execution**: Execute SQL queries với parameterized support
- **Schema Analysis**: Phân tích chi tiết database structure
- **Migration Management**: Execute và track database migrations
- **Query Optimization**: Analyze query execution plans và suggestions
- **SQL Generation**: Auto-generate complex SQL queries

## 📁 Cấu trúc Project (Best Practices)

```
mcp/
├── src/
│   ├── agents/                    # Agent modules
│   │   ├── api/                  # API Testing Agent
│   │   │   ├── api.tools.ts      # Tool definitions
│   │   │   ├── api.handlers.ts   # Business logic
│   │   │   ├── api-agent.server.ts # Server entry point
│   │   │   └── index.ts          # Module exports
│   │   └── database/             # Database Agent
│   │       ├── database.tools.ts
│   │       ├── database.handlers.ts
│   │       ├── database-agent.server.ts
│   │       ├── services/         # Service layer
│   │       │   ├── schema.service.ts
│   │       │   ├── migration.service.ts
│   │       │   ├── query.service.ts
│   │       │   └── index.ts
│   │       └── index.ts
│   └── shared/                   # Shared utilities
│       ├── types/                # TypeScript types
│       │   ├── api.types.ts
│       │   ├── database.types.ts
│       │   └── index.ts
│       ├── utils/                # Utility functions
│       │   ├── http.utils.ts
│       │   ├── validation.utils.ts
│       │   ├── database.utils.ts
│       │   └── index.ts
│       └── index.ts
├── tests/
│   └── integration/
│       └── test-agents.js        # Integration tests
├── dist/                         # Compiled output
├── package.json
├── tsconfig.json
└── README.md
```

### Architectural Highlights

✅ **Separation of Concerns**: Tools, handlers, and services are separated  
✅ **Type Safety**: Comprehensive TypeScript interfaces  
✅ **Reusability**: Shared utilities and types  
✅ **Testability**: Modular structure enables easy testing  
✅ **Scalability**: Easy to add new tools or agents  
✅ **Maintainability**: Clear module boundaries and exports

## 📦 Cài đặt

```bash
cd mcp
npm install
npm run build
```

## ⚙️ Cấu hình

Thêm vào file cấu hình MCP của bạn:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`  
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`  
**Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "api-tester": {
      "command": "node",
      "args": ["/absolute/path/to/mcp/dist/agents/api/api-agent.server.js"]
    },
    "database-agent": {
      "command": "node",
      "args": ["/absolute/path/to/mcp/dist/agents/database/database-agent.server.js"]
    }
  }
}
```

**⚠️ Lưu ý:** 
- Thay đổi `/absolute/path/to/mcp` bằng đường dẫn tuyệt đối thực tế trên máy của bạn
- Ví dụ: `/Users/yourname/projects/go-mvc/mcp`
- Restart MCP client sau khi cấu hình

### Quick Test Commands

**API Agent:**
```
Use api_get to fetch https://jsonplaceholder.typicode.com/users/1
```

**Database Agent:**
```
Use db_connect to connect to database localhost:5432, 
database: testdb, user: postgres, password: secret
```

## 🔧 API Testing Agent

### Tools có sẵn

#### 1. `api_get`
Gửi GET request đến REST API endpoint.

**Ví dụ:**
```
Use api_get to fetch https://jsonplaceholder.typicode.com/users/1
```

#### 2. `api_post`
Gửi POST request với JSON body.

**Ví dụ:**
```
Use api_post to create a post at https://jsonplaceholder.typicode.com/posts 
with body: {title: "Test", body: "Content", userId: 1}
```

#### 3. `api_put`
Gửi PUT request để update toàn bộ resource.

#### 4. `api_patch`
Gửi PATCH request để update một phần resource.

#### 5. `api_delete`
Gửi DELETE request để xóa resource.

#### 6. `api_test`
Test API với các validations tự động.

**Ví dụ:**
```
Use api_test to test https://jsonplaceholder.typicode.com/users/1, 
expecting status 200, max response time 2000ms, 
and validate response has id (number), name (string), email (string)
```

### Parameters

| Parameter | Type | Mô tả |
|-----------|------|-------|
| `url` | string | Full URL (required) |
| `method` | string | GET, POST, PUT, PATCH, DELETE |
| `headers` | object | HTTP headers |
| `queryParams` | object | Query parameters |
| `body` | object/array/string | Request body |
| `timeout` | number | Timeout (ms, default: 30000) |
| `expectedStatus` | number | Expected status code (test only) |
| `maxResponseTime` | number | Max response time (test only) |
| `jsonSchema` | object | Schema validation (test only) |

## 🗄️ Database Agent

### Tools có sẵn

#### 1. `db_connect`
Connect đến PostgreSQL database và test connection.

**Ví dụ:**
```
Use db_connect to connect to localhost:5432, database: testdb, 
user: postgres, password: secret
```

**Response:**
```json
{
  "success": true,
  "message": "Connected to database successfully",
  "database": "testdb",
  "user": "postgres",
  "version": "PostgreSQL 14.5..."
}
```

#### 2. `db_query`
Execute SQL query với optional parameters.

**Ví dụ:**
```
Use db_query to execute: "SELECT * FROM users WHERE id = $1" 
with params: [1]
```

**Response:**
```json
{
  "success": true,
  "rows": [...],
  "rowCount": 1,
  "executionTime": 5
}
```

#### 3. `db_schema`
Analyze database schema chi tiết.

**Ví dụ:**
```
Use db_schema to analyze the current database
```

**Response:**
- Danh sách tables với columns (name, type, nullable, default)
- Indexes (name, columns, unique)
- Foreign keys (column, references, onDelete, onUpdate)

#### 4. `db_migrate`
Execute database migrations.

**Ví dụ (inline SQL):**
```
Use db_migrate with statements: [
  "CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(255))",
  "CREATE INDEX idx_users_name ON users(name)"
]
```

**Ví dụ (file-based):**
```
Use db_migrate with migrationFile: "./migrations/001_create_users.sql"
and name: "create_users_table"
```

**Rollback:**
```
Use db_migrate with rollback: true and statements: [
  "DROP TABLE users"
]
```

#### 5. `db_analyze`
Analyze query execution plan và performance.

**Ví dụ:**
```
Use db_analyze to analyze query: 
"SELECT * FROM orders JOIN users ON orders.user_id = users.id WHERE status = 'pending'"
```

**Response:**
- Execution plan (EXPLAIN output)
- Estimated cost
- Actual execution time
- Rows returned
- Optimization suggestions

#### 6. `db_generate_sql`
Generate SQL queries from natural language hoặc structured params.

**Natural Language:**
```
Use db_generate_sql with description: 
"Get all users with their orders from last month"
```

**Structured:**
```
Use db_generate_sql with operation: "SELECT", 
tables: ["orders", "users"],
joins: [{type: "INNER", table: "users", on: "orders.user_id = users.id"}],
conditions: ["orders.created_at > NOW() - INTERVAL '1 month'"]
```

**Hỗ trợ:**
- SELECT với JOINs, aggregations
- INSERT, UPDATE, DELETE
- Complex conditions
- Query explanation

### Database Connection Flow

1. **Connect**: Use `db_connect` với credentials
2. **Execute Operations**: Use các tools khác (query, schema, migrate, etc.)
3. Connection được maintain automatically
4. **Reconnect**: Dùng `db_connect` lại để switch databases

## 📊 Response Format

### API Testing Response
```json
{
  "success": true,
  "request": {
    "method": "GET",
    "url": "https://...",
    "headers": {...},
    "queryParams": {...},
    "body": null
  },
  "response": {
    "status": 200,
    "statusText": "OK",
    "headers": {...},
    "body": {...},
    "responseTime": "245ms"
  }
}
```

### Database Schema Analysis
```json
{
  "success": true,
  "schema": {
    "tables": [
      {
        "name": "users",
        "columns": [...],
        "indexes": [...],
        "foreignKeys": [...]
      }
    ]
  },
  "summary": {
    "totalTables": 5,
    "totalColumns": 45,
    "totalIndexes": 12,
    "totalForeignKeys": 8
  }
}
```

### Query Debug Analysis
```json
{
  "success": true,
  "analysis": {
    "query": "SELECT * FROM ...",
    "executionPlan": {...},
    "performanceIssues": [
      "Sequential scan detected - consider adding indexes",
      "High execution time: 1234.56ms"
    ],
    "suggestedIndexes": [
      "CREATE INDEX idx_orders_user_id ON orders(user_id);"
    ],
    "optimizedQuery": "SELECT id, name FROM ..."
  }
}
```

## 💡 Ví dụ sử dụng

### API Testing Examples

#### Simple GET Request
```
Use api_get to fetch https://jsonplaceholder.typicode.com/users/1
```

#### POST with Authentication
```
Use api_post to send data to https://api.example.com/items 
with headers: Authorization=Bearer token123, Content-Type=application/json
and body: {name: "Item 1", quantity: 5}
```

#### Comprehensive API Test
```
Use api_test to thoroughly test POST to https://api.example.com/users
with body: {name: "John", email: "john@example.com"}
expecting status 201, max response time 3000ms,
and validate response has id (number), name (string), email (string)
```

### Database Examples

#### Create Migration
```
Use db_create_migration named "add_status_to_orders"
with up SQL: "ALTER TABLE orders ADD COLUMN status VARCHAR(50) DEFAULT 'pending';"
and down SQL: "ALTER TABLE orders DROP COLUMN status;"
```

#### Analyze Schema
```
Use db_analyze_schema on localhost:5432/ecommerce_db
user: postgres, password: secret
```

#### Suggest Indexes
```
Use db_suggest_indexes for localhost:5432/ecommerce_db
with query pattern: "SELECT * FROM orders WHERE user_id = ? AND created_at > ?"
```

#### Generate Complex SQL
```
Use db_generate_sql to create SELECT query
from tables: ["orders", "users", "products"]
with joins: [
  {type: "INNER", table: "users", on: "orders.user_id = users.id"},
  {type: "INNER", table: "products", on: "orders.product_id = products.id"}
]
and conditions: ["orders.status = 'completed'", "orders.created_at > '2024-01-01'"]
and aggregations: ["COUNT(*) as total_orders", "SUM(orders.amount) as total_revenue"]
```

#### Debug Slow Query
```
Use db_debug_query to analyze:
"SELECT o.*, u.name, p.title 
FROM orders o 
JOIN users u ON o.user_id = u.id 
JOIN products p ON o.product_id = p.id 
WHERE o.status = 'pending'"
on localhost:5432/ecommerce_db
```

## 🔍 Use Cases

### API Testing
- Test endpoints trong development
- Validate API integrations
- Performance testing
- Schema validation
- Debug API issues
- Document API behavior

### Database Management
- **Migration Management**: Tạo và quản lý database migrations
- **Schema Analysis**: Hiểu rõ cấu trúc database
- **Performance Optimization**: Tìm và fix slow queries
- **Index Optimization**: Suggest và implement optimal indexes
- **Query Development**: Generate complex SQL queries
- **Database Debugging**: Analyze và optimize query performance

## 🛠️ Development

```bash
# Install dependencies
npm install

# Build
npm run build

# Watch mode (development)
npm run dev

# Test API agent
node dist/index.js

# Test Database agent
node dist/database-agent.js
```

## 📋 Requirements

- **Node.js**: 18.0.0+
- **PostgreSQL**: 12+ (cho Database Agent)
- **MCP Client**: Claude Desktop, VS Code với Cline, hoặc bất kỳ MCP-compatible client nào

## 🔒 Security

- Database credentials được truyền qua tool parameters (không lưu trữ)
- Support tất cả authentication methods cho APIs
- Request validation
- Timeout protection
- Error message sanitization

## 🎓 Best Practices

### API Testing
1. Bắt đầu với GET requests đơn giản
2. Test public APIs trước (JSONPlaceholder, HTTPBin)
3. Luôn validate status codes
4. Include Content-Type headers cho POST/PUT
5. Set reasonable timeouts
6. Use schema validation cho production APIs

### Database Management
1. **Always backup** trước khi run migrations
2. Test migrations trên staging environment trước
3. Review suggested indexes trước khi implement
4. Analyze queries trên production-like data
5. Monitor index usage sau khi create
6. Keep migrations reversible (có down SQL)
7. Use transactions cho data migrations

## 📚 Technical Details

### Dependencies
- `@modelcontextprotocol/sdk`: ^1.0.4 - MCP protocol
- `axios`: ^1.7.9 - HTTP client
- `pg`: ^8.13.1 - PostgreSQL client
- `typescript`: ^5.7.2 - TypeScript compiler

### Architecture
- **Modular design**: 2 independent agents
- **Type-safe**: Full TypeScript implementation
- **Error handling**: Comprehensive error handling
- **Performance**: Efficient database connection pooling
- **Extensible**: Easy to add new tools

## 🚀 Performance

- API requests: Configurable timeout (default 30s)
- Database connections: Connection pooling (max 10)
- Response time tracking
- Query performance analysis
- Efficient resource management

## 📈 Future Enhancements

### API Agent
- Response caching
- Request collections
- GraphQL support
- WebSocket testing
- File upload support

### Database Agent
- Support multiple databases (MySQL, SQLite, MongoDB)
- Advanced migration management
- Database backup/restore
- Data seeding
- Schema comparison
- Migration rollback automation

## 🆘 Troubleshooting

### API Agent
- **Server not found**: Check absolute path trong config
- **Request timeouts**: Tăng timeout parameter
- **Schema validation fails**: Check data types

### Database Agent
- **Connection failed**: Verify database credentials
- **Permission denied**: Check user permissions
- **Schema analysis slow**: Database có nhiều tables
- **Migration conflicts**: Check existing schema

## 📄 License

MIT

## 🤝 Contributing

Contributions welcome! Feel free to submit issues or pull requests.

---

**Version**: 1.0.0  
**Status**: Production Ready ✅  
**Node Required**: 18+  
**Protocol**: MCP (Model Context Protocol)
