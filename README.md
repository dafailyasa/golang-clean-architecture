# Auth Service - Hotel Microservice

Auth Service is a core microservice responsible for handling user authentication, authorization, and identity management within the broader Go Hotel Microservices ecosystem. 

Designed following **Clean Architecture** and **Domain-Driven Design (DDD)** principles, this service abstracts infrastructure dependencies behind application ports, granting superior reliability and isolated testing capabilities.

## 🏗 Architecture 

This service enforces the **Clean Architecture** structure:
- **Domain Layer (`internal/domain`)**: Core business entities (`User`), Value Objects (`Email`), and domain policies.
- **Application Layer (`internal/application`)**: UseCases encapsulating application-specific business rules.
- **Infrastructure Layer (`internal/infrastructure`)**: Outermost layer dictating how data persists (MySQL with GORM) and identity provider mappings (Keycloak).
- **Presentation Layer (`internal/presentation`)**: HTTP interface resolving delivery mechanisms, JSON parsing, API routing via Gorilla Mux, and middleware chaining.

## 🚀 Tech Stack

- **Go 1.2x**
- **HTTP/Router**: [Gorilla Mux](https://github.com/gorilla/mux)
- **Database**: MySQL with [GORM](https://gorm.io/)
- **Identity Provider**: Keycloak API Integration
- **Validations**: [go-playground/validator](https://github.com/go-playground/validator)
- **Testing**: Testify Mocks

## 🛠 Features

- Generic **Pagination Module** integrating decoupled scopes for scalable data-table integration (e.g., Search, Sort, Limit, Page).
- User Management: 
  - Register new users
  - Update user metadata
  - Retrieve details by Unique Keycloak UUID or system ID
  - List / Search users dynamically
  - Delete User
- Native IDP Session Handling: Login / Renew Access Tokens.
- Graceful HTTP Shutdown and Robust Error Standardizations mapping Infrastructure and Business errors securely via internal configurations (`pkg/errors`, `pkg/response`).

## 📁 Project Structure

```text
auth-service/
├── cmd/                            # Main application entry point
├── config/                         # Configuration loading & environmental binding
├── internal/                       # Private application & domain code
│   ├── application/                # Application UseCases & Port Definitions
│   │   ├── dto/                    # Data Transfer Objects
│   │   ├── port/                   # Interface abstraction (DB/IdentityProvider)
│   │   ├── service/                # Logic adapting application & domain
│   │   └── usecase/                # System boundaries defining explicit features
│   ├── domain/                     # Core Business Logic & Enterprise patterns
│   ├── infrastructure/             # Concrete Implementations mapped to ports
│   │   ├── keycloack/              # Keycloak Identity Provider integration
│   │   └── persistence/mysql/      # Data schema bindings & migrations
│   └── presentation/               # HTTP Interface logic mapping JSON/REST
│       ├── http/handler/           # Delivery implementation (Controllers)
│       ├── http/middlewares/       # Intercepts mappings (CORS, Validations)
│       └── http/router/            # Aggressive URI bounding & Prefix grouping
└── pkg/                            # Shared utilities and helpers logic
    ├── app_http/                   # Internal standard Resty/HTTP wrappers
    ├── errors/                     # Standalone structured error abstractions
    ├── pagination/                 # Generic isolated HTTP pagination handling
    └── response/                   # REST Standardized JSON response formatting
```

## ⚙️ Setup and Installation

**1. Clone and navigate to the Auth Service**
```bash
cd auth-service
```

**2. Copy Configuration**
Ensure you construct a `config.yaml` file natively at the service root containing credentials resolving:
- Target Database Configuration (MySQL DSN setup)
- Environment and Logging definitions
- Keycloak API target and Master Client contexts

**3. Run the microservice**
```bash
go run cmd/main.go
```

## 📍 API Reference (V1)

| Method | Endpoint                    | Description                          |
|--------|-----------------------------|--------------------------------------|
| `POST`   | `/api/v1/users`             | Register / Create a user account     |
| `GET`    | `/api/v1/users`             | Search and paginates system accounts |
| `GET`    | `/api/v1/users/{id}`        | Read user profile mappings           |
| `PUT`    | `/api/v1/users/{id}`        | Update existing details natively     |
| `DELETE` | `/api/v1/users/{id}`        | Purge record off database entirely   |
| `POST`   | `/api/v1/users/login`       | Authenticate granting Access Token   |
| `POST`   | `/api/v1/users/refresh-token`| Yield fresh keycloak Session cycle  |

*Note: Listing endpoint accepts query params `?page=1&limit=10&search=xyz&sortBy=created_at&sortOrder=desc`*
