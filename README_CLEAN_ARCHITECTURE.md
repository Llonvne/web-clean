# Clean Architecture Implementation

This project implements a complete Clean Architecture pattern following Uncle Bob's principles for building maintainable, testable, and scalable applications.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Clean Architecture                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────── │
│  │   Interface     │    │   Application   │    │     Domain     │
│  │    Layer        │ -> │     Layer       │ -> │     Layer      │
│  │  (HTTP/REST)    │    │   (Use Cases)   │    │  (Entities)    │
│  └─────────────────┘    └─────────────────┘    └─────────────── │
│           │                       │                       │     │
│           v                       v                       │     │
│  ┌─────────────────┐                                     │     │
│  │ Infrastructure  │ <-----------------------------------┘     │
│  │     Layer       │        (Dependency Inversion)             │
│  │ (Database/Web)  │                                           │
│  └─────────────────┘                                           │
└─────────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
├── cmd/                           # Application entry points
│   └── main.go                   # Dependency injection & wiring
├── domain/                       # Shared domain interfaces
│   └── log.go                    # Logging interface
├── internal/                     # Internal application code
│   ├── domain/                   # Domain Layer (Core Business Logic)
│   │   ├── entity/               # Business entities
│   │   │   └── user.go          # User domain entity
│   │   ├── repository/           # Repository interfaces
│   │   │   └── user_repository.go # User repository contract
│   │   └── usecase/              # Use case interfaces
│   │       └── user_usecase.go   # User business operations
│   ├── application/              # Application Layer (Business Logic)
│   │   └── service/              # Use case implementations
│   │       └── user_service.go   # User business logic
│   ├── infrastructure/           # Infrastructure Layer (External Concerns)
│   │   └── repository/           # Repository implementations
│   │       └── user_repository_impl.go # Database implementation
│   └── interface/                # Interface Layer (Delivery Mechanisms)
│       └── http/                 # HTTP handlers
│           └── user_handler.go   # REST API handlers
├── infra/                        # Infrastructure framework
└── repository/                   # Legacy repository (to be migrated)
```

## Core Principles Implemented

### 1. Dependency Inversion Principle ✅
- **Domain layer** defines interfaces (contracts)
- **Infrastructure layer** implements these interfaces
- Dependencies point inward toward the domain
- No circular dependencies

### 2. Separation of Concerns ✅
- **Domain Layer**: Business rules and entities
- **Application Layer**: Use cases and business logic
- **Infrastructure Layer**: Database, external services
- **Interface Layer**: HTTP, CLI, etc.

### 3. Independent Layers ✅
- Each layer has a single responsibility
- Inner layers don't depend on outer layers
- Easy to swap implementations
- Framework independent core business logic

### 4. Testability ✅
- Business logic is isolated from infrastructure
- Dependencies are injected via interfaces
- Easy to mock external dependencies
- Unit tests can run without database/web server

## Layer Responsibilities

### Domain Layer (`internal/domain/`)
- **Purpose**: Core business logic and rules
- **Contains**: Entities, Value Objects, Domain Services, Repository Interfaces
- **Dependencies**: None (innermost layer)
- **Example**: User entity with business validation

```go
// User entity with business rules
type User struct {
    ID        uuid.UUID
    Email     string
    Username  string
    Name      string
}

func (u *User) IsValid() bool {
    return u.ID != uuid.Nil && u.Email != "" && u.Username != ""
}
```

### Application Layer (`internal/application/`)
- **Purpose**: Orchestrates business logic and use cases
- **Contains**: Services implementing use case interfaces
- **Dependencies**: Domain layer only
- **Example**: User service with business operations

```go
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*entity.User, error) {
    // Business rule: Check if user already exists
    existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
    if existing != nil {
        return nil, ErrUserAlreadyExists
    }
    
    // Create and validate entity
    user := entity.NewUser(req.Email, req.Username, req.Name)
    if !user.IsValid() {
        return nil, ErrInvalidUserData
    }
    
    return s.userRepo.Create(ctx, user)
}
```

### Infrastructure Layer (`internal/infrastructure/`)
- **Purpose**: Implementation of external concerns
- **Contains**: Repository implementations, Database models, External service clients
- **Dependencies**: Domain interfaces
- **Example**: PostgreSQL user repository implementation

```go
type UserRepositoryImpl struct {
    db database.Database
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
    model := &UserModel{}
    model.FromEntity(user)
    return r.db.Transaction(func(tx *gorm.DB) error {
        return tx.WithContext(ctx).Create(model).Error
    })
}
```

### Interface Layer (`internal/interface/`)
- **Purpose**: Handle external communication protocols
- **Contains**: HTTP handlers, CLI commands, gRPC servers
- **Dependencies**: Use case interfaces
- **Example**: REST API handlers

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid_request"})
        return
    }
    
    user, err := h.userUseCase.CreateUser(c.Request.Context(), useCaseReq)
    if err != nil {
        h.handleError(c, err)
        return
    }
    
    c.JSON(http.StatusCreated, toUserResponse(user))
}
```

## API Endpoints

The following RESTful endpoints are available:

```http
# Health Check
GET /health

# API Documentation
GET /api/v1/

# User Management
POST   /api/v1/users           # Create user
GET    /api/v1/users           # List users (paginated)
GET    /api/v1/users/:id       # Get user by ID
PUT    /api/v1/users/:id       # Update user profile
DELETE /api/v1/users/:id       # Delete user
```

## Benefits of This Architecture

### 🔧 Maintainability
- Clear separation of concerns
- Easy to locate and modify specific functionality
- Consistent patterns throughout the codebase

### 🧪 Testability
- Business logic isolated from infrastructure
- Easy to unit test without external dependencies
- Mockable interfaces for integration testing

### 🔄 Flexibility
- Easy to swap database implementations
- Can add new delivery mechanisms (gRPC, CLI, etc.)
- Framework-independent business logic

### 📈 Scalability
- Clear boundaries enable team scaling
- Independent deployment of layers possible
- Microservices-ready architecture

### 🔒 Reliability
- Dependency inversion prevents tight coupling
- Error handling isolated to appropriate layers
- Fail-fast principle with proper error propagation

## Running the Application

```bash
# Install dependencies
go mod download

# Run the application
go run cmd/main.go

# The server will start on the configured port
# Health check: GET http://localhost:8080/health
# API documentation: GET http://localhost:8080/api/v1/
```

## Testing Strategy

```bash
# Unit tests (business logic)
go test ./internal/application/service/...

# Integration tests (repository layer)
go test ./internal/infrastructure/repository/...

# API tests (interface layer)
go test ./internal/interface/http/...
```

## Migration from Legacy Code

The existing code has been restructured to follow Clean Architecture:

1. **Domain entities** extracted from scattered models
2. **Repository interfaces** defined in domain layer
3. **Business logic** centralized in application services
4. **HTTP concerns** isolated in interface layer
5. **Dependency injection** properly configured in main.go

This ensures the codebase follows Clean Architecture principles completely while maintaining backward compatibility where needed.

## Next Steps

1. Add comprehensive unit tests for all layers
2. Implement additional domain entities (orders, products, etc.)
3. Add integration tests with test database
4. Implement authentication and authorization
5. Add API documentation with Swagger/OpenAPI
6. Migrate remaining legacy code to Clean Architecture

---

**Note**: This implementation demonstrates a complete adherence to Clean Architecture principles as defined by Robert C. Martin, ensuring maintainable, testable, and scalable code.