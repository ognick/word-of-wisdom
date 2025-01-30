# Word of Wisdom

## Description

The **Word of Wisdom** project provides clients with inspiring quotes through TCP and HTTP servers, protected by a Proof of Work (PoW) mechanism to prevent DDoS attacks.  It utilizes the crypto/argon2 package to generate and verify the "challenges".

## Project Architecture

- **Clean Architecture**: A clear separation of layers: domain, application, and infrastructure.
- **Domain-Driven Design (DDD)**: Focus on modeling domain entities and business logic.

## Monorepository Structure

- **arch/**: Architecture-related documentation and diagrams.
- **cmd/**: Entry point for applications.
- **deploy/**: Deployment-related configurations and scripts.
- **internal/**: Internal packages not intended for external use.
  - **common/**: Common utilities shared across the project.
    - **config/**: Shared configuration utilities.
    - **constants/**: Global constants.
  - **server/**
    - **app/**: Core application logic.
        - **docs/**: Documentation related to the server application.
    - **internal/**: Internal utilities for the server.
        - **config/**: Configuration handling.
        - **domain/**: Business logic layer.
            - **models/**: Data models.
            - **interfaces/**: Interfaces for dependencies.
                - **usecases/**: Use case interfaces.
        - **services/**: Service layer implementation.
            - **wisdom/**: Wisdom-related services.
                - **repository/**: Data persistence layer.
                - **api/**: API integrations.
                    - **tcp/**: TCP-based API.
                        - **v1/**: Version 1.
                    - **http/**: HTTP-based API.
                        - **v1/**: Version 1.
                            - **dto/**: Data transfer objects.
                - **usecase/**: Wisdom-related use case implementations.
- **pkg/**: Shared reusable packages.
  - **pow/**: Proof-of-work related utilities.
  - **logger/**: Logging utilities.
  - **shutdown/**: Graceful shutdown handling.
  - **lifecycle/**: Application lifecycle management.


## Technologies and Approaches

### Core Technologies
- **Programming Language**: Go
- **Networking Protocols**: TCP and HTTP
- **Security**: Proof of Work (PoW) for DDoS attack prevention
- **Containerization**: Docker and Docker Compose for environment management
- **CI/CD**: GitHub Actions for automated build and deployment processes
- **Testing**: Unit tests and static code analysis using golangci-lint
- **Dependency Management**: Wire for automatic dependency injection
- **Architecture Analysis Tools**: arch-go for architecture validation and dependency analysis

## How to Run

Run test and benchmarks:
```
make test
```

Run demonstration:
```
make run
```
