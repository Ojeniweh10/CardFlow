# CardFlow - Virtual Card Issuance Platform

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Fiber](https://img.shields.io/badge/Fiber-2.51+-00ACD7?style=flat)](https://gofiber.io)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-316192?style=flat&logo=postgresql)](https://postgresql.org)
[![Redis](https://img.shields.io/badge/Redis-7.2+-DC382D?style=flat&logo=redis)](https://redis.io)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

> **Enterprise-grade backend system for virtual card issuance, management, and transaction processing.**

CardFlow is a fintech platform that enables businesses and individuals to issue and manage virtual payment cards through a simple, secure API. The platform integrates with card-issuing merchant partners while providing enhanced features including KYC/KYB verification, transaction monitoring, and comprehensive card lifecycle management.

---

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Tech Stack](#-tech-stack)
- [Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Database Setup](#database-setup)
  - [Running the Application](#running-the-application)
- [API Documentation](#-api-documentation)
- [Development](#-development)
  - [Project Structure](#project-structure)
  - [Running Tests](#running-tests)
  - [Code Quality](#code-quality)
- [Deployment](#-deployment)
- [Security](#-security)
- [Contributing](#-contributing)
- [License](#-license)

---

## ✨ Features

### Core Capabilities
- 🔐 **Secure Authentication** - JWT-based auth with refresh tokens and optional MFA
- 👤 **User Management** - Comprehensive user registration and profile management
- ✅ **KYC/KYB Verification** - Automated identity verification with third-party provider integration
- 💳 **Virtual Card Issuance** - Single-use and multi-use card generation
- 📊 **Transaction Management** - Real-time transaction processing and monitoring
- 🔔 **Notifications** - Multi-channel notifications (Email, SMS, Push)
- 👨‍💼 **Admin Portal APIs** - Complete administrative control and reporting
- 📝 **Audit Logging** - Comprehensive audit trails for compliance

### Security & Compliance
- ✅ PCI-DSS compliant (no storage of full card numbers)
- ✅ GDPR compliant with data protection measures
- ✅ AML/KYC compliance
- ✅ End-to-end encryption (TLS 1.3)
- ✅ Data encryption at rest (AES-256)
- ✅ Rate limiting and DDoS protection
- ✅ RBAC (Role-Based Access Control)

### Performance & Reliability
- ⚡ Sub-500ms API response time (95th percentile)
- 🔄 Horizontal scalability
- 🛡️ Circuit breaker pattern for external APIs
- 📈 99.9% uptime SLA
- 🔁 Automated retry mechanisms
- 📊 Comprehensive monitoring and alerting

---

## 🏗️ Architecture

CardFlow follows a **clean, layered architecture** pattern:

```
┌─────────────────────────────────────────────────────────┐
│                   Client Applications                    │
│              (Web, Mobile, Third-party APIs)             │
└────────────────────────┬────────────────────────────────┘
                         │ HTTPS/REST
                         │
┌────────────────────────▼────────────────────────────────┐
│              API Gateway & Load Balancer                 │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│         Presentation Layer (HTTP Handlers)               │
├──────────────────────────────────────────────────────────┤
│            Service Layer (Business Logic)                │
├──────────────────────────────────────────────────────────┤
│         Repository Layer (Data Access)                   │
├──────────────────────────────────────────────────────────┤
│    Database Layer (PostgreSQL + Redis)                   │
└──────────────────────────────────────────────────────────┘
         │                              │
         │                              │
    ┌────▼──────┐              ┌────────▼─────────┐
    │  Merchant │              │   KYC Provider   │
    │ Partner   │              │   (korapay) │
    │    API    │              └──────────────────┘
    └───────────┘
```

### Design Patterns
- **Repository Pattern** - Data access abstraction
- **Service Layer Pattern** - Business logic encapsulation
- **Factory Pattern** - Object creation
- **Strategy Pattern** - Provider abstraction (KYC, notifications)
- **Circuit Breaker** - External service resilience
- **Observer Pattern** - Event-driven notifications

---

## 🛠️ Tech Stack

### Core Technologies
| Component | Technology | Version |
|-----------|------------|---------|
| Language | Go (Golang) | 1.21+ |
| Web Framework | Fiber | 2.51+ |
| Database | PostgreSQL | 15+ |
| Cache/Queue | Redis | 7.2+ |
| ORM | GORM | 1.25+ |
| Migration | golang-migrate | 4.16+ |

### Key Libraries
```
- github.com/gofiber/fiber/v2          # High-performance web framework
- gorm.io/gorm                         # ORM
- github.com/golang-jwt/jwt/v5         # JWT authentication
- github.com/go-redis/redis/v9         # Redis client
- golang.org/x/crypto/bcrypt           # Password hashing
- github.com/go-playground/validator   # Request validation
- github.com/joho/godotenv            # Environment configuration
- go.uber.org/zap                     # Structured logging
- github.com/stretchr/testify         # Testing framework
```

### Infrastructure
- **Containerization**: Docker
- **Orchestration**: Kubernetes
- **CI/CD**: GitHub Actions
- **Monitoring**: Prometheus + Grafana
- **Logging**: ELK Stack / CloudWatch
- **Secrets Management**: HashiCorp Vault
- **API Documentation**: Swagger/OpenAPI

---

## 🚀 Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go** 1.21 or higher ([Download](https://golang.org/dl/))
- **PostgreSQL** 15+ ([Download](https://www.postgresql.org/download/))
- **Redis** 7.2+ ([Download](https://redis.io/download))
- **Docker** (optional, for containerized development) ([Download](https://www.docker.com/get-started))
- **Make** (optional, for using Makefile commands)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/cardflow.git
   cd cardflow
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Install development tools**
   ```bash
   # Install golang-migrate for database migrations
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   
   # Install golangci-lint for code quality
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # Install swag for API documentation
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

### Configuration

1. **Create environment file**
   ```bash
   cp .env.example .env
   ```

2. **Configure environment variables**
   
   Edit `.env` with your configuration:
   ```bash
   # Application
   APP_ENV=development
   APP_PORT=8080
   APP_NAME=CardFlow
   APP_LOG_LEVEL=debug
   
   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=cardflow_dev
   DB_USER=cardflow
   DB_PASSWORD=your_secure_password
   DB_SSL_MODE=disable
   
   # Redis
   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=
   REDIS_DB=0
   
   # JWT
   JWT_SECRET=your_jwt_secret_min_32_chars_long
   JWT_ACCESS_EXPIRY=900           # 15 minutes in seconds
   JWT_REFRESH_EXPIRY=604800       # 7 days in seconds
   
   # Merchant Partner API
   MERCHANT_API_URL=https://sandbox.merchant.com
   MERCHANT_API_KEY=your_merchant_api_key
   MERCHANT_API_SECRET=your_merchant_api_secret
   MERCHANT_WEBHOOK_SECRET=your_webhook_secret
   
   # KYC Provider (korapay example)
   KYC_PROVIDER=onfido
   KYC_API_URL=https://api.korapay.com/v3
   KYC_API_TOKEN=your_korapay_api_token
   KYC_WEBHOOK_TOKEN=your_korapat_webhook_token
   
   # Email Service (SendGrid example)
   EMAIL_PROVIDER=sendgrid
   EMAIL_API_KEY=your_sendgrid_api_key
   EMAIL_FROM_ADDRESS=noreply@cardflow.com
   EMAIL_FROM_NAME=CardFlow
   
   # SMS Service (Twilio example)
   SMS_PROVIDER=twilio
   SMS_ACCOUNT_SID=your_twilio_account_sid
   SMS_AUTH_TOKEN=your_twilio_auth_token
   SMS_FROM_NUMBER=+1234567890
   
   # Encryption
   ENCRYPTION_KEY=your_32_byte_encryption_key_here
   
   # Rate Limiting
   RATE_LIMIT_ENABLED=true
   RATE_LIMIT_MAX=100
   RATE_LIMIT_WINDOW=60            # seconds
   ```

### Database Setup

1. **Create the database**
   ```bash
   createdb cardflow_db
   ```

2. **Run migrations**
   ```bash
   # Using make
   make migrate-up
   
   # Or directly
   migrate -path ./migrations \
           -database "postgresql://cardflow:password@localhost:5432/cardflow_dev?sslmode=disable" \
           up
   ```

3. **Seed the database (optional)**
   ```bash
   make seed
   ```

### Running the Application

#### Using Go directly
```bash
go run cmd/api/main.go
```

#### Using Make
```bash
make run
```

#### Using Docker Compose (Recommended for development)
```bash
# Start all services (API, PostgreSQL, Redis)
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down
```

The API will be available at `http://localhost:8080`

**Health Check**: `curl http://localhost:8080/health`

---

## 📚 API Documentation

### Interactive Documentation

Once the application is running, access the interactive API documentation:

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **ReDoc**: http://localhost:8080/redoc

### Quick API Reference

#### Authentication
```bash
# Register a new user
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "first_name": "John",
  "last_name": "Doe"
}

# Login
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}

# Response
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900
  }
}
```

#### Card Management
```bash
# Issue a new card
POST /api/v1/cards
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "card_type": "multi-use",
  "currency": "USD",
  "spending_limit": {
    "amount": 1000.00,
    "period": "monthly"
  },
  "label": "Business Expenses"
}

# List user's cards
GET /api/v1/cards
Authorization: Bearer {access_token}

# Get card details
GET /api/v1/cards/{card_id}
Authorization: Bearer {access_token}

# Freeze a card
POST /api/v1/cards/{card_id}/freeze
Authorization: Bearer {access_token}
```

#### Transactions
```bash
# List transactions
GET /api/v1/transactions?page=1&limit=50&card_id={card_id}
Authorization: Bearer {access_token}

# Get transaction details
GET /api/v1/transactions/{transaction_id}
Authorization: Bearer {access_token}
```

For complete API documentation, see swagger ui

---

## 👨‍💻 Development

### Project Structure

```
cardflow-backend/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration management
│   ├── handlers/                   # HTTP request handlers
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── kyc_handler.go
│   │   ├── card_handler.go
│   │   └── transaction_handler.go
│   ├── services/                   # Business logic layer
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── kyc_service.go
│   │   ├── card_service.go
│   │   ├── transaction_service.go
│   │   └── notification_service.go
│   ├── repositories/               # Data access layer
│   │   ├── user_repository.go
│   │   ├── card_repository.go
│   │   └── transaction_repository.go
│   ├── models/                     # Data models
│   │   ├── user.go
│   │   ├── card.go
│   │   └── transaction.go
│   ├── middleware/                 # HTTP middleware
│   │   ├── auth_middleware.go
│   │   ├── rate_limiter.go
│   │   ├── logger_middleware.go
│   │   └── cors_middleware.go
│   ├── integrations/               # External API clients
│   │   ├── merchant_client.go
│   │   ├── kyc_client.go
│   │   ├── email_client.go
│   │   └── sms_client.go
│   └── utils/                      # Utility functions
│       ├── crypto.go
│       ├── validator.go
│       └── logger.go
├── migrations/                     # Database migrations
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   └── ...
├── tests/                          # Test files
│   ├── integration/
│   ├── unit/
│   └── mocks/
├── docs/                           # Documentation
│   ├── API_DOCS.md
│   ├── DEPLOYMENT.md
│   └── SECURITY.md
├── scripts/                        # Utility scripts
│   ├── seed.go
│   └── migrate.sh
├── deployments/                    # Deployment configurations
│   ├── docker/
│   │   └── Dockerfile
│   └── kubernetes/
│       ├── deployment.yaml
│       ├── service.yaml
│       └── ingress.yaml
├── .github/
│   └── workflows/
│       └── ci-cd.yml
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
├── .env.example
├── .gitignore
└── README.md
```

### Running Tests

#### Unit Tests
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detection
go test -race ./...

# Run specific package tests
go test ./internal/services/...
```

#### Integration Tests
```bash
# Run integration tests (requires running services)
make test-integration
```

#### Test Coverage Report
```bash
# Generate HTML coverage report
make coverage-html

# View in browser
open coverage.html
```

### Code Quality

#### Linting
```bash
# Run linter
make lint

# Auto-fix issues where possible
make lint-fix
```

#### Code Formatting
```bash
# Format code
make fmt

# Check formatting
make fmt-check
```

#### Generate API Documentation
```bash
# Generate Swagger docs
make swagger

# This creates/updates docs/swagger.yaml and docs/swagger.json
```

### Database Migrations

#### Create a new migration
```bash
make migrate-create name=add_new_table
```

#### Apply migrations
```bash
# Apply all pending migrations
make migrate-up

# Apply specific number of migrations
migrate -path ./migrations -database "${DATABASE_URL}" up 2

# Rollback last migration
make migrate-down

# Reset database (careful in production!)
make migrate-reset
```

#### Check migration status
```bash
make migrate-status
```

### Useful Make Commands

```bash
make help              # Show all available commands
make build             # Build the application
make run               # Run the application
make test              # Run tests
make test-coverage     # Run tests with coverage
make lint              # Run linter
make fmt               # Format code
make docker-build      # Build Docker image
make docker-run        # Run Docker container
make clean             # Clean build artifacts
```

---

## 🚢 Deployment

### Docker Deployment

#### Build Docker image
```bash
docker build -t cardflow-api:latest -f deployments/docker/Dockerfile .
```

#### Run with Docker Compose
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### Kubernetes Deployment

#### Prerequisites
- Kubernetes cluster (1.25+)
- kubectl configured
- Helm (optional)

#### Deploy to Kubernetes
```bash
# Create namespace
kubectl create namespace cardflow

# Apply configurations
kubectl apply -f deployments/kubernetes/

# Check deployment status
kubectl get pods -n cardflow
kubectl get services -n cardflow

# View logs
kubectl logs -f deployment/cardflow-api -n cardflow
```

#### Using Helm (if available)
```bash
helm install cardflow ./deployments/helm/cardflow \
  --namespace cardflow \
  --create-namespace \
  --values values.prod.yaml
```

### Environment-Specific Deployments

```bash
# Development
make deploy-dev

# Staging
make deploy-staging

# Production (requires manual approval)
make deploy-prod
```


---

## 🔒 Security

### Security Best Practices

- ✅ **Never commit secrets** to version control
- ✅ **Use environment variables** for sensitive configuration
- ✅ **Rotate credentials regularly** (every 90 days minimum)
- ✅ **Enable MFA** for admin accounts
- ✅ **Keep dependencies updated** - run `make update-deps` regularly
- ✅ **Run security scans** - `make security-scan`
- ✅ **Review audit logs** regularly

### Reporting Security Issues

If you discover a security vulnerability, please email ojeniwehalexander@gmail.com immediately. Do **not** create a public GitHub issue.

### Security Checklist

- [ ] All environment variables properly set
- [ ] TLS/SSL certificates configured
- [ ] Firewall rules properly configured
- [ ] Database access restricted to application only
- [ ] API rate limiting enabled
- [ ] Monitoring and alerting configured
- [ ] Regular security audits scheduled
- [ ] Incident response plan documented


---

## 🤝 Contributing

We welcome contributions! Please follow these guidelines:

### Development Workflow

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Make your changes**
4. **Write/update tests**
5. **Ensure tests pass**
   ```bash
   make test
   make lint
   ```
6. **Commit your changes**
   ```bash
   git commit -m "feat: add amazing feature"
   ```
7. **Push to your fork**
   ```bash
   git push origin feature/amazing-feature
   ```
8. **Open a Pull Request**

### Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new feature
fix: bug fix
docs: documentation changes
style: formatting, missing semi colons, etc
refactor: code restructuring
test: adding tests
chore: updating build tasks, package manager configs, etc
```

### Code Review Process

- All submissions require code review
- CI/CD checks must pass
- Test coverage must not decrease
- Documentation must be updated

### Pull Request Checklist

- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] No linting errors
- [ ] All tests passing
- [ ] Commit messages follow convention

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👥 Team & Support

### Core Team
- **Project Lead**: [Alexander Ojeniweh](mailto:ojeniwehalexander@gmail.com)
- **Backend Lead**: [Alexander Ojeniweh](mailto:ojeniwehalexander@gmail.com)

### Support Channels
- 📧 **Email**: ojeniwehalexander@gmail.com
- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/Ojeniweh10/CardFlow/issues)

### Acknowledgments

Special thanks to:
- The Fiber framework team
- The GORM maintainers
- All our contributors

---

## 📊 Project Status

### Current Version: `v1.0.0`

### Roadmap

#### ✅ Phase 1: Foundation (Completed)
- [x] Core authentication system
- [x] User management
- [x] Database setup and migrations
- [x] Basic CI/CD pipeline

#### 🚧 Phase 2: Core Features (In Progress)
- [x] KYC integration
- [x] Card issuance
- [ ] Transaction processing
- [ ] Notification service

#### 📅 Phase 3: Enhancement (Planned Q2 2026)
- [ ] Multi-currency support
- [ ] Advanced fraud detection
- [ ] Real-time analytics dashboard
- [ ] Mobile SDK

#### 📅 Phase 4: Scale (Planned Q3 2026)
- [ ] Multi-region deployment
- [ ] Advanced caching strategies
- [ ] Performance optimizations
- [ ] GraphQL API

---

## 📈 Metrics

![GitHub stars](https://img.shields.io/github/stars/Ojeniweh10/CardFlow?style=social)
![GitHub forks](https://img.shields.io/github/forks/Ojeniweh10/CardFlow?style=social)
![GitHub issues](https://img.shields.io/github/issues/Ojeniweh10/CardFlow)
![GitHub pull requests](https://img.shields.io/github/issues-pr/Ojeniweh10/CardFlow)

---

<div align="center">


Made with ❤️ by the CardFlow Team

</div>
