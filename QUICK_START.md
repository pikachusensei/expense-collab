# Quick Start Guide

## 🚀 Get Running in 5 Minutes

### Option 1: Docker Compose (Recommended)

```bash
# Clone/navigate to project
cd /Users/shreyansh/expense-go-collab-backend

# Start all services (PostgreSQL + App)
docker-compose up -d

# Wait for services to be ready (check logs)
docker-compose logs -f

# Test the API
curl http://localhost:8080/health
```

**Result**: API running at `http://localhost:8080`

---

### Option 2: Local Development

```bash
# Prerequisites: Go 1.21+, PostgreSQL running locally

# Install dependencies
go mod downloadv

# Update .env for local PostgreSQL
# (default: localhost:5432, user: postgres, password: postgres)

# Build and run
make build
make run

# Or just
go run cmd/main.go
```

**Result**: API running at `http://localhost:8080`

---

## 📝 First API Calls

### 1. Register a User
```bash
curl -X POST http://localhost:8080/api/users/register \ 
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "name": "Alice"
  }'
```

**Response**:
```json
{
  "id": 1,
  "email": "alice@example.com",
  "name": "Alice",
  "created_at": "2024-01-22T10:00:00Z"
}
```

### 2. Register Another User
```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "bob@example.com",
    "name": "Bob"
  }'
```

### 3. Login
```bash
curl "http://localhost:8080/api/users/login?email=alice@example.com"
```

### 4. Create a Group
```bash
curl -X POST http://localhost:8080/api/groups \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Dinner Party",
    "description": "Weekend dinner expenses",
    "creator_id": 1
  }'
```

**Response** (note the group id):
```json
{
  "id": 1,
  "name": "Dinner Party",
  "description": "Weekend dinner expenses",
  "creator_id": 1,
  "created_at": "2024-01-22T10:00:00Z"
}
```

### 5. Add Member to Group
```bash
curl -X POST http://localhost:8080/api/members \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "user_id": 2
  }'
```

### 6. Create an Expense
```bash
curl -X POST http://localhost:8080/api/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "paid_by_id": 1,
    "amount": 120,
    "description": "Restaurant bill"
  }'
```

### 7. Add a Split
```bash
curl -X POST http://localhost:8080/api/splits \
  -H "Content-Type: application/json" \
  -d '{
    "expense_id": 1,
    "user_id": 2,
    "amount": 60
  }'
```

### 8. Check Balances
```bash
curl http://localhost:8080/api/balance/group/1
```

**Response** (negative = owes money, positive = owed money):
```json
[
  {
    "user_id": 1,
    "user_name": "Alice",
    "amount": 60
  },
  {
    "user_id": 2,
    "user_name": "Bob",
    "amount": -60
  }
]
```

### 9. Health Check
```bash
curl http://localhost:8080/health
```

### 10. View Metrics
```bash
curl http://localhost:8080/metrics
```

---

## 🔌 Complete Workflow Example

### Scenario: Two friends sharing a trip

```bash
# Step 1: Register users
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","name":"John"}'
# Save ID: 1

curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"email":"jane@example.com","name":"Jane"}'
# Save ID: 2

# Step 2: Create group
curl -X POST http://localhost:8080/api/groups \
  -H "Content-Type: application/json" \
  -d '{"name":"Paris Trip","description":"Summer vacation","creator_id":1}'
# Save ID: 1

# Step 3: Add Jane to group
curl -X POST http://localhost:8080/api/members \
  -H "Content-Type: application/json" \
  -d '{"group_id":1,"user_id":2}'

# Step 4: John pays for hotel ($200)
curl -X POST http://localhost:8080/api/expenses \
  -H "Content-Type: application/json" \
  -d '{"group_id":1,"paid_by_id":1,"amount":200,"description":"Hotel"}'
# Save ID: 1

# Step 5: Add split - Jane owes half
curl -X POST http://localhost:8080/api/splits \
  -H "Content-Type: application/json" \
  -d '{"expense_id":1,"user_id":2,"amount":100}'

# Step 6: John pays for meals ($150)
curl -X POST http://localhost:8080/api/expenses \
  -H "Content-Type: application/json" \
  -d '{"group_id":1,"paid_by_id":1,"amount":150,"description":"Meals"}'
# Save ID: 2

# Step 7: Add split - Jane owes half
curl -X POST http://localhost:8080/api/splits \
  -H "Content-Type: application/json" \
  -d '{"expense_id":2,"user_id":2,"amount":75}'

# Step 8: Check final balances
curl http://localhost:8080/api/balance/group/1

# Result: Jane owes John $175 (100 + 75)
```

---

## 📊 Docker Commands

```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f app
docker-compose logs -f postgres

# Stop services
docker-compose down

# Remove everything including data
docker-compose down -v

# Rebuild images
docker-compose up -d --build

# Database shell
docker-compose exec postgres psql -U postgres -d expense_tracker
```

---

## 🛠️ Development Commands

```bash
# Using Makefile
make help           # Show all commands
make build          # Build binary
make run            # Run application
make docker-up      # Start containers
make docker-down    # Stop containers
make clean          # Clean artifacts
make fmt            # Format code

# Manual commands
go mod download     # Download dependencies
go run cmd/main.go  # Run directly
go build -o app cmd/main.go  # Build binary
```

---

## 🔍 Testing Endpoints

### Test User Endpoints
```bash
# Get all users
curl http://localhost:8080/api/users

# Get specific user
curl http://localhost:8080/api/users/1

# Update user
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Updated"}'

# Delete user
curl -X DELETE http://localhost:8080/api/users/1
```

### Test Group Endpoints
```bash
# Get all groups
curl http://localhost:8080/api/groups

# Get user's groups
curl http://localhost:8080/api/groups/user/1

# Update group
curl -X PUT http://localhost:8080/api/groups/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Name","description":"Updated desc"}'
```

### Test Expense Endpoints
```bash
# Get all group expenses
curl http://localhost:8080/api/expenses/group/1

# Get user expenses
curl http://localhost:8080/api/expenses/user/1

# Get specific expense
curl http://localhost:8080/api/expenses/1
```

### Test Balance Endpoints
```bash
# Get user balance in group
curl http://localhost:8080/api/balance/user/1/group/1

# Get all balances in group
curl http://localhost:8080/api/balance/group/1
```

---

## 🔑 Important Endpoints Summary

| Endpoint | Example |
|----------|---------|
| Register | `POST /api/users/register` |
| Login | `GET /api/users/login?email=user@example.com` |
| Create Group | `POST /api/groups` |
| Create Expense | `POST /api/expenses` |
| Add Split | `POST /api/splits` |
| View Balances | `GET /api/balance/group/1` |
| **Create Settlement** | **`POST /api/settle`** |
| **View Settlements** | **`GET /api/settle/group/1`** |
| Health Check | `GET /health` |
| Metrics | `GET /metrics` |

---

## 📚 Documentation Links

- **Full API Docs**: `API_DOCUMENTATION.md`
- **Quick Reference**: `API_ENDPOINTS.md`
- **AWS Deployment**: `AWS_ECS_DEPLOYMENT.md`
- **Project Summary**: `PROJECT_SUMMARY.md`
- **Implementation Status**: `IMPLEMENTATION_CHECKLIST.md`

---

## 🐛 Troubleshooting

### Port Already in Use
```bash
# Change port in .env or kill process
docker-compose down
# Change PORT in .env
docker-compose up -d
```

### Database Connection Error
```bash
# Wait for PostgreSQL to start
docker-compose ps
docker-compose logs postgres

# If not healthy, restart
docker-compose down
docker-compose up -d
```

### API Not Responding
```bash
# Check application logs
docker-compose logs app

# Verify health endpoint
curl http://localhost:8080/health

# Check if port is correct
docker-compose ps
```

### Build Errors
```bash
# Clear dependencies
go clean
go mod tidy
go mod download

# Rebuild
go build -o app cmd/main.go
```

---

## 🌐 Access Points

After `docker-compose up -d`:

- **API Base URL**: `http://localhost:8080`
- **Health Check**: `http://localhost:8080/health`
- **Metrics**: `http://localhost:8080/metrics`
- **Database**: `localhost:5432`
  - User: `postgres`
  - Password: `postgres`
  - Database: `expense_tracker`

---

## 💳 Settlement / Payment API

Record and track payments between users to settle shared expenses.

### Create a Payment Settlement
```bash
curl -X POST http://localhost:8080/api/settle \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "from_user_id": 2,
    "to_user_id": 3,
    "amount": 150.50,
    "description": "Payment for dinner"
  }'
```

**Response** (201 Created):
```json
{
  "id": 1,
  "group_id": 1,
  "from_user_id": 2,
  "from_user_name": "Bob",
  "to_user_id": 3,
  "to_user_name": "Charlie",
  "amount": 150.50,
  "description": "Payment for dinner",
  "created_at": "2026-01-22T13:17:51.996257Z"
}
```

### Get Settlement by ID
```bash
curl http://localhost:8080/api/settle/1
```

### Get All Settlements in a Group
```bash
curl http://localhost:8080/api/settle/group/1
```

**Response**:
```json
{
  "settlements": [
    {
      "id": 1,
      "group_id": 1,
      "from_user_id": 2,
      "from_user_name": "Bob",
      "to_user_id": 3,
      "to_user_name": "Charlie",
      "amount": 150.50,
      "description": "Payment for dinner",
      "created_at": "2026-01-22T13:17:51.996257Z"
    },
    {
      "id": 2,
      "group_id": 1,
      "from_user_id": 3,
      "from_user_name": "Charlie",
      "to_user_id": 2,
      "to_user_name": "Bob",
      "amount": 75.25,
      "description": "Lunch reimbursement",
      "created_at": "2026-01-22T13:18:20.123456Z"
    }
  ],
  "total_amount": 225.75
}
```

### Get User's Settlements  
```bash
curl http://localhost:8080/api/settle/user/2
```

### Get All Settlements
```bash
curl http://localhost:8080/api/settle
```

### Settlement Workflow Example
```bash
# Step 1: Check group balances to see who owes whom
curl http://localhost:8080/api/balance/group/1

# Step 2: Create settlement payment
curl -X POST http://localhost:8080/api/settle \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "from_user_id": 2,
    "to_user_id": 1,
    "amount": 175.00,
    "description": "Final payment for Paris trip expenses"
  }'

# Step 3: View all settlements in group
curl http://localhost:8080/api/settle/group/1

# Step 4: View specific settlement  
curl http://localhost:8080/api/settle/1
```

---

## ✅ Verification Checklist

After starting services:

- [ ] Health endpoint returns `{"status":"ok"}`
- [ ] Can register a user
- [ ] Can login with email
- [ ] Can create a group
- [ ] Can add members to group
- [ ] Can create expenses
- [ ] Can add splits
- [ ] Can view balances
- [ ] Can create payment settlements
- [ ] Metrics endpoint accessible
- [ ] No errors in logs

---

## 🚀 Next Steps

1. **Read Documentation**
   - Full API Docs: `API_DOCUMENTATION.md`
   - Quick Reference: `API_ENDPOINTS.md`

2. **Test Endpoints**
   - Use curl commands above
   - Test all CRUD operations

3. **Frontend Integration**
   - Integrate with frontend using API endpoints
   - Reference response formats in documentation

4. **Deploy to AWS**
   - Follow `AWS_ECS_DEPLOYMENT.md`
   - Push to ECR
   - Deploy to ECS

---

## 📞 Support

- Check `API_DOCUMENTATION.md` for endpoint details
- Review `PROJECT_SUMMARY.md` for architecture overview
- See `IMPLEMENTATION_CHECKLIST.md` for features status

---

**Version**: 1.1.0
**Last Updated**: January 22, 2026
**Status**: Ready to Use ✅
