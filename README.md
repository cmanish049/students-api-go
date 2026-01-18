# Students API

A RESTful API service for managing student records built with Go and SQLite. This project demonstrates clean architecture principles with clear separation of concerns between HTTP handlers, business logic, and data storage.

## Features

- ✅ Create, Read, Update, and Delete (CRUD) operations for students
- ✅ SQLite database for data persistence
- ✅ Request validation using validator/v10
- ✅ Structured JSON responses
- ✅ YAML-based configuration
- ✅ Graceful server shutdown
- ✅ Structured logging with slog
- ✅ Clean architecture with dependency injection

## Tech Stack

- **Language**: Go 1.25.5
- **Database**: SQLite
- **HTTP Router**: Go standard library (http.ServeMux)
- **Configuration**: cleanenv (YAML/ENV)
- **Validation**: go-playground/validator
- **Logging**: Go standard library (log/slog)

## Project Structure

```
students-api/
├── cmd/
│   └── students-api/
│       └── main.go              # Application entry point
├── config/
│   └── local.yaml               # Local configuration file
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loading logic
│   ├── http/
│   │   └── handlers/
│   │       └── student/
│   │           └── student.go   # HTTP handlers for student operations
│   ├── storage/
│   │   ├── storage.go           # Storage interface definition
│   │   ├── postgres/            # PostgreSQL implementation (placeholder)
│   │   └── sqlite/
│   │       └── sqlite.go        # SQLite implementation
│   ├── types/
│   │   └── types.go             # Data type definitions
│   └── utils/
│       └── response/
│           └── response.go      # HTTP response utilities
├── storage/                     # Database file location
├── go.mod                       # Go module dependencies
└── README.md                    # This file
```

## Installation

### Prerequisites

- Go 1.25.5 or higher
- SQLite3

### Setup

1. Clone the repository:
```bash
git clone https://github.com/cmanish049/students-api.git
cd students-api
```

2. Install dependencies:
```bash
go mod download
```

3. Create the storage directory:
```bash
mkdir -p storage
```

## Configuration

The application uses YAML configuration files. Create or modify `config/local.yaml`:

```yaml
env: "dev"
storage_path: "storage/storage.db"
http_server:
  address: "localhost:8082"
```

### Configuration Options

- `env`: Environment name (dev, production)
- `storage_path`: Path to SQLite database file
- `http_server.address`: Server address and port

### Configuration Loading

The application loads configuration in the following priority:

1. `CONFIG_PATH` environment variable
2. `--config` command-line flag

Example:
```bash
# Using environment variable
export CONFIG_PATH=config/local.yaml
go run cmd/students-api/main.go

# Using command-line flag
go run cmd/students-api/main.go --config=config/local.yaml
```

## Running the Application

### Development Mode

```bash
go run cmd/students-api/main.go --config=config/local.yaml
```

### Build and Run

```bash
# Build the binary
go build -o bin/students-api cmd/students-api/main.go

# Run the binary
./bin/students-api --config=config/local.yaml
```

The server will start on `http://localhost:8082` (or the address specified in your config file).

## API Endpoints

### Student Model

```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "age": 20
}
```

### Endpoints

#### Create a Student

```http
POST /api/students
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 20
}
```

**Success Response** (201 Created):
```json
{
  "id": 1
}
```

**Error Response** (400 Bad Request):
```json
{
  "status": "Error",
  "error": "field Name is required field"
}
```

#### Get Student by ID

```http
GET /api/students/{id}
```

**Success Response** (200 OK):
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "age": 20
}
```

#### Get All Students

```http
GET /api/students
```

**Success Response** (200 OK):
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 20
  },
  {
    "id": 2,
    "name": "Jane Smith",
    "email": "jane@example.com",
    "age": 22
  }
]
```

#### Update a Student

```http
PUT /api/students/{id}
Content-Type: application/json

{
  "name": "John Updated",
  "email": "john.updated@example.com",
  "age": 21
}
```

**Success Response** (200 OK):
```json
{
  "message": "student updated successfully"
}
```

#### Delete a Student

```http
DELETE /api/students/{id}
```

**Success Response** (200 OK):
```json
{
  "message": "student deleted successfully"
}
```

## Testing with cURL

### Create a student
```bash
curl -X POST http://localhost:8082/api/students \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","age":20}'
```

### Get all students
```bash
curl http://localhost:8082/api/students
```

### Get student by ID
```bash
curl http://localhost:8082/api/students/1
```

### Update a student
```bash
curl -X PUT http://localhost:8082/api/students/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Updated","email":"john.updated@example.com","age":21}'
```

### Delete a student
```bash
curl -X DELETE http://localhost:8082/api/students/1
```

## Validation Rules

The API validates student data with the following rules:

- `name`: Required field
- `email`: Required field (must be unique in database)
- `age`: Required field

## Error Handling

The API returns consistent error responses:

```json
{
  "status": "Error",
  "error": "error description"
}
```

Common HTTP status codes:
- `200 OK`: Successful GET/PUT/DELETE operation
- `201 Created`: Successful POST operation
- `400 Bad Request`: Invalid input or validation error
- `500 Internal Server Error`: Server-side error

## Database Schema

The SQLite database contains a single `students` table:

```sql
CREATE TABLE IF NOT EXISTS students (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    age INTEGER NOT NULL
);
```

## Architecture

### Clean Architecture Principles

The project follows clean architecture patterns:

1. **Handlers Layer** (`internal/http/handlers`): HTTP request/response handling
2. **Storage Interface** (`internal/storage`): Abstraction for data persistence
3. **Storage Implementation** (`internal/storage/sqlite`): Concrete database implementation
4. **Types** (`internal/types`): Domain models
5. **Config** (`internal/config`): Configuration management
6. **Utils** (`internal/utils`): Shared utilities

### Dependency Injection

The application uses dependency injection to maintain loose coupling:

```go
// Storage interface is injected into handlers
router.HandleFunc("POST /api/students", student.New(db))
```

This allows for easy testing and swapping of storage implementations (e.g., SQLite to PostgreSQL).

## Graceful Shutdown

The application implements graceful shutdown to handle SIGINT and SIGTERM signals:

- Active requests complete processing
- Server shutdown timeout: 5 seconds
- Database connections are properly closed

## Development

### Adding New Features

1. **Add new storage method**: Update `internal/storage/storage.go` interface
2. **Implement in SQLite**: Add method to `internal/storage/sqlite/sqlite.go`
3. **Create handler**: Add handler in `internal/http/handlers/student/student.go`
4. **Register route**: Add route in `cmd/students-api/main.go`

### Adding PostgreSQL Support

The project structure includes a `postgres/` directory for future PostgreSQL implementation:

1. Implement the `Storage` interface in `internal/storage/postgres/`
2. Update the main.go to switch between SQLite and PostgreSQL based on configuration

## Dependencies

- `github.com/go-playground/validator/v10`: Request validation
- `github.com/ilyakaznacheev/cleanenv`: Configuration management
- `github.com/mattn/go-sqlite3`: SQLite driver
- Go standard library for HTTP server and logging

## License

This project is available for educational and personal use.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Author

[@cmanish049](https://github.com/cmanish049)

## Roadmap

- [ ] Add unit tests
- [ ] Add integration tests
- [ ] Implement PostgreSQL support
- [ ] Add authentication and authorization
- [ ] Add pagination for GET /api/students
- [ ] Add filtering and sorting
- [ ] Add API documentation with Swagger
- [ ] Add Docker support
- [ ] Add CI/CD pipeline
