# API Documentation

## 📋 Table of Contents
- [Overview](#overview)
- [Authentication](#authentication)
- [Base URLs](#base-urls)
- [HTTP Status Codes](#http-status-codes)
- [Response Format](#response-format)
- [Rate Limiting](#rate-limiting)
- [API Endpoints](#api-endpoints)
- [Error Handling](#error-handling)
- [Examples](#examples)

## 🌐 Overview

The Go MVC API is a RESTful web service that provides comprehensive functionality for user management, authentication, and business operations. The API follows RESTful conventions and returns JSON responses.

### API Characteristics
- **Protocol**: HTTP/HTTPS
- **Data Format**: JSON
- **Authentication**: JWT Bearer Tokens
- **Versioning**: URL-based versioning (`/api/v1/`)
- **Rate Limiting**: 100 requests per minute per IP
- **CORS**: Configurable cross-origin support

## 🔐 Authentication

### JWT Authentication
The API uses JWT (JSON Web Token) for authentication. Include the token in the Authorization header:

```http
Authorization: Bearer <your-jwt-token>
```

### Getting a Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-01T12:00:00Z",
    "user": {
      "id": "uuid-here",
      "email": "user@example.com",
      "name": "John Doe"
    }
  },
  "message": "Login successful"
}
```

## 🌍 Base URLs

| Environment | Base URL |
|-------------|----------|
| Development | `http://localhost:8080/api/v1` |
| Staging | `https://staging-api.example.com/api/v1` |
| Production | `https://api.example.com/api/v1` |

## 📊 HTTP Status Codes

| Code | Status | Description |
|------|--------|-------------|
| 200 | OK | Request successful |
| 201 | Created | Resource created successfully |
| 204 | No Content | Request successful, no content returned |
| 400 | Bad Request | Invalid request parameters |
| 401 | Unauthorized | Authentication required |
| 403 | Forbidden | Access denied |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Resource already exists |
| 422 | Unprocessable Entity | Validation error |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |

## 📄 Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "success": true,
  "data": {
    // Response data here
  },
  "message": "Operation successful",
  "meta": {
    "timestamp": "2024-01-01T12:00:00Z",
    "request_id": "req-uuid-here"
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "email",
        "message": "Email is required"
      }
    ]
  },
  "meta": {
    "timestamp": "2024-01-01T12:00:00Z",
    "request_id": "req-uuid-here"
  }
}
```

### Pagination Response
```json
{
  "success": true,
  "data": [
    // Array of items
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10,
    "has_next": true,
    "has_prev": false
  }
}
```

## 🚦 Rate Limiting

Rate limiting is implemented to prevent abuse:

- **Limit**: 100 requests per minute per IP address
- **Headers**: Rate limit information is included in response headers
- **Reset**: Limit resets every minute

### Rate Limit Headers
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

## 🛠️ API Endpoints

### Health & Monitoring

#### Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "services": {
    "database": "healthy",
    "redis": "healthy",
    "external_api": "healthy"
  }
}
```

#### Metrics (Prometheus)
```http
GET /metrics
```

**Response:** Prometheus metrics format

#### Trace Test
```http
GET /api/v1/trace-test
```

**Purpose:** Generate a test trace for monitoring

### Authentication

#### Register User
```http
POST /api/v1/auth/register
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe",
  "phone": "+1234567890"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid-here",
      "email": "user@example.com",
      "name": "John Doe",
      "created_at": "2024-01-01T12:00:00Z"
    }
  },
  "message": "User registered successfully"
}
```

#### Login
```http
POST /api/v1/auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Logout
```http
POST /api/v1/auth/logout
```

**Headers:**
```http
Authorization: Bearer <token>
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
```

**Headers:**
```http
Authorization: Bearer <refresh-token>
```

### User Management

#### Get Current User
```http
GET /api/v1/users/me
```

**Headers:**
```http
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid-here",
    "email": "user@example.com",
    "name": "John Doe",
    "phone": "+1234567890",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

#### Update User Profile
```http
PUT /api/v1/users/me
```

**Headers:**
```http
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "John Smith",
  "phone": "+1234567891"
}
```

#### Get User by ID
```http
GET /api/v1/users/{id}
```

**Headers:**
```http
Authorization: Bearer <token>
```

#### List Users (Admin)
```http
GET /api/v1/users?page=1&limit=10&search=john
```

**Headers:**
```http
Authorization: Bearer <admin-token>
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)
- `search` (optional): Search term for name/email
- `sort` (optional): Sort field (name, email, created_at)
- `order` (optional): Sort order (asc, desc)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid-1",
      "email": "user1@example.com",
      "name": "User One",
      "created_at": "2024-01-01T12:00:00Z"
    },
    {
      "id": "uuid-2", 
      "email": "user2@example.com",
      "name": "User Two",
      "created_at": "2024-01-01T11:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3,
    "has_next": true,
    "has_prev": false
  }
}
```

#### Delete User (Admin)
```http
DELETE /api/v1/users/{id}
```

**Headers:**
```http
Authorization: Bearer <admin-token>
```

### Password Management

#### Change Password
```http
POST /api/v1/auth/change-password
```

**Headers:**
```http
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword123"
}
```

#### Forgot Password
```http
POST /api/v1/auth/forgot-password
```

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

#### Reset Password
```http
POST /api/v1/auth/reset-password
```

**Request Body:**
```json
{
  "token": "reset-token-here",
  "new_password": "newpassword123"
}
```

## ❌ Error Handling

### Error Codes

| Code | Description |
|------|-------------|
| `VALIDATION_ERROR` | Request validation failed |
| `AUTHENTICATION_ERROR` | Authentication failed |
| `AUTHORIZATION_ERROR` | Access denied |
| `NOT_FOUND` | Resource not found |
| `CONFLICT` | Resource already exists |
| `RATE_LIMIT_EXCEEDED` | Too many requests |
| `INTERNAL_ERROR` | Server error |

### Validation Errors
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "email",
        "message": "Email is required"
      },
      {
        "field": "password",
        "message": "Password must be at least 8 characters"
      }
    ]
  }
}
```

### Authentication Errors
```json
{
  "success": false,
  "error": {
    "code": "AUTHENTICATION_ERROR",
    "message": "Invalid credentials"
  }
}
```

## 📝 Examples

### Complete User Registration Flow

1. **Register a new user:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123",
    "name": "John Doe"
  }'
```

2. **Login to get token:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

3. **Get user profile:**
```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer <token-from-login>"
```

4. **Update profile:**
```bash
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "John Smith",
    "phone": "+1234567890"
  }'
```

### Pagination Example

```bash
# Get first page of users
curl -X GET "http://localhost:8080/api/v1/users?page=1&limit=5" \
  -H "Authorization: Bearer <admin-token>"

# Search users
curl -X GET "http://localhost:8080/api/v1/users?search=john&page=1&limit=10" \
  -H "Authorization: Bearer <admin-token>"

# Sort users by creation date
curl -X GET "http://localhost:8080/api/v1/users?sort=created_at&order=desc" \
  -H "Authorization: Bearer <admin-token>"
```

### Error Handling Example

```bash
# Invalid request (missing required fields)
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email"
  }'

# Response:
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "email",
        "message": "Email format is invalid"
      },
      {
        "field": "password",
        "message": "Password is required"
      },
      {
        "field": "name",
        "message": "Name is required"
      }
    ]
  }
}
```

## 🧪 Testing the API

### Using cURL

Test all endpoints with the provided cURL examples above.

### Using Postman

1. Import the OpenAPI specification from `/api/openapi/`
2. Set up environment variables for base URL and tokens
3. Use the pre-configured requests

### Using HTTPie

```bash
# Install HTTPie
pip install httpie

# Register user
http POST localhost:8080/api/v1/auth/register \
  email=test@example.com password=test123 name="Test User"

# Login
http POST localhost:8080/api/v1/auth/login \
  email=test@example.com password=test123

# Get profile (replace TOKEN with actual token)
http GET localhost:8080/api/v1/users/me \
  Authorization:"Bearer TOKEN"
```

## 📚 Additional Resources

- **OpenAPI Specification**: Available at `/api/openapi/swagger.json`
- **Interactive Documentation**: Swagger UI at `/api/docs` (if enabled)
- **Monitoring**: Grafana dashboard at `http://localhost:3000`
- **Tracing**: Jaeger UI at `http://localhost:16686`

For more detailed information about the API implementation, see the [Architecture Documentation](ARCHITECTURE.md).