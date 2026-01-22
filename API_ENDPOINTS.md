# API Endpoints Quick Reference

## Base URL
`http://localhost:8080`

## Status Codes
- **200** - OK
- **201** - Created
- **204** - No Content
- **400** - Bad Request
- **404** - Not Found
- **500** - Internal Server Error

---

## USER ENDPOINTS (8 endpoints)

| Method | Endpoint | Description | Auth | Body |
|--------|----------|-------------|------|------|
| POST | `/api/users/register` | Register new user | No | `{email, name}` |
| GET | `/api/users/login?email=...` | Login with email | No | - |
| GET | `/api/users/{id}` | Get user details | Yes | - |
| GET | `/api/users` | List all users | Yes | - |
| PUT | `/api/users/{id}` | Update user name | Yes | `{name}` |
| DELETE | `/api/users/{id}` | Delete user | Yes | - |

---

## GROUP ENDPOINTS (6 endpoints)

| Method | Endpoint | Description | Auth | Body |
|--------|----------|-------------|------|------|
| POST | `/api/groups` | Create group | Yes | `{name, description, creator_id}` |
| GET | `/api/groups/{id}` | Get group details | Yes | - |
| GET | `/api/groups` | List all groups | Yes | - |
| GET | `/api/users/{user_id}/groups` | Get user's groups | Yes | - |
| PUT | `/api/groups/{id}` | Update group | Yes | `{name, description}` |
| DELETE | `/api/groups/{id}` | Delete group | Yes | - |

---

## GROUP MEMBERS ENDPOINTS (3 endpoints)

| Method | Endpoint | Description | Auth | Body |
|--------|----------|-------------|------|------|
| POST | `/api/groups/members` | Add member to group | Yes | `{group_id, user_id}` |
| GET | `/api/groups/{group_id}/members` | Get group members | Yes | - |
| DELETE | `/api/groups/{group_id}/members/{user_id}` | Remove member | Yes | - |

---

## EXPENSE ENDPOINTS (6 endpoints)

| Method | Endpoint | Description | Auth | Body |
|--------|----------|-------------|------|------|
| POST | `/api/expenses` | Create expense | Yes | `{group_id, paid_by_id, amount, description}` |
| GET | `/api/expenses/{id}` | Get expense details | Yes | - |
| GET | `/api/groups/{group_id}/expenses` | Get group expenses | Yes | - |
| GET | `/api/users/{user_id}/expenses` | Get user's expenses | Yes | - |
| PUT | `/api/expenses/{id}` | Update expense | Yes | `{amount, description}` |
| DELETE | `/api/expenses/{id}` | Delete expense | Yes | - |

---

## EXPENSE SPLIT ENDPOINTS (4 endpoints)

| Method | Endpoint | Description | Auth | Body |
|--------|----------|-------------|------|------|
| POST | `/api/expense-splits` | Add split to expense | Yes | `{expense_id, user_id, amount}` |
| GET | `/api/expenses/{expense_id}/splits` | Get expense splits | Yes | - |
| GET | `/api/users/{user_id}/splits` | Get user's splits | Yes | - |
| PUT | `/api/expense-splits/{id}` | Update split amount | Yes | `{amount}` |

---

## BALANCE ENDPOINTS (2 endpoints)

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/users/{user_id}/groups/{group_id}/balance` | Get user balance in group | Yes |
| GET | `/api/groups/{group_id}/balances` | Get all balances in group | Yes |

---

## SYSTEM ENDPOINTS (2 endpoints)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/metrics` | Prometheus metrics |

---

## TOTAL API COUNT: 32 Endpoints

---

## Response Format

### Success Response
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  ...other fields
}
```

### Error Response
```json
{
  "error": "Error message describing what went wrong"
}
```

### List Response
```json
[
  { ...object1 },
  { ...object2 },
  { ...objectN }
]
```

### Balance Response
```json
{
  "user_id": 1,
  "group_id": 1,
  "balance": -25.50
}
```

### Group Balances Response
```json
[
  {
    "user_id": 1,
    "user_name": "John Doe",
    "amount": 25.50
  },
  {
    "user_id": 2,
    "user_name": "Jane Doe",
    "amount": -25.50
  }
]
```

---

## Balance Interpretation

- **Negative Balance** (-): User owes money to the group
- **Positive Balance** (+): User is owed money by the group
- **Zero Balance** (0): All settled

---

## Request Examples

### Register User
```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "name": "John Doe"
  }'
```

### Create Group
```bash
curl -X POST http://localhost:8080/api/groups \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Trip to Paris",
    "description": "Summer vacation",
    "creator_id": 1
  }'
```

### Create Expense
```bash
curl -X POST http://localhost:8080/api/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "paid_by_id": 1,
    "amount": 120.50,
    "description": "Hotel"
  }'
```

### Add Expense Split
```bash
curl -X POST http://localhost:8080/api/expense-splits \
  -H "Content-Type: application/json" \
  -d '{
    "expense_id": 1,
    "user_id": 2,
    "amount": 60.25
  }'
```

### Get Group Balances
```bash
curl http://localhost:8080/api/groups/1/balances
```

### Add User to Group
```bash
curl -X POST http://localhost:8080/api/groups/members \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "user_id": 2
  }'
```

---

## Common HTTP Status Codes

### 2xx Success
- **200 OK** - Request successful
- **201 Created** - Resource created successfully
- **204 No Content** - Successful deletion

### 4xx Client Error
- **400 Bad Request** - Invalid input or validation error
- **404 Not Found** - Resource doesn't exist

### 5xx Server Error
- **500 Internal Server Error** - Server processing error

---

## Metrics Available

### HTTP Metrics
- `http_requests_total` - Total requests by method, path, status
- `http_request_duration_seconds` - Request latency histogram
- `http_errors_total` - Error count by type

### Access Point
`GET /metrics`

Returns Prometheus format metrics for monitoring and alerting.

---

## Authentication Notes

Current implementation:
- No authentication required for register/login
- User ID required for authenticated endpoints
- Pass user ID as parameter in subsequent requests
- For production: Implement JWT tokens

---

## Data Types

| Type | Example | Notes |
|------|---------|-------|
| Integer | `1` | User/Group/Expense IDs |
| String | `"john@example.com"` | Email, name, description |
| Float | `120.50` | Expense amounts, balances |
| DateTime | `"2024-01-22T10:00:00Z"` | ISO 8601 format |

---

## Content-Type

All requests and responses use:
```
Content-Type: application/json
```

---

## Rate Limiting

Currently: No rate limiting
Production: Implement as needed

---

## Pagination

Currently: Not implemented
Future: Add limit/offset parameters

---

## Sorting

Currently: Default ordering by `created_at DESC`
Future: Add sortBy parameter

---

## Filtering

Currently: None
Future: Add filter parameters

---

## API Documentation

Full documentation: See `API_DOCUMENTATION.md`

---

## Deployment

- **Docker**: `docker-compose up`
- **AWS ECS**: See `AWS_ECS_DEPLOYMENT.md`
- **Local**: `go run cmd/main.go`

---

**Last Updated**: January 22, 2024
**API Version**: 1.0.0
