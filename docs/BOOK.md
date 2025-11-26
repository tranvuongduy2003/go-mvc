# Go MVC - Complete Developer Guide

**Enterprise-Grade Go Web Application with Clean Architecture**

Version 1.0.0

---

## About This Book

This comprehensive guide covers everything you need to know to build, deploy, and maintain enterprise-grade Go web applications using Clean Architecture, Domain-Driven Design, and modern development practices.

## Who This Book Is For

- **Backend Developers** looking to master Go web development
- **Architects** designing scalable microservices
- **DevOps Engineers** deploying Go applications
- **AI/ML Engineers** leveraging AI-powered code generation
- **Teams** adopting Clean Architecture and DDD

## What You'll Learn

- **Clean Architecture** principles and implementation
- **Domain-Driven Design** with rich domain models
- **CQRS and Event Sourcing** patterns
- **AI-Powered Development** with automated code generation
- **Enterprise Patterns** (Repository, Unit of Work, Specification)
- **Observability** with Prometheus, Grafana, and Jaeger
- **Testing Strategies** from unit to integration tests
- **Deployment** with Docker and Kubernetes

---

## Table of Contents

### Part I: Getting Started

#### [Chapter 1: Introduction](book/01-getting-started/01-introduction.md)
- Project Overview
- Key Features
- Technology Stack
- Prerequisites

#### [Chapter 2: Quick Start](book/01-getting-started/02-quick-start.md)
- Installation
- Running the Application
- First API Call
- Development Environment Setup

#### [Chapter 3: Quick Reference](book/01-getting-started/03-quick-reference.md)
- Common Commands
- Project Structure Overview
- Configuration Guide
- Troubleshooting

---

### Part II: Architecture

#### [Chapter 4: Clean Architecture Overview](book/02-architecture/01-architecture-overview.md)
- Clean Architecture Principles
- Layer Responsibilities
- Dependency Rules
- Data Flow

#### [Chapter 5: Project Structure](book/02-architecture/02-project-structure.md)
- Directory Organization
- Module Layout
- Package Design
- Naming Conventions

#### [Chapter 6: Domain Layer](book/02-architecture/03-domain-layer.md)
- Entities and Aggregates
- Value Objects
- Domain Events
- Repository Interfaces

#### [Chapter 7: Application Layer](book/02-architecture/04-application-layer.md)
- Commands and Queries
- Use Case Handlers
- DTOs and Mapping
- Application Services

#### [Chapter 8: Infrastructure Layer](book/02-architecture/05-infrastructure-layer.md)
- Database Integration
- External Services
- Repository Implementation
- Message Queuing

#### [Chapter 9: Presentation Layer](book/02-architecture/06-presentation-layer.md)
- HTTP Handlers
- Middleware
- Request/Response Flow
- API Versioning

#### [Chapter 10: Dependency Injection](book/02-architecture/07-dependency-injection.md)
- Uber FX Framework
- Module Organization
- Lifecycle Management
- Testing with DI

#### [Chapter 11: Architecture Examples](book/02-architecture/08-architecture-examples.md)
- Real-world Scenarios
- Best Practices
- Common Patterns
- Anti-patterns to Avoid

---

### Part III: Development Guide

#### [Chapter 12: Development Workflow](book/03-development-guide/01-development-workflow.md)
- Setting Up Development Environment
- Hot Reload with Air
- Debugging Techniques
- Development Tools

#### [Chapter 13: Database Migrations](book/03-development-guide/02-migrations.md)
- Migration Strategy
- Creating Migrations
- Running Migrations
- Rollback Procedures

#### [Chapter 14: Testing](book/03-development-guide/03-testing.md)
- Testing Philosophy
- Unit Tests
- Integration Tests
- Test Coverage

#### [Chapter 15: API Development](book/03-development-guide/04-api-development.md)
- RESTful API Design
- Request Validation
- Error Handling
- API Documentation

---

### Part IV: Features & Integration

#### [Chapter 16: Authentication & Authorization](book/04-features/01-authentication.md)
- JWT Authentication
- Role-Based Access Control (RBAC)
- Permission System
- Security Best Practices

#### [Chapter 17: Background Jobs](book/04-features/02-background-jobs.md)
- Job Queue Architecture
- Creating Jobs
- Job Scheduling
- Error Handling

#### [Chapter 18: Email Service](book/04-features/03-email-service.md)
- Email Configuration
- Template System
- Sending Emails
- Email Testing

#### [Chapter 19: File Upload](book/04-features/04-file-upload.md)
- Upload Configuration
- Storage Backends
- File Validation
- Security Considerations

#### [Chapter 20: Message Deduplication](book/04-features/05-message-deduplication.md)
- Inbox/Outbox Pattern
- Idempotency
- Event Processing
- Consistency Guarantees

#### [Chapter 21: NATS Messaging](book/04-features/06-nats-messaging.md)
- NATS Integration
- Pub/Sub Patterns
- Request/Reply
- Message Persistence

#### [Chapter 22: Distributed Tracing](book/04-features/07-tracing.md)
- OpenTelemetry Setup
- Jaeger Integration
- Trace Context
- Performance Analysis

---

### Part V: AI-Powered Development

#### [Chapter 23: AI Quick Start](book/05-ai-development/01-ai-quick-start.md)
- AI Development Overview
- Prerequisites
- First AI-Generated API
- Workflow

#### [Chapter 24: API Generation Rules](book/05-ai-development/02-api-generation-rules.md)
- User Story Template
- Generation Process
- Layer-by-Layer Guidelines
- Code Conventions

#### [Chapter 25: Code Generation Guidelines](book/05-ai-development/03-code-generation-guidelines.md)
- Domain Layer Generation
- Application Layer Generation
- Infrastructure Layer Generation
- Presentation Layer Generation
- Integration & Testing

---

### Part VI: Deployment & Operations

#### [Chapter 26: Deployment Guide](book/06-operations/01-deployment.md)
- Deployment Strategies
- Docker Setup
- Kubernetes Deployment
- CI/CD Pipeline

#### [Chapter 27: Monitoring & Observability](book/06-operations/02-monitoring.md)
- Prometheus Metrics
- Grafana Dashboards
- Alerting Rules
- Log Aggregation

#### [Chapter 28: Performance Optimization](book/06-operations/03-performance.md)
- Profiling
- Caching Strategies
- Database Optimization
- Scaling Strategies

---

## Appendices

### [Appendix A: AI Coding Standards](book/appendix/AI_CODING_STANDARDS.md) ⭐ **NEW**
- Self-Documenting Code Principles
- Comment Policy (Minimal Comments Only)
- AI Code Generation Rules
- Code Review Checklist

### [Appendix B: Glossary](book/appendix/glossary.md)
- Technical Terms
- Acronyms
- Concepts

### [Appendix C: Resources](book/appendix/resources.md)
- Further Reading
- Online Resources
- Community Links

### [Appendix D: Migration from Legacy](book/appendix/migration.md)
- Migration Strategies
- Step-by-Step Guide
- Common Challenges

---

## How to Use This Book

### For Beginners
Start with **Part I** to understand the basics, then move through **Part II** to grasp the architecture. Practice with **Part III** before exploring advanced features.

### For Experienced Developers
Jump to **Part II** for architecture details, then explore specific features in **Part IV** and **Part V** for AI-powered development.

### For DevOps Engineers
Focus on **Part VI** for deployment and operations, but review **Part II** to understand the application structure.

### For Teams
Use this book as a reference guide and onboarding material. Each chapter can be studied independently.

---

## Conventions Used in This Book

**Code Examples**
```go
// Code examples are syntax-highlighted
func Example() {}
```

**Terminal Commands**
```bash
# Shell commands start with $
$ make build
```

**Important Notes**
> 💡 **Note**: Important information is highlighted like this.

**Warnings**
> ⚠️ **Warning**: Critical information appears in warning boxes.

**Best Practices**
> ✅ **Best Practice**: Recommended approaches are marked like this.

---

## About the Author

This book is maintained by the Go MVC project team and contributors.

## Contributing

Found an error? Want to improve the documentation? Contributions are welcome!
- GitHub: https://github.com/tranvuongduy2003/go-mvc
- Issues: https://github.com/tranvuongduy2003/go-mvc/issues

---

## License

This documentation is part of the Go MVC project, licensed under MIT License.

Copyright © 2024-2025 Go MVC Project

---

**Ready to start?** Begin with [Chapter 1: Introduction](book/01-getting-started/01-introduction.md)
