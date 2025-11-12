# Backend Engineering Interview Assignment (Golang)

## Author

**Name:** Mochamad Sohibul Iman  
**Email:** [iman@imansohibul.my.id](mailto:iman@imansohibul.my.id)  
**LinkedIn:** [www.linkedin.com/in/imansohibul](https://www.linkedin.com/in/imansohibul)

## 🏗️ Project Structure & Software Architecture

This project follows a **modular clean architecture** pattern. It ensures high maintainability, testability, and clear separation of concerns.

### 🧱 Architectural Layers

| Layer        | Responsibility |
|--------------|----------------|
| **Entity**   | Core domain logic: models and business rules. No framework or external dependency here. |
| **Usecase**  | Orchestrates application flow: how data moves and is transformed. Calls repositories and domain logic. |
| **Repository** | Data persistence and third-party integration. Implements storage logic (PostgreSQL, etc.). |
| **Delivery (REST)** | Handles HTTP requests and responses using Echo. Maps JSON ↔ DTO ↔ Entities. |
| **Config**   | Dependency wiring (DI), configuration loading, and server setup. |
| **DB Migrate** | Database version control using SQL migrations. |

---

## 📂 Project Structure
```text
.
📦 otp-services
├── cmd/                     # Application entrypoints
│   └── main.go              # Main function as entrypoint for REST API, consumer, cron-job, etc
├── config/                  # Configuration management and dependency injection
│   ├── common.go            # Common configuration
│   └── server.go            # Server configuration
├── db/
│   └── migrate/             # DB migrations using golang-migrate (up/down SQL files)
│       ├── 20251111124517_create_otps_table.down.sql
│       └── 20251111124517_create_otps_table.up.sql
├── entity/                  # Domain entities and business rules
│   ├── error_test.go        # Error entity tests
│   ├── error.go             # Error entity definitions
│   └── otp.go               # OTP entity
├── generated/
│   └── api.gen.go           # Generated API code
├── internal/
│   ├── handler/             # HTTP handlers (controllers)
│   │   ├── middleware/      # Custom middleware
│   │   ├── mock/            # Handler mocks for testing
│   │   ├── otp_test.go      # OTP handler tests
│   │   ├── otp.go           # OTP handler
│   │   ├── server.go        # Server setup, routing, and middleware
│   │   └── usecase.go       # Use case interfaces
│   ├── repository/          # Data access layer (Postgres, etc.)
│   │   ├── otp_repository_test.go
│   │   ├── otp_repository.go
│   │   ├── repository_test.go
│   │   ├── repository.go    # Repository implementation
│   │   ├── transaction_manager_test.go
│   │   ├── transaction_manager.go
│   │   └── types.go         # Repository types
│   └── usecase/             # Application use cases (interactors)
│       ├── mock/            # Use case mocks for testing
│       ├── otp_generator_test.go
│       ├── otp_generator.go # OTP generation logic
│       ├── otp_test.go
│       ├── otp.go           # OTP use case
│       └── repository.go    # Repository interfaces
├── .env                     # Environment configuration
├── .gitignore               # Git ignore file
├── api.yml                  # API specification (OpenAPI/Swagger)
├── coverage.out             # Test coverage output
├── docker-compose.yml       # Defines services (DB) for development
├── env.sample               # Sample environment configuration
├── go.mod                   # Go module dependencies
├── go.sum                   # Go module checksums
├── Makefile                 # Common scripts for building, testing, running
└── README.md                # Project documentation
```

## 🚀 Project Setup

### Prerequisites
- Go 1.25 or higher
- MySQL 8.0 or higher
- Make

### 1. Clone the Repository
```bash
git clone <repository-url>
cd otp-services
```

### 2. Environment Configuration
Copy the sample environment file and configure it:
```bash
cp env.sample .env
```

Update `.env` with your configuration:
```env
SERVICE_DB_USERNAME=mysqldev
SERVICE_DB_PASSWORD=mysqldev
SERVICE_DB_HOST=127.0.0.1
SERVICE_DB_PORT=3306
SERVICE_DB_NAME=otp-service-dev
```

### 3. Install Dependencies
```bash
make init
```

This command will:
- Clean generated files
- Generate API code from `api.yml`
- Generate mocks for testing
- Download Go dependencies

### 4. Database Setup
Create the database:
```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS \`otp-service-dev\`;"
```

Run migrations:
```bash
make migrate MIGRATE_ARGS=up
```

### 5. Run the Application
```bash
go run cmd/main.go
```

Or build and run:
```bash
make all
./build/main
```

### 6. Run Tests
```bash
make test
```

This will run all tests and generate a coverage report.

## 📝 Available Make Commands

| Command | Description |
|---------|-------------|
| `make init` | Initialize project (clean, generate code, install dependencies) |
| `make build/main` | Build the application binary |
| `make clean` | Remove generated files |
| `make generate` | Generate API code and mocks |
| `make test` | Run tests with coverage |
| `make migrate MIGRATE_ARGS=up` | Run database migrations up |
| `make migrate MIGRATE_ARGS=down` | Run database migrations down |
| `make migrate MIGRATE_ARGS=down N=1` | Rollback last migration |
| `make create-db-migration MIGRATE_NAME=<name>` | Create new migration files |

## 🗄️ Database Migration Examples

### Create a new migration
```bash
make create-db-migration MIGRATE_NAME=add_users_table
```

### Run all migrations
```bash
make migrate MIGRATE_ARGS=up
```

### Rollback last migration
```bash
make migrate MIGRATE_ARGS=down N=1
```

### Rollback all migrations
```bash
make migrate MIGRATE_ARGS=down
```

## 🐳 Docker Setup (Optional)

Run with Docker Compose:
```bash
docker-compose up -d
```

This will start:
- MySQL Database

## 🔧 Development Workflow

1. Make changes to your code
2. If you modify `api.yml`, run `make generate` to regenerate API code
3. Run tests: `make test`
4. Build: `make build/main`
5. Run: `./build/main`

## 📦 Project Dependencies

The project uses:
- **oapi-codegen**: OpenAPI code generation
- **mockgen**: Mock generation for testing
- **golang-migrate**: Database migrations

These tools are automatically installed when running `make init`.
