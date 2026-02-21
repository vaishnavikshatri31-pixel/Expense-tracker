# Design Document - Expense Tracker with Bill Splitting

## Overview
A production-ready REST API for tracking shared expenses among groups of friends.

## Architecture
- **Language**: Go (Golang)
- **Framework**: Gin Gonic
- **Database**: PostgreSQL
- **Money Handling**: `shopspring/decimal` library (avoids floating point errors)

### Layered Architecture
1. **Models**: Entity definitions.
2. **Handlers**: API request/response processing.
3. **Service**: Core business logic (Split calculations).
4. **Repository**: Database abstraction layer.
5. **Algorithms**: Transaction minimization logic.

## Money Handling Strategy
We use the `decimal.Decimal` type for all money operations. 
Database columns are defined as `NUMERIC(12,2)`.
Key rules:
- No `float64` for money.
- Rounding to 2 decimal places using `DivRound`.
- Remainder handling: In equal splits, any rounding differences (e.g., ₹10 / 3 = 3.333...₹) are assigned to the first member of the group to ensure the total sum remains exact.

## Settlement Algorithm: Greedy Match
The goal is to minimize the total number of transactions.
1. Each user's net balance is calculated: `Total Paid - Total Owed`.
2. Users are split into `Debtors` (negative balance) and `Creditors` (positive balance).
3. We sort both lists by absolute amount.
4. We match the largest debtor with the largest creditor and "transfer" the smaller of the two amounts.
5. This reduces the number of participants at each step, achieving an efficient (though not always mathematically optimal) solution that is standard for this problem.

## Database Schema
`users` <-> `groups` via `group_members` (M:N)
`expenses` -> `expense_splits` (1:N)
`user_balances` tracks the current net standing of each user per group.
