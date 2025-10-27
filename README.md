# Evently - Event Management API

A Go-based event management API built with Gin, SQLite, and JWT authentication.

## Prerequisites

- Go 1.24.5 or higher
- Air (for live reload during development)
- SQLite3

## Getting Started

### 1. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install Air (optional but recommended for development)
go install github.com/cosmtrek/air@latest
```

### 2. Environment Setup

Create a `.env` file in the root directory:

```env
PORT=8080
JWT_SECRET=your-secret-key-here
DATABASE_URL=./data.db
```

### 3. Run Database Migrations

```bash
go run cmd/migrate/main.go up
```

### 4. Run the Application

**Option A: Using Air (Recommended for Development)**

```bash
air
```

**Option B: Using Go directly**

```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

## Project Structure

```
.
├── cmd/
│   ├── api/           # Main API application
│   └── migrate/       # Database migrations
├── internals/
│   ├── database/      # Database models and queries
│   └── env/           # Environment configuration
├── server/            # Server configurations
├── data.db            # SQLite database
├── .air.toml          # Air configuration
└── go.mod             # Go modules

```

## API Endpoints

- Authentication endpoints
- Event management endpoints

## Development

The project uses Air for live reload during development. Any changes to `.go` files will automatically rebuild and restart the server.

## License

MIT
