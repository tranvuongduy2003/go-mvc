# AI API Generation Quick Start Guide

## 🚀 Cách Sử Dụng AI để Sinh API Tự Động

### 📋 Tóm Tắt

Dự án này đã được cấu hình để AI có thể tự động sinh ra một bộ API hoàn chỉnh chỉ từ User Story. Bạn chỉ cần:

1. **Viết User Story** theo template có sẵn
2. **Đưa cho AI** với instruction đơn giản
3. **Nhận code hoàn chỉnh** với tất cả layers
4. **Review và test** code được sinh ra
5. **Deploy** lên production

## ⚡ Quick Start (5 phút)

### Bước 1: Chuẩn Bị User Story (2 phút)

```bash
# Copy template
cp docs/USER_STORY_TEMPLATE.md my_feature.md

# Hoặc sử dụng example có sẵn
cp docs/USER_STORY_TEMPLATE.md product_management.md
```

Điền thông tin theo format:

```markdown
## User Story: Create Product

### Business Description
- **Actor**: Admin User
- **Action**: Create
- **Object**: Product
- **Purpose**: Allow admin to add new products to catalog

### Functional Requirements
#### Inputs
name: string (required, 2-100 chars) - Product name
price: decimal (required, > 0) - Product price in USD
category_id: uuid (required, must exist) - Reference to category
description: string (optional, max 1000 chars) - Product description

#### Business Rules
- Product name must be unique within same category
- Price must be positive number
- Category must exist and be active
- Only admin users can create products

### Technical Specifications
- **HTTP Method**: POST
- **Endpoint**: /api/v1/products
- **Authentication**: Required (JWT)
- **Authorization**: Admin role required

### Error Scenarios
#### Client Errors (4xx)
400: Invalid JSON format, validation errors
401: No authentication token
403: User is not admin
409: Product name already exists in category

### Performance Requirements
- **Response Time**: < 500ms
- **Throughput**: 100 RPS
```

### Bước 2: Sử Dụng AI (1 phút)

Copy instruction này và kèm theo User Story:

```markdown
Hãy sinh ra một bộ API hoàn chỉnh theo Clean Architecture từ User Story bên dưới.
Sử dụng các quy tắc trong docs/AI_API_GENERATION_RULES.md và docs/CODE_GENERATION_GUIDELINES.md.
Dự án sử dụng Go với Gin framework, GORM ORM, PostgreSQL database, và FX dependency injection.

[Dán User Story ở đây]
```

### Bước 3: Review Code (2 phút)

AI sẽ sinh ra code cho các files:

```
✅ Domain Layer:
   ├── internal/core/domain/product/product.go
   ├── internal/core/domain/product/value_objects.go
   └── internal/core/ports/repositories/product_repository.go

✅ Application Layer:
   ├── internal/application/commands/product/create_product_command.go
   ├── internal/application/queries/product/get_product_query.go
   ├── internal/application/dto/product/product_dto.go
   ├── internal/application/services/product_service.go
   └── internal/application/validators/product/product_validator.go

✅ Infrastructure Layer:
   ├── internal/adapters/persistence/postgres/models/product.go
   ├── internal/adapters/persistence/postgres/repositories/product_repository.go
   └── internal/adapters/persistence/postgres/migrations/001_create_products_table.up.sql

✅ Presentation Layer:
   ├── internal/handlers/http/rest/v1/product_handler.go
   └── Route setup code

✅ Integration:
   └── internal/di/modules/product.go
```

Review checklist:
- [ ] Tất cả files được tạo đúng vị trí
- [ ] Business rules được implement trong domain layer
- [ ] Validation rules đầy đủ và chính xác
- [ ] Database constraints và indexes hợp lý
- [ ] Authorization checks đúng requirements
- [ ] Error handling cover tất cả scenarios

## 📚 Tài Liệu Chi Tiết

### Core Documents
- **[AI_API_GENERATION_RULES.md](./AI_API_GENERATION_RULES.md)**: Quy tắc toàn diện cho AI
- **[USER_STORY_TEMPLATE.md](./USER_STORY_TEMPLATE.md)**: Template và examples
- **[CODE_GENERATION_GUIDELINES.md](./CODE_GENERATION_GUIDELINES.md)**: Guidelines cho từng layer
- **[DEVELOPMENT.md](./DEVELOPMENT.md)**: Quy trình development với AI

### Architecture Documents
- **[ARCHITECTURE.md](./ARCHITECTURE.md)**: Clean Architecture overview
- **[PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md)**: Cấu trúc project
- **[API.md](./API.md)**: API documentation standards

## 💡 Tips cho Kết Quả Tốt Nhất

### ✅ User Story Nên Có
- **Chi tiết cụ thể**: Mô tả rõ ràng từng field và validation
- **Business rules đầy đủ**: Tất cả constraints và logic
- **Error scenarios complete**: Cover tất cả lỗi có thể xảy ra
- **Authorization rõ ràng**: Ai có quyền làm gì
- **Performance requirements thực tế**: Response time, throughput

### ❌ Tránh
- **Mơ hồ**: "Create something for users"
- **Thiếu validation**: AI sẽ không biết validate gì
- **Bỏ qua errors**: Code sẽ không robust
- **Không có authorization**: Tạo ra security holes
- **Performance không rõ**: API sẽ chậm

## 🎯 Examples & Use Cases

### Example 1: Simple CRUD
```markdown
## User Story: Manage Categories

### Business Description
- **Actor**: Admin User
- **Action**: CRUD Operations
- **Object**: Product Categories
- **Purpose**: Organize products into categories

### Functional Requirements
#### Inputs (Create/Update)
name: string (required, 2-50 chars, unique) - Category name
description: string (optional, max 500 chars) - Category description
is_active: boolean (default true) - Category status

#### Business Rules
- Category name must be unique globally
- Cannot delete category if products exist
- Only admin can modify categories

### Technical Specifications
- **Endpoints**: 
  - POST /api/v1/categories
  - GET /api/v1/categories
  - GET /api/v1/categories/{id}
  - PUT /api/v1/categories/{id}
  - DELETE /api/v1/categories/{id}
- **Authentication**: Required
- **Authorization**: Admin role required
```

### Example 2: Complex Business Logic
```markdown
## User Story: Process Order Payment

### Business Description
- **Actor**: System (Payment Gateway Webhook)
- **Action**: Process Payment
- **Object**: Order Payment
- **Purpose**: Update order status based on payment result

### Functional Requirements
#### Inputs
order_id: uuid (required, must exist, must be pending_payment)
payment_status: enum (success, failed, pending) - Payment result
payment_id: string (required) - External payment reference
amount: decimal (required, must match order total) - Payment amount

#### Business Rules
- Order must be in pending_payment status
- Payment amount must exactly match order total
- Successful payment → order status = paid, inventory confirmed
- Failed payment → order status = payment_failed, inventory released
- Pending payment → no status change, schedule retry

### Integration Requirements
#### External APIs
- Payment gateway for verification
- Inventory service for reservation/release
- Email service for notifications

#### Message Queue
- Queue inventory confirmation/release
- Queue customer notification emails
- Queue order status change events
```

### Example 3: File Upload Feature
```markdown
## User Story: Upload Product Images

### Business Description
- **Actor**: Admin User
- **Action**: Upload
- **Object**: Product Images
- **Purpose**: Allow attaching images to products for better catalog presentation

### Functional Requirements
#### Inputs
product_id: uuid (required, must exist) - Product reference
image_file: file (required, jpg/png, max 5MB) - Image file
alt_text: string (optional, max 200 chars) - Alternative text
is_primary: boolean (default false) - Primary image flag

#### Business Rules
- Product must exist and be editable by user
- Only one primary image per product
- Maximum 10 images per product
- Image must be valid format và size
- Generate thumbnails automatically

### Technical Specifications
- **HTTP Method**: POST multipart/form-data
- **Endpoint**: /api/v1/products/{id}/images
- **Authentication**: Required
- **Authorization**: Admin role or product owner

### Performance Requirements
- **Response Time**: < 2s for upload processing
- **File Size**: Max 5MB per image
- **Concurrent Uploads**: Support 10 concurrent uploads
- **Storage**: Use cloud storage with CDN
```

## 🔧 Advanced Scenarios

### Multi-Entity Features
Khi User Story liên quan đến nhiều entities:

```markdown
## User Story: Complete Order Checkout

### Entities Involved
- Order (main aggregate)
- OrderItem (part of order)
- Payment (separate aggregate)  
- ShippingAddress (value object)
- Customer (existing entity)

### Business Rules Across Entities
- Order total = sum of OrderItem totals + shipping + tax
- Payment amount must match order total exactly
- Customer must be authenticated and verified
- Shipping address must be valid và deliverable
- Inventory must be available for all items
```

**AI sẽ sinh ra:**
- Separate aggregates với clear boundaries
- Cross-aggregate validation logic
- Saga pattern for distributed transactions
- Event-driven communication between aggregates
- Compensating transactions for failure scenarios

### Event-Driven Features
Cho features cần async processing:

```markdown
### Integration Requirements
#### Domain Events
- OrderCreated → trigger inventory reservation
- PaymentCompleted → trigger order fulfillment
- OrderShipped → trigger customer notification
- PaymentFailed → trigger inventory release

#### Message Handlers
- InventoryReservationHandler
- OrderFulfillmentHandler
- CustomerNotificationHandler
- InventoryReleaseHandler
```

## 🚀 Deployment & Production

### Testing Generated Code
```bash
# Unit tests
make test-unit

# Integration tests  
make test-integration

# Load testing
make load-test

# Security scan
make security-scan
```

### Performance Optimization
AI-generated code includes:
- Database indexes for query optimization
- Caching strategies for frequently accessed data
- Pagination for large result sets
- Rate limiting for API protection
- Connection pooling for database efficiency

### Monitoring & Observability
Generated code includes:
- Prometheus metrics collection
- Structured logging với context
- Distributed tracing support
- Health checks for dependencies
- Error tracking và alerting

## 📞 Support & Troubleshooting

### Common Issues

**Issue**: AI generates code but missing some files
**Solution**: Check User Story completeness, re-run with more specific requirements

**Issue**: Generated code doesn't compile
**Solution**: Review import statements, check Go version compatibility

**Issue**: Database migration fails
**Solution**: Check for naming conflicts, verify PostgreSQL version

**Issue**: Authorization not working
**Solution**: Verify JWT middleware setup, check role definitions

**Issue**: Performance issues
**Solution**: Review database indexes, check N+1 query problems

### Getting Help

1. **Check Documentation**: Review all docs in `/docs` folder
2. **Validate User Story**: Ensure it follows the template completely
3. **Review Generated Code**: Use the provided checklists
4. **Test Incrementally**: Test each layer separately
5. **Check Existing Examples**: Look at current working code in the project

## 📝 Changelog & Updates

### v1.0 - Initial AI Generation System
- Complete Clean Architecture code generation
- User Story template và validation
- Comprehensive documentation
- Integration với existing project structure
- Testing và quality assurance workflows

### Future Enhancements
- Visual User Story builder
- Real-time code generation preview
- Automated testing generation
- Performance benchmarking integration
- Multi-language API client generation

---

**🎉 Bây giờ bạn đã sẵn sàng sử dụng AI để sinh API tự động! Chỉ cần 1 User Story chi tiết, bạn sẽ có ngay 1 bộ API production-ready.**