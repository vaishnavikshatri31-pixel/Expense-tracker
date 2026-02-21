## Expense Tracker & Bill Splitting API
[VIEW ON GIT](https://github.com/vaishnavikshatri31-pixel/Expense-tracker)

A high-performance REST API built in Go for tracking shared expenses and optimizing settlements (similar to Splitwise).

## Features
- **User Management**: Create and fetch user profiles by UUID or Phone Number.
- **Group Management**: Organize users into groups (e.g., Roommates, Trips).
- **Expense Splitting**:
  - **Equal Split**: Automatically handles rounding remainders.
  - **Share-wise Split**: Custom amounts per user.
- **Settlement Optimization**: Greedy algorithm to minimize the total number of transactions needed to settle debts.
- **Precision Money Handling**: Uses the `shopspring/decimal` library to prevent floating-point errors.

## Tech Stack
- **Go** (1.21+)
- **Gin Gonic** (Web Framework)
- **PostgreSQL** (Relational Database)
- **Shopspring Decimal** (Financial precision)

## Getting Started

### 1. Database Setup
Create a PostgreSQL database named `expense_tracker` and run the migration:
```bash
psql -U postgres -d expense_tracker -f migrations/schema.sql
```

### 2. Configuration
Set the database URL in your environment:
```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/expense_tracker?sslmode=disable"
```

### 3. Run the API
```bash
go run cmd/api/main.go
```

## API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| POST | `/users` | Create a new user profile |
| GET | `/users/:id` | Fetch user profile by UUID |
| POST | `/groups` | Create a new group |
| POST | `/groups/:id/members` | Add member to group by phone number |
| POST | `/expenses` | Add a new expense with splitting logic |
| GET | `/groups/:id/settlements` | Calculate minimum transactions required |
| POST | `/groups/:id/settle` | Reset balances after physical payment |

## Settlement Optimization Example

Example:

A owes ₹100 to B,
B owes ₹100 to C

Normal transactions:
A → B (₹100),
B → C (₹100)

Optimized:
A → C (100)

Reduced from 2 transactions to 1.

## Documentation
- [Design Document](docs/design.md)
- [Prompts](prompts/ai-prompts.md)
