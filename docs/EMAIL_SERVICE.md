# Email Service Documentation

## 📧 Overview

The Go MVC application includes a comprehensive email service system that supports both SMTP and API-based email delivery. The service is designed with clean architecture principles, providing flexible email functionality for authentication, notifications, and user communication.

## 📋 Table of Contents
- [Architecture](#architecture)
- [Configuration](#configuration)
- [Development Setup](#development-setup)
- [Email Features](#email-features)
- [Implementation Details](#implementation-details)
- [Testing](#testing)
- [Production Deployment](#production-deployment)
- [Troubleshooting](#troubleshooting)

## 🏗️ Architecture

### Service Structure

```
internal/
├── adapters/external/
│   ├── email_service.go     # API-based email service (SendGrid, SES)
│   └── smtp_service.go      # SMTP email service implementation
├── application/services/
│   └── auth_service.go      # Integration with authentication flows
└── di/modules/
    └── auth.go              # Dependency injection configuration
```

### Email Service Types

| Service Type | Use Case | Configuration |
|-------------|----------|---------------|
| **SMTP Service** | Development, small-scale production | Direct SMTP server connection |
| **API Service** | Large-scale production | Third-party email providers (SendGrid, SES) |

## ⚙️ Configuration

### Environment Configuration

**Development** (`configs/development.yaml`):
```yaml
external:
  email_service:
    provider: "smtp"              # Use SMTP for development
    smtp:
      host: "localhost"           # MailCatcher host
      port: 1025                  # MailCatcher SMTP port
      username: ""                # No auth for MailCatcher
      password: ""                # No auth for MailCatcher  
      from: "noreply@go-mvc.dev"  # From address
      tls: false                  # No TLS for local development
    api_key: ""                   # Not used for SMTP
    from: "noreply@go-mvc.dev"    # Default from address
```

**Production** (`configs/production.yaml`):
```yaml
external:
  email_service:
    provider: "api"               # Use API service for production
    smtp:
      host: ""                    # Not used for API
      port: 0                     # Not used for API
      username: ""                # Not used for API
      password: ""                # Not used for API
      from: "noreply@yourapp.com" # Production from address
      tls: false                  # Not used for API
    api_key: "${SENDGRID_API_KEY}" # API key from environment
    from: "noreply@yourapp.com"    # Production from address
```

### Environment Variables

```bash
# For production API-based email service
export SENDGRID_API_KEY="your-sendgrid-api-key"
export EMAIL_FROM="noreply@yourapp.com"
```

## 🛠️ Development Setup

### 1. MailCatcher Integration

MailCatcher is automatically configured for development email testing:

```bash
# MailCatcher services (automatically started with docker-compose up)
# Web Interface: http://localhost:1080
# SMTP Server: localhost:1025
```

### 2. Starting Development Environment

```bash
# Start all services including MailCatcher
docker-compose up -d

# Verify MailCatcher is running
curl -s http://localhost:1080/messages
```

### 3. Testing Email Functionality

```bash
# Register a new user (may trigger verification email)
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testPassword123",
    "name": "Test User"
  }'

# Request password reset email
curl -X POST http://localhost:8080/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com"}'

# Request verification email
curl -X POST http://localhost:8080/api/v1/auth/resend-verification \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com"}'
```

### 4. Viewing Captured Emails

Open http://localhost:1080 in your browser to view all captured emails.

## 📧 Email Features

### Authentication Email Flows

| Feature | Endpoint | Email Type | Description |
|---------|----------|------------|-------------|
| **Password Reset** | `POST /api/v1/auth/reset-password` | Password Reset | Sends reset link with token |
| **Email Verification** | `POST /api/v1/auth/resend-verification` | Verification | Sends verification link |
| **Password Confirm** | `POST /api/v1/auth/confirm-reset` | - | Processes reset token (no email) |
| **Email Verify** | `POST /api/v1/auth/verify-email` | - | Processes verification token (no email) |

### Email Templates

#### Password Reset Email
```
Subject: Password Reset Request

Hello {Name},

You requested a password reset. Please click the link below to reset your password:

http://localhost:8080/api/v1/auth/confirm-reset?token={ResetToken}

This link will expire in 1 hour. If you didn't request this, please ignore this email.

Best regards,
The Team
```

#### Email Verification Email
```
Subject: Verify Your Email Address

Hello {Name},

Please click the link below to verify your email address:

http://localhost:8080/api/v1/auth/verify-email?token={VerificationToken}

If you didn't create an account, please ignore this email.

Best regards,
The Team
```

## 💻 Implementation Details

### SMTP Service Implementation

**File**: `internal/adapters/external/smtp_service.go`

**Key Features**:
- Direct SMTP connection using Go's `net/smtp` package
- Support for authenticated and non-authenticated SMTP
- TLS support for secure connections
- Built-in email templates for common use cases
- Error handling and logging
- Clean shutdown support

**Example Usage**:
```go
// Dependency injection
smtpService := external.NewSMTPService(&config.SMTPConfig{
    Host: "localhost",
    Port: "1025",
    From: "noreply@go-mvc.dev",
}, logger)

// Send email
err := smtpService.SendPasswordResetEmail(ctx, "user@example.com", "John", "reset-token")
```

### Integration with Authentication Service

**File**: `internal/application/services/auth_service.go`

The SMTP service is integrated into authentication workflows:

```go
// Password reset flow
func (s *authService) ResetPassword(ctx context.Context, email string) error {
    // ... token generation logic ...
    
    // Send password reset email
    if err := s.smtpService.SendPasswordResetEmail(ctx, userEntity.Email(), userEntity.Name(), resetToken); err != nil {
        s.logger.Errorf("Failed to send password reset email: %v", err)
        // Continue - we don't want to fail the reset process if email fails
    }
    
    return nil
}
```

### Dependency Injection

**File**: `internal/di/modules/auth.go`

```go
// SMTP service provider
func NewSMTPService(cfg *config.AppConfig, logger *logger.Logger) *external.SMTPService {
    return external.NewSMTPService(&cfg.External.EmailService.SMTP, logger)
}

// Auth service with SMTP dependency
func NewAuthService(params AuthServiceParams) services.AuthService {
    return appServices.NewAuthService(
        params.UserRepo,
        params.JWTService,
        params.PasswordHasher,
        params.CacheService,
        params.SMTPService, // SMTP service injection
        params.Logger,
    )
}
```

## 🧪 Testing

### Unit Testing

**Testing SMTP Service**:
```go
func TestSMTPService_SendEmail(t *testing.T) {
    // Mock SMTP server for testing
    mockServer := startMockSMTPServer()
    defer mockServer.Close()
    
    service := external.NewSMTPService(&config.SMTPConfig{
        Host: mockServer.Host,
        Port: mockServer.Port,
        From: "test@example.com",
    }, logger)
    
    err := service.SendEmail(ctx, []string{"user@example.com"}, "Test", "Body")
    assert.NoError(t, err)
}
```

### Integration Testing

**Testing Email Flows**:
```go
func TestAuthService_ResetPassword(t *testing.T) {
    // Setup test dependencies with mock SMTP service
    mockSMTP := &MockSMTPService{}
    authService := NewAuthService(userRepo, jwtService, hasher, cache, mockSMTP, logger)
    
    err := authService.ResetPassword(ctx, "user@example.com")
    
    assert.NoError(t, err)
    assert.True(t, mockSMTP.EmailSent)
    assert.Equal(t, "user@example.com", mockSMTP.LastRecipient)
}
```

### Manual Testing with MailCatcher

1. **Start MailCatcher**: `docker-compose up -d`
2. **Trigger Email**: Make API requests to email endpoints
3. **Verify Delivery**: Check http://localhost:1080
4. **Test Content**: Verify email templates and links work correctly

## 🚀 Production Deployment

### Environment Setup

```bash
# Set production email service configuration
export EMAIL_PROVIDER="api"
export SENDGRID_API_KEY="your-production-api-key"
export EMAIL_FROM="noreply@yourproductiondomain.com"

# Or use SMTP for production
export EMAIL_PROVIDER="smtp"
export SMTP_HOST="your-smtp-server.com"
export SMTP_PORT="587"
export SMTP_USERNAME="your-smtp-username"
export SMTP_PASSWORD="your-smtp-password"
export SMTP_FROM="noreply@yourproductiondomain.com"
export SMTP_TLS="true"
```

### Production Considerations

1. **Email Provider Choice**:
   - **SendGrid**: Recommended for high volume, good deliverability
   - **AWS SES**: Cost-effective, good AWS integration
   - **SMTP**: Good for existing email infrastructure

2. **Security**:
   - Use environment variables for sensitive data
   - Enable TLS for SMTP connections
   - Validate email addresses before sending
   - Implement rate limiting for email sending

3. **Monitoring**:
   - Track email delivery success/failure rates
   - Monitor email bounce rates
   - Set up alerts for email service failures

4. **Performance**:
   - Email sending should not block user operations
   - Consider async email sending for high-volume applications
   - Implement retry logic for failed email deliveries

## 🔧 Troubleshooting

### Common Issues

#### 1. MailCatcher Not Receiving Emails

**Symptoms**: Emails not appearing in MailCatcher interface
**Solutions**:
```bash
# Check MailCatcher container status
docker ps | grep mailcatcher

# Restart MailCatcher
docker-compose restart mailcatcher

# Check SMTP configuration
curl -s http://localhost:1080/messages
```

#### 2. SMTP Connection Errors

**Symptoms**: Connection refused, timeout errors
**Solutions**:
```bash
# Test SMTP connection manually
telnet localhost 1025

# Check network connectivity
docker network ls
docker network inspect go-mvc_dev-network

# Verify configuration
grep -A 10 "smtp:" configs/development.yaml
```

#### 3. Email Template Issues

**Symptoms**: Malformed emails, missing content
**Solutions**:
- Check template rendering in `smtp_service.go`
- Verify user data (name, email) availability
- Test with minimal template first

#### 4. Production Email Delivery Issues

**Symptoms**: Emails not delivered, marked as spam
**Solutions**:
- Verify SPF, DKIM, DMARC records
- Check email provider logs
- Validate sender reputation
- Use email testing tools (Mail Tester, etc.)

### Debugging Commands

```bash
# Check email service logs
make logs | grep -i "email\|smtp"

# Test email endpoints directly
curl -X POST http://localhost:8080/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"email": "debug@example.com"}' -v

# View MailCatcher messages via API
curl -s http://localhost:1080/messages.json | jq .

# Check specific email content
curl -s http://localhost:1080/messages/1.plain
```

### Configuration Validation

```bash
# Validate YAML configuration
docker run --rm -v $(pwd)/configs:/configs mikefarah/yq:latest \
  eval '.external.email_service' /configs/development.yaml

# Test configuration loading
go run cmd/main.go --config-check
```

## 📚 Additional Resources

- [MailCatcher Documentation](https://mailcatcher.me/)
- [Go net/smtp Package](https://pkg.go.dev/net/smtp)
- [SendGrid Go SDK](https://github.com/sendgrid/sendgrid-go)
- [Email Deliverability Best Practices](https://sendgrid.com/blog/email-deliverability-best-practices/)