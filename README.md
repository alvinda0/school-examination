# MY-APP

REST API menggunakan Go dengan Clean Architecture.

## 📁 Struktur Proyek

```
.
├── .github/              # GitHub configuration
├── cmd/
│   └── server/           # Application entry point (main.go)
├── internal/             # Private application code
│   ├── api/              # DTOs & Response models
│   ├── appern/           # Application errors
│   ├── config/           # Configuration management
│   ├── database/         # Database connections
│   ├── handlers/         # HTTP handlers (controllers)
│   ├── mapper/           # Data mappers
│   ├── middleware/       # HTTP middleware
│   ├── model/            # Domain models
│   ├── repository/       # Data access layer
│   ├── routes/           # Route definitions
│   ├── services/         # Business logic
│   └── utils/            # Utilities
├── migrations/           # Database migrations
├── pkg/                  # Public libraries
├── test/                 # Tests
├── .env                  # Environment variables
├── .env.example          # Environment variables template
├── go.mod                # Go module definition
└── README.md             # This file
```

## 🏗️ Arsitektur

Proyek ini menggunakan **Clean Architecture** dengan layer:

1. **Handler Layer** (`internal/handlers/`) - Menerima HTTP request
2. **Service Layer** (`internal/services/`) - Business logic
3. **Repository Layer** (`internal/repository/`) - Data access
4. **Model Layer** (`internal/model/`) - Domain entities

## 🚀 Getting Started

### Prerequisites

- Go 1.26+
- PostgreSQL 12+

### Installation

1. Clone repository:
```bash
git clone <repository-url>
cd my-app
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Edit `.env` sesuai konfigurasi database Anda

4. Install dependencies:
```bash
go mod download
```

5. Jalankan migration (manual):
```bash
psql -h localhost -p 5433 -U postgres -d postgres -f migrations/000001_create_users.sql
```

6. Jalankan aplikasi:
```bash
go run cmd/server/main.go
```

Server akan berjalan di `http://localhost:8080`

## 📡 API Endpoints

### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users` | Get all users |
| GET | `/api/v1/users?role_id=1` | Get users by role |
| GET | `/api/v1/users/{id}` | Get user by ID |
| POST | `/api/v1/users` | Create new user |
| PUT | `/api/v1/users/{id}` | Update user |
| DELETE | `/api/v1/users/{id}` | Delete user |

### Roles

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/role` | Get all roles |

## 📝 Request Examples

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "age": 25}'
```

### Get All Users
```bash
curl http://localhost:8080/api/v1/users
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe", "age": 30}'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## 🧪 Testing

```bash
go test ./test/...
```

## 📦 Build

```bash
go build -o bin/server cmd/server/main.go
```

## 🔧 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5433 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | 123456 |
| DB_NAME | Database name | postgres |
| DB_SSLMODE | SSL mode | disable |
| SERVER_PORT | Server port | 8080 |

## 📄 License

MIT

