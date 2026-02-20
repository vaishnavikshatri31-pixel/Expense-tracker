# Verification Guide: CRUD Operations

This guide provides step-by-step instructions to verify the core functionality of the Expense Tracker API using `curl`.

### Prerequisites
- Ensure the API is running: `.\go\bin\go.exe run cmd/api/main.go`
- Ensure PostgreSQL is running and the schema is initialized.

---

### 1. User Management (Create & Read)

**Create User A:**
```bash
curl -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{"name": "Alice", "phone": "1234567890"}'
```
*Note: Save the `id` from the response.*

**Create User B:**
```bash
curl -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{"name": "Bob", "phone": "9876543210"}'
```

**Fetch User Profile:**
Replace `{id}` with the ID returned from the creation step.
```bash
curl -X GET http://localhost:8080/users/{id}
```

---

### 2. Group Management

**Create a Group:**
Replace `{user_id}` with Alice's ID.
```bash
curl -X POST http://localhost:8080/groups \
-H "Content-Type: application/json" \
-d '{"name": "Goa Trip", "created_by": "{user_id}"}'
```
*Note: Save the group `id`.*

**Add Bob to the Group:**
Replace `{group_id}` with the group's ID.
```bash
curl -X POST http://localhost:8080/groups/{group_id}/members \
-H "Content-Type: application/json" \
-d '{"phone": "9876543210"}'
```

---

### 3. Expense Management (Splitting)

**Add Equal Split Expense:**
Alice pays $100 for Everyone. Replace IDs accordingly.
```bash
curl -X POST http://localhost:8080/expenses \
-H "Content-Type: application/json" \
-d '{
  "group_id": "{group_id}",
  "paid_by": "{alice_id}",
  "description": "Lunch",
  "total_amount": 100.00,
  "split_type": "EQUAL"
}'
```

**Add Share-wise Split Expense:**
Alice pays $50. Bob owes $20, Alice owes $30.
```bash
curl -X POST http://localhost:8080/expenses \
-H "Content-Type: application/json" \
-d '{
  "group_id": "{group_id}",
  "paid_by": "{alice_id}",
  "description": "Drinks",
  "total_amount": 50.00,
  "split_type": "SHARE",
  "splits": [
    {"user_id": "{bob_id}", "amount": 20.00},
    {"user_id": "{alice_id}", "amount": 30.00}
  ]
}'
```

---

### 4. Verification & Settlement

**Check Net Settlements:**
This uses the greedy algorithm to show exactly who owes whom.
```bash
curl -X GET http://localhost:8080/groups/{group_id}/settlements
```

**Settle the Group:**
Mark everything as physically paid and reset balances.
```bash
curl -X POST http://localhost:8080/groups/{group_id}/settle
```
