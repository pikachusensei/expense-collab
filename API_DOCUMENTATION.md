# Group Expense Tracker - API Documentation

## API Overview

Base URL: `http://localhost:8080`

All requests and responses use JSON format.

---

## Authentication

Currently, authentication is simple email-based:
- No password required
- Just provide email for registration and login
- User ID is returned for subsequent requests

---

## Response Codes

- **200 OK**: Successful GET, PUT requests
- **201 Created**: Successful POST requests
- **204 No Content**: Successful DELETE requests
- **400 Bad Request**: Invalid input or validation errors
- **404 Not Found**: Resource not found
- **500 Internal Server Error**: Server error

---

## User Endpoints

### 1. Register User
- **Endpoint**: `POST /api/users/register`
- **Description**: Create a new user account
- **Request Body**:
```json
{
  "email": "user@example.com",
  "name": "John Doe"
}
```
- **Response** (201):
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "created_at": "2024-01-22T10:00:00Z"
}
```

### 2. Login User
- **Endpoint**: `GET /api/users/login?email=user@example.com`
- **Description**: Login with email
- **Response** (200):
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "created_at": "2024-01-22T10:00:00Z"
}
```

### 3. Get User by ID
- **Endpoint**: `GET /api/users/{id}`
- **Description**: Retrieve specific user details
- **Response** (200): User object

### 4. Get All Users
- **Endpoint**: `GET /api/users`
- **Description**: Get list of all users
- **Response** (200): Array of user objects

### 5. Update User
- **Endpoint**: `PUT /api/users/{id}`
- **Description**: Update user information
- **Request Body**:
```json
{
  "name": "Jane Doe"
}
```
- **Response** (200): Updated user object

### 6. Delete User
- **Endpoint**: `DELETE /api/users/{id}`
- **Description**: Delete user account
- **Response** (204): No content

---

## Group Endpoints

### 1. Create Group
- **Endpoint**: `POST /api/groups`
- **Description**: Create a new expense group
- **Request Body**:
```json
{
  "name": "Trip to Paris",
  "description": "Summer vacation expenses",
  "creator_id": 1
}
```
- **Response** (201):
```json
{
  "id": 1,
  "name": "Trip to Paris",
  "description": "Summer vacation expenses",
  "creator_id": 1,
  "created_at": "2024-01-22T10:00:00Z"
}
```

### 2. Get Group by ID
- **Endpoint**: `GET /api/groups/{id}`
- **Response** (200): Group object

### 3. Get All Groups
- **Endpoint**: `GET /api/groups`
- **Response** (200): Array of group objects

### 4. Get User's Groups
- **Endpoint**: `GET /api/users/{user_id}/groups`
- **Description**: Get all groups a user belongs to
- **Response** (200): Array of group objects

### 5. Update Group
- **Endpoint**: `PUT /api/groups/{id}`
- **Request Body**:
```json
{
  "name": "Updated Group Name",
  "description": "Updated description"
}
```
- **Response** (200): Updated group object

### 6. Delete Group
- **Endpoint**: `DELETE /api/groups/{id}`
- **Response** (204): No content

---

## Group Member Endpoints

### 1. Add Member to Group
- **Endpoint**: `POST /api/groups/members`
- **Description**: Add a user to a group
- **Request Body**:
```json
{
  "group_id": 1,
  "user_id": 2
}
```
- **Response** (201):
```json
{
  "id": 1,
  "group_id": 1,
  "user_id": 2,
  "added_at": "2024-01-22T10:00:00Z"
}
```

### 2. Get Group Members
- **Endpoint**: `GET /api/groups/{group_id}/members`
- **Description**: Get all members of a group
- **Response** (200): Array of group member objects

### 3. Remove Member from Group
- **Endpoint**: `DELETE /api/groups/{group_id}/members/{user_id}`
- **Response** (204): No content

---

## Expense Endpoints

### 1. Create Expense
- **Endpoint**: `POST /api/expenses`
- **Description**: Create a new expense
- **Request Body**:
```json
{
  "group_id": 1,
  "paid_by_id": 1,
  "amount": 120.50,
  "description": "Restaurant bill"
}
```
- **Response** (201):
```json
{
  "id": 1,
  "group_id": 1,
  "paid_by_id": 1,
  "amount": 120.50,
  "description": "Restaurant bill",
  "created_at": "2024-01-22T10:00:00Z"
}
```

### 2. Get Expense by ID
- **Endpoint**: `GET /api/expenses/{id}`
- **Response** (200): Expense object

### 3. Get Group Expenses
- **Endpoint**: `GET /api/groups/{group_id}/expenses`
- **Description**: Get all expenses in a group
- **Response** (200): Array of expense objects

### 4. Get User Expenses
- **Endpoint**: `GET /api/users/{user_id}/expenses`
- **Description**: Get all expenses paid by a user
- **Response** (200): Array of expense objects

### 5. Update Expense
- **Endpoint**: `PUT /api/expenses/{id}`
- **Request Body**:
```json
{
  "amount": 125.00,
  "description": "Updated restaurant bill"
}
```
- **Response** (200): Updated expense object

### 6. Delete Expense
- **Endpoint**: `DELETE /api/expenses/{id}`
- **Description**: Delete expense and all associated splits
- **Response** (204): No content

---

## Expense Split Endpoints

### 1. Add Split to Expense
- **Endpoint**: `POST /api/expense-splits`
- **Description**: Create a split for an expense
- **Request Body**:
```json
{
  "expense_id": 1,
  "user_id": 2,
  "amount": 60.25
}
```
- **Response** (201):
```json
{
  "id": 1,
  "expense_id": 1,
  "user_id": 2,
  "amount": 60.25
}
```

### 2. Get Splits for Expense
- **Endpoint**: `GET /api/expenses/{expense_id}/splits`
- **Response** (200): Array of split objects

### 3. Get Splits for User
- **Endpoint**: `GET /api/users/{user_id}/splits`
- **Description**: Get all splits involving a user
- **Response** (200): Array of split objects

### 4. Update Split
- **Endpoint**: `PUT /api/expense-splits/{id}`
- **Request Body**:
```json
{
  "amount": 65.50
}
```
- **Response** (200): Updated split object

---

## Balance & Settlement Endpoints

### 1. Get User Balance in Group
- **Endpoint**: `GET /api/users/{user_id}/groups/{group_id}/balance`
- **Description**: Get how much a user owes/is owed in a group
- **Response** (200):
```json
{
  "user_id": 1,
  "group_id": 1,
  "balance": -25.50
}
```
**Note**: 
- Negative balance = User owes money
- Positive balance = User is owed money

### 2. Get All Balances in Group
- **Endpoint**: `GET /api/groups/{group_id}/balances`
- **Description**: Get balances for all members in a group
- **Response** (200):
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

## Monitoring Endpoints

### 1. Health Check
- **Endpoint**: `GET /health`
- **Description**: Check if server is running
- **Response** (200):
```json
{
  "status": "ok"
}
```

### 2. Prometheus Metrics
- **Endpoint**: `GET /metrics`
- **Description**: Get Prometheus metrics
- **Metrics Available**:
  - `http_requests_total` - Total requests by method, path, status
  - `http_request_duration_seconds` - Request latency
  - `http_errors_total` - Total errors by status

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message describing what went wrong"
}
```

### Common Errors

- **Invalid email**: "email and name are required"
- **Duplicate email**: "user with email already exists"
- **User not found**: "user not found"
- **Invalid group**: "group not found"
- **Invalid amount**: "amount must be greater than 0"
- **Member already exists**: "user is already a member of this group"

---

## Example Workflows

### Workflow 1: Setting up a group expense

1. **Register users**:
   ```
   POST /api/users/register
   POST /api/users/register
   ```

2. **Create group**:
   ```
   POST /api/groups
   ```

3. **Add members**:
   ```
   POST /api/groups/members (for each user)
   ```

4. **Create expense**:
   ```
   POST /api/expenses
   ```

5. **Add splits**:
   ```
   POST /api/expense-splits (for each participant)
   ```

6. **Check balances**:
   ```
   GET /api/groups/{group_id}/balances
   ```

### Workflow 2: Login and view user's groups

1. **Login**:
   ```
   GET /api/users/login?email=user@example.com
   ```

2. **Get user's groups**:
   ```
   GET /api/users/{user_id}/groups
   ```

3. **Get group expenses**:
   ```
   GET /api/groups/{group_id}/expenses
   ```

4. **Check user's balance in each group**:
   ```
   GET /api/users/{user_id}/groups/{group_id}/balance
   ```

---

## Deployment

### Docker Deployment

1. **Build image**:
   ```bash
   docker build -t expense-tracker:latest .
   ```

2. **Run with Docker Compose**:
   ```bash
   docker-compose up -d
   ```

3. **Access application**:
   - API: `http://localhost:8080`
   - Metrics: `http://localhost:8080/metrics`
   - Health: `http://localhost:8080/health`

### AWS ECS Deployment

1. Push image to ECR:
   ```bash
   aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
   docker tag expense-tracker:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/expense-tracker:latest
   docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/expense-tracker:latest
   ```

2. Create ECS task definition and service

3. Deploy to ECS cluster

---

## Rate Limiting

Currently no rate limiting implemented. For production, add:
- Middleware for rate limiting
- Request throttling
- DDoS protection

---

## Version

API Version: 1.0.0
Last Updated: January 22, 2024
