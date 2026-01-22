# Expense Go Collab Backend

A production-ready Group Expense Tracker backend built in Go with PostgreSQL.

## Features

- User management with email-based authentication
- Group creation and management
- Expense tracking and splitting
- Real-time balance calculation
- Prometheus metrics integration
- Docker containerization
- PostgreSQL database

## Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Docker (optional)

### Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```

3. Create `.env` file with configuration:
```
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=expense_tracker
```

4. Run the application:
```bash
go run cmd/main.go
```

## Docker Setup

### Build Docker Image

```bash
docker build -t expense-tracker:latest .
```

### Run with Docker Compose

```bash
docker-compose up
```

## API Endpoints

### User Management

- **Register User**: `POST /api/users/register`
  - Request: `{"email": "user@example.com", "name": "User Name"}`
  - Response: User object with ID

- **Login User**: `GET /api/users/login?email=user@example.com`
  - Response: User object with ID

- **Get User**: `GET /api/users/{id}`
- **Get All Users**: `GET /api/users`
- **Update User**: `PUT /api/users/{id}`
- **Delete User**: `DELETE /api/users/{id}`

### Group Management

- **Create Group**: `POST /api/groups`
  - Request: `{"name": "Group Name", "description": "...", "creator_id": 1}`

- **Get Group**: `GET /api/groups/{id}`
- **Get All Groups**: `GET /api/groups`
- **Get User Groups**: `GET /api/users/{user_id}/groups`
- **Update Group**: `PUT /api/groups/{id}`
- **Delete Group**: `DELETE /api/groups/{id}`

### Group Members

- **Add Member**: `POST /api/groups/members`
  - Request: `{"group_id": 1, "user_id": 2}`

- **Remove Member**: `DELETE /api/groups/{group_id}/members/{user_id}`
- **Get Members**: `GET /api/groups/{group_id}/members`

### Expenses

- **Create Expense**: `POST /api/expenses`
  - Request: `{"group_id": 1, "paid_by_id": 1, "amount": 100.50, "description": "..."}`

- **Get Expense**: `GET /api/expenses/{id}`
- **Get Group Expenses**: `GET /api/groups/{group_id}/expenses`
- **Get User Expenses**: `GET /api/users/{user_id}/expenses`
- **Update Expense**: `PUT /api/expenses/{id}`
- **Delete Expense**: `DELETE /api/expenses/{id}`

### Expense Splits

- **Add Split**: `POST /api/expense-splits`
  - Request: `{"expense_id": 1, "user_id": 2, "amount": 50.25}`

- **Get Splits**: `GET /api/expenses/{expense_id}/splits`
- **Get User Splits**: `GET /api/users/{user_id}/splits`
- **Update Split**: `PUT /api/expense-splits/{id}`

### Balance & Settlement

- **Get User Balance**: `GET /api/users/{user_id}/groups/{group_id}/balance`
  - Response: `{"user_id": 1, "group_id": 1, "balance": -25.50}`

- **Get Group Balances**: `GET /api/groups/{group_id}/balances`
  - Response: Array of user balances

## Monitoring & Metrics

### Health Check

- **Health**: `GET /health`
  - Response: `{"status": "ok"}`

### Prometheus Metrics

- **Metrics Endpoint**: `GET /metrics`

Available metrics:
- `http_requests_total` - Total HTTP requests by method, path, and status
- `http_request_duration_seconds` - Request latency histogram
- `http_errors_total` - Total HTTP errors by status code

## Project Structure

```
.
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── database.go        # Database configuration
│   ├── handler/               # HTTP handlers
│   │   ├── middleware.go      # Prometheus metrics middleware
│   │   ├── user_handler.go
│   │   ├── group_handler.go
│   │   ├── expense_handler.go
│   │   └── balance_handler.go
│   ├── model/                 # Data models
│   │   ├── user_model.go
│   │   ├── group_model.go
│   │   ├── expense_model.go
│   │   ├── expense_split_model.go
│   │   ├── group_member_model.go
│   │   └── balance_model.go
│   ├── repository/            # Repository interfaces
│   │   └── interfaces.go
│   ├── repositorypg/          # PostgreSQL implementations
│   │   ├── user_repository.go
│   │   ├── group_repository.go
│   │   ├── expense_repository.go
│   │   ├── group_member_repository.go
│   │   ├── expense_split_repository.go
│   │   └── balance_repository.go
│   └── service/               # Business logic
│       ├── user_service.go
│       ├── group_service.go
│       ├── expense_service.go
│       └── balance_service.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── .env
└── README.md
```

## Authentication

Currently uses simple email-based login without validation. For production, implement:
- JWT tokens
- Password hashing
- Email verification

## Error Handling

All endpoints return consistent error responses:
```json
{
  "error": "Error message"
}
```

## Future Enhancements

- Notification system for payment reminders
- Transaction history and audit logs
- Advanced settlement algorithms
- Mobile app integration
- Payment gateway integration
