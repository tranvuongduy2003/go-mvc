# User Story Template for AI API Generation

## 📝 How to Use This Template

Copy template bên dưới và điền thông tin chi tiết cho chức năng bạn muốn AI sinh API. Template này đảm bảo AI có đủ thông tin để sinh ra code hoàn chỉnh theo Clean Architecture.

---

## User Story: [Tên chức năng - VD: Create Product, Update User Profile, etc.]

### Business Description
- **Actor**: [Ai sẽ sử dụng - User, Admin, System, Guest, etc.]
- **Action**: [Hành động gì - Create, Update, Delete, Get, List, Search, Upload, etc.]
- **Object**: [Đối tượng gì - Product, Order, User, Profile, Document, etc.]
- **Purpose**: [Mục đích/lợi ích - Tại sao cần chức năng này, giá trị business]

### Functional Requirements

#### Inputs
```
field1: type (validation rules) - [Mô tả field]
field2: type (validation rules) - [Mô tả field]
...

Example:
name: string (required, 2-100 chars) - Product name
price: decimal (required, > 0) - Product price in USD
category_id: uuid (required, must exist) - Reference to category
description: string (optional, max 1000 chars) - Product description
```

#### Outputs
```
Định nghĩa response structure và format

Example:
- Created product with ID, name, price, timestamps
- Response includes generated ID, validation status, creation timestamp
- Return full product details for confirmation
```

#### Business Rules
```
Rule 1: [Mô tả business logic cụ thể]
Rule 2: [Mô tả constraints và limitations]
Rule 3: [Mô tả relationships và dependencies]
...

Example:
- Product name must be unique within same category
- Price must be positive number with 2 decimal places
- Category must exist and be in active status
- Only admin users can create products
- Maximum 10 products can be created per hour per user
```

### Technical Specifications
- **HTTP Method**: [GET/POST/PUT/PATCH/DELETE]
- **Endpoint**: /api/v1/[resource] - [Mô tả endpoint structure]
- **Authentication**: [Required/Optional] - [Chi tiết authentication requirements]
- **Authorization**: [Role-based rules] - [Specific permissions needed]
- **Pagination**: [Yes/No] - [If yes, specify default page size và max limit]
- **Filtering**: [Yes/No] - [If yes, specify filterable fields]
- **Sorting**: [Yes/No] - [If yes, specify sortable fields]

### Database Impact

#### Tables
```
Primary table: [main table name]
Related tables: [list of related tables]

Example:
Primary table: products
Related tables: categories, product_images, product_reviews
```

#### Relationships
```
Describe foreign key relationships và cardinality

Example:
products.category_id -> categories.id (many-to-one)
products.id <- product_images.product_id (one-to-many)
products.id <- product_reviews.product_id (one-to-many)
```

#### Indexes
```
List required indexes for performance

Example:
- idx_products_name_category (name, category_id) - for uniqueness check
- idx_products_category_id (category_id) - for filtering
- idx_products_created_at (created_at) - for sorting
```

#### Migration Requirements
```
Describe database schema changes needed

Example:
- Create products table with all fields
- Add foreign key constraint to categories table
- Create unique constraint on (name, category_id)
- Add indexes for performance
```

### Validation Rules

#### Required Fields
```
List all mandatory fields

Example:
- name (cannot be empty or whitespace only)
- price (must be provided)
- category_id (must be valid UUID)
```

#### Format Validation
```
Specify format requirements for fields

Example:
- email: valid email format with @ và domain
- phone: international format (+country code)
- price: decimal with max 2 decimal places
- url: valid HTTP/HTTPS URL format
```

#### Business Validation
```
Specify business logic validation rules

Example:
- Product name unique within same category (database check)
- Category must exist và be active (database check)
- Price must be within allowed range for category
- User must have permission to create products
```

#### Size Limits
```
Specify size constraints for fields

Example:
- name: 2-100 characters
- description: maximum 1000 characters
- image file: maximum 5MB
- bulk operations: maximum 100 items per request
```

### Error Scenarios

#### Client Errors (4xx)
```
400 Bad Request:
- [Specific validation error scenarios]
- [Invalid JSON format scenarios]
- [Missing required parameters]

401 Unauthorized:
- [Authentication failure scenarios]
- [Invalid or expired tokens]

403 Forbidden:
- [Authorization failure scenarios]
- [Insufficient permissions]

404 Not Found:
- [Resource not found scenarios]
- [Invalid resource IDs]

409 Conflict:
- [Resource conflict scenarios]
- [Uniqueness constraint violations]

422 Unprocessable Entity:
- [Business rule violation scenarios]
- [Complex validation failures]

Example:
400: Invalid JSON format, missing required fields, invalid data types
401: No JWT token provided, expired token, invalid token signature
403: User doesn't have admin role, cannot modify other user's data
404: Product not found, category not found
409: Product name already exists in category
422: Price exceeds category maximum, inactive category selected
```

#### Server Errors (5xx)
```
500 Internal Server Error:
- [Database connection failures]
- [External service failures]
- [Unexpected application errors]

503 Service Unavailable:
- [Temporary service outages]
- [Database maintenance scenarios]

Example:
500: Database connection timeout, external email service failure
503: Database maintenance mode, Redis cache unavailable
```

### Performance Requirements
- **Response Time**: [Expected latency] - [Acceptable response time under normal load]
- **Throughput**: [Expected RPS] - [Number of requests per second to handle]
- **Caching Strategy**: [Cache approach] - [What to cache và TTL]
- **Rate Limiting**: [Requests per period] - [Rate limiting rules]

```
Example:
Response Time: < 500ms for 95% of requests under normal load
Throughput: 100 RPS sustained, 200 RPS peak
Caching: Cache category lookups for 5 minutes, cache product lists for 1 minute
Rate Limiting: 10 product creations per minute per user
```

### Integration Requirements

#### External APIs
```
List third-party services needed

Example:
- Payment gateway for price validation
- Image processing service for product images
- Email service for notifications
- Analytics service for tracking
```

#### Message Queue
```
Specify async processing needs

Example:
- Queue product creation event for search index update
- Queue email notification to admin users
- Queue image processing for product photos
```

#### Notifications
```
Specify notification requirements

Example:
- Email notification to admin when new product created
- SMS notification for high-value products
- Push notification to mobile users about new products
```

### Security Considerations
- **Input Sanitization**: [XSS prevention, SQL injection protection]
- **Data Encryption**: [Fields requiring encryption]
- **Audit Logging**: [What actions to log]
- **Rate Limiting**: [Abuse prevention measures]

```
Example:
Input Sanitization: HTML encode all text inputs, validate file uploads
Data Encryption: Encrypt sensitive product cost data
Audit Logging: Log all product creation, modification, deletion with user context
Rate Limiting: Prevent spam product creation, limit API calls per user
```

### Testing Requirements
- **Unit Tests**: [Domain logic to test]
- **Integration Tests**: [API endpoints to test]
- **Performance Tests**: [Load testing requirements]

```
Example:
Unit Tests: Product name validation, price calculation, business rule enforcement
Integration Tests: Full product creation workflow, error handling scenarios
Performance Tests: 100 concurrent users creating products
```

---

## 📚 Example User Stories

### Example 1: Simple CRUD
```markdown
## User Story: Create Product Category

### Business Description
- **Actor**: Admin User
- **Action**: Create
- **Object**: Product Category
- **Purpose**: Organize products into categories for better catalog management

### Functional Requirements
#### Inputs
name: string (required, 2-50 chars) - Category name
description: string (optional, max 500 chars) - Category description
is_active: boolean (default true) - Category status

#### Outputs
- Created category with ID, name, description, status, timestamps
- Return full category details for confirmation

#### Business Rules
- Category name must be unique across all categories
- Only admin users can create categories
- New categories are active by default

### Technical Specifications
- **HTTP Method**: POST
- **Endpoint**: /api/v1/categories
- **Authentication**: Required (JWT)
- **Authorization**: Admin role required
- **Pagination**: N/A

### Database Impact
#### Tables
Primary table: categories

#### Relationships
categories.id <- products.category_id (one-to-many)

#### Indexes
- idx_categories_name (name) - for uniqueness và search
- idx_categories_is_active (is_active) - for filtering

#### Migration Requirements
- Create categories table
- Add unique constraint on name
- Add indexes for performance

### Validation Rules
#### Required Fields
- name (2-50 characters, no special characters except spaces và hyphens)

#### Business Validation
- Category name must be unique (database check)
- User must have admin role

### Error Scenarios
#### Client Errors (4xx)
400: Invalid name format, name too short/long
401: No authentication token
403: User is not admin
409: Category name already exists

#### Server Errors (5xx)
500: Database connection error

### Performance Requirements
- **Response Time**: < 200ms
- **Throughput**: 50 RPS
- **Caching Strategy**: Cache active categories for 10 minutes
```

### Example 2: Complex Business Logic
```markdown
## User Story: Process Order Payment

### Business Description
- **Actor**: System (automated)
- **Action**: Process
- **Object**: Order Payment
- **Purpose**: Handle payment processing and update order status accordingly

### Functional Requirements
#### Inputs
order_id: uuid (required, must exist) - Order to process payment for
payment_method: string (required, enum: credit_card, paypal, bank_transfer) - Payment method
payment_details: object (required, varies by method) - Payment method specific data
idempotency_key: string (required, unique) - Prevent duplicate processing

#### Outputs
- Payment result with status, transaction ID, updated order status
- Payment receipt details
- Error details if payment fails

#### Business Rules
- Order must exist và be in "pending_payment" status
- Payment amount must match order total exactly
- Idempotency key prevents duplicate payments
- Failed payments are retried up to 3 times
- Successful payments update order status to "paid"
- Inventory is reserved during payment processing

### Technical Specifications
- **HTTP Method**: POST
- **Endpoint**: /api/v1/orders/{order_id}/payments
- **Authentication**: Required (JWT or API key)
- **Authorization**: User must own order or be admin
- **Pagination**: N/A

### Database Impact
#### Tables
Primary table: payments
Related tables: orders, payment_methods, transactions

#### Relationships
payments.order_id -> orders.id (many-to-one)
payments.payment_method_id -> payment_methods.id (many-to-one)
payments.id <- transactions.payment_id (one-to-many)

#### Indexes
- idx_payments_order_id (order_id) - for order lookups
- idx_payments_idempotency_key (idempotency_key) - for duplicate prevention
- idx_payments_status_created_at (status, created_at) - for reporting

### Integration Requirements
#### External APIs
- Stripe API for credit card processing
- PayPal API for PayPal payments
- Bank API for bank transfers

#### Message Queue
- Queue order status update event
- Queue inventory release event on failure
- Queue email receipt sending

#### Notifications
- Email payment confirmation to customer
- SMS notification for high-value payments
- Admin notification for failed payments
```

---

## 🚀 Quick Start Guide

1. **Copy template trên** và save thành file `.md`
2. **Điền thông tin chi tiết** cho chức năng bạn muốn
3. **Đưa User Story cho AI** với instruction: "Hãy sinh API hoàn chỉnh theo Clean Architecture từ User Story này"
4. **AI sẽ tự động sinh**:
   - Domain entities và value objects
   - Repository interfaces và implementations
   - Commands/Queries với handlers
   - DTOs và validators
   - HTTP handlers với Swagger docs
   - Database migrations
   - Dependency injection setup

## 💡 Tips for Better Results

### ✅ DO
- **Chi tiết cụ thể**: Càng nhiều thông tin càng tốt
- **Business rules rõ ràng**: Mô tả logic nghiệp vụ cụ thể
- **Error scenarios đầy đủ**: Cover tất cả case có thể xảy ra
- **Performance requirements thực tế**: Đưa ra số liệu hợp lý
- **Security considerations**: Luôn nghĩ về security

### ❌ DON'T
- **Mơ hồ**: "Create something for users"
- **Thiếu validation rules**: AI không biết validate gì
- **Bỏ qua error handling**: Dẫn đến code không robust
- **Quên performance**: API sẽ chậm và không scalable
- **Bỏ qua security**: Tạo ra security vulnerabilities

## 🔗 Related Documents

- [AI_API_GENERATION_RULES.md](./AI_API_GENERATION_RULES.md) - Chi tiết rules cho AI
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Clean Architecture overview
- [DEVELOPMENT.md](./DEVELOPMENT.md) - Development workflow
- [API.md](./API.md) - API documentation standards