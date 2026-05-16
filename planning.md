- Clean Architecture
- Domain Driven Design (DDD)

```bash
.
├── api/                            # API contracts and specifications
│   ├── openapi/                    # OpenAPI / Swagger files
│   ├── proto/                      # gRPC protobuf definitions
│   └── jsonschema/                 # JSON schema definitions
│
├── build/                          # Packaging and deployment files
│   ├── docker/                     # Docker related files
│   ├── kubernetes/                 # Kubernetes manifests
│   ├── ci/                         # CI/CD pipelines
│   └── scripts/                    # Build scripts
│
├── cmd/                            # Application entry points
│   └── app/
│       └── main.go                 # Main application bootstrap
│
├── docs/                           # Project documentation
│   ├── architecture/
│   ├── api/
│   └── adr/
│
├── domain/                         # Domain layer (pure business logic)
│   ├── aggregate/                  # Aggregate roots
│   ├── entity/                     # Domain entities
│   ├── valueobject/                # Immutable value objects
│   ├── enum/                       # Enums/constants
│   ├── event/                      # Domain events
│   ├── repository/                 # Repository interfaces
│   └── service/                    # Domain services
│
├── application/                    # Application layer
│   ├── usecase/                    # Business use cases
│   └── repository/                 # Application repository contracts

│
├── infrastructure/                 # Infrastructure implementations
│   ├── persistence/                # Database implementations
│   │   ├── mysql/
│   │   ├── postgres/
│   │   └── redis/
│   │
│   ├── messaging/                  # Kafka/RabbitMQ implementations
│   ├── service/                    # External service implementations
│   ├── config/                     # Config loader
│   ├── logger/                     # Logging implementation
│   └── observability/              # Metrics/tracing
│
├── delivery/                       # Transport layer
│   ├── http/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── request/
│   │   └── response/
│   │
│   ├── grpc/
│   │   ├── handler/
│   │   └── interceptor/
│   │
│   └── router/
│
├── migrations/                     # Database migration files
│
├── pkg/                            # Shared reusable packages
│   ├── errorcodes/
│   ├── middleware/
│   ├── utils/
│   ├── response/
│   └── constants/
│
├── test/                           # Integration and E2E tests
│   ├── integration/
│   ├── e2e/
│   └── testdata/
│
├── .env.example                    # Environment example
├── Makefile                        # Build automation
├── Dockerfile                      # Docker image definition
├── docker-compose.yml              # Local development environment
├── go.mod
├── go.sum
└── README.md
```