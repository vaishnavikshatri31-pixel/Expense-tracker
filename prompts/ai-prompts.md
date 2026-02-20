# AI Assistance Log

AI was used as a coding assistant to help generate boilerplate and validate architectural decisions. The core system design, database modeling, and settlement optimization logic were defined before prompting.

### 1. User Module Assistance Prompt

I provided the following structured requirement:

> Create a Go-based REST endpoint for managing users with UUID as primary key and unique phone number constraint. Ensure PostgreSQL schema includes unique index and proper error handling for duplicates.

AI helped generate:
- Handler skeleton
- Repository pattern template
- Validation logic

### 2. Group Management Assistance Prompt

I defined:
- Many-to-many relationship between users and groups
- Add member by searching phone number
- Enforce foreign key constraints

AI helped scaffold:
- SQL schema draft
- Repository function template

### 3. Expense Splitting Logic Assistance

I specified:
- Equal split with remainder handling
- Share-wise split with validation
- No floating-point arithmetic

AI assisted in:
- Structuring service layer
- Decimal handling integration

### 4. Settlement Algorithm Assistance

I designed the greedy approach:
- Compute net balances
- Separate debtors and creditors
- Match highest amounts first

AI helped refine:
- Edge case handling
- Clean implementation formatting

### 5. Money Handling Strategy Validation

I requested validation of:
- Using NUMERIC(12,2) in PostgreSQL
- Mapping to decimal.Decimal
- Avoiding float64 errors

AI confirmed best practices and provided reference implementation.
