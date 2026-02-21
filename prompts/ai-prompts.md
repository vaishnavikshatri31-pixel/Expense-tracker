## Prompts
# Main Project Prompt
I am building a production-style REST API in Go (Golang) for an Expense Tracker with Bill Splitting functionality similar to Splitwise.

Requirements:

1. Features:
- Create user profile (name, phone number as unique identifier)
- Create group (trip group, roommates group etc.)
- Add members to group using phone number
- Add expense to group
- Split expense:
    a) Equally among members
    b) Share-wise (custom amount per user)
- Calculate net balances
- Implement settlement algorithm that minimizes number of transactions
- Update database after settlement

2. Technical Requirements:
- Use Go programming language
- Use Gin framework
- Use PostgreSQL
- Use proper layered architecture:
    - models
    - handlers
    - services
    - repository
    - algorithms
- Use proper decimal/money handling (avoid float)
- Store money in smallest currency unit (paise/cents) OR use decimal library
- Implement greedy settlement algorithm
- Provide REST endpoints with JSON input/output

3. Deliverables:
- Complete database schema
- Settlement algorithm implementation
- Example scenarios
- Clear comments explaining logic
- Clean project structure

Important:
- Focus on correctness of money calculations
- Avoid floating point precision errors
- Write production-quality structured code
- Include edge case handling

Generate full backend code with explanation.

# PROMPT FOR USER PROFILE FEATURE
Generate Go REST API code using Gin to create and fetch user profiles.

User fields:
- ID (UUID)
- Name
- Phone number (unique)
- CreatedAt

Implement:
- POST /users
- GET /users/:id

Ensure:
- Phone number uniqueness validation
- Proper error handling
- Database schema for PostgreSQL
- Layered architecture

# PROMPT FOR GROUP CREATION + ADD BY PHONE NUMBER
Generate Go REST API code to:

1. Create group
2. Add member to group using phone number

Requirements:
- Group has:
    ID
    Name
    CreatedBy
- Many-to-many relationship between users and groups
- When adding member:
    - Search user by phone number
    - Add to group if exists
    - Return error if phone not registered

Provide:
- Database schema
- Repository methods
- Service layer logic
- API endpoints:
    POST /groups
    POST /groups/:id/members

 # PROMPT FOR SPLITTING EXPENSE (EQUAL + SHAREWISE)
 
Generate Go service logic to add expense in a group.

Expense fields:
- ID
- GroupID
- PaidBy
- TotalAmount (use decimal or smallest unit)
- SplitType (EQUAL or SHARE)
- CreatedAt

Implement:

1. Equal split:
   Divide total amount equally among group members

2. Share-wise split:
   Accept JSON input like:
   [
     {user_id: X, amount: 500},
     {user_id: Y, amount: 300}
   ]

Validate:
- Sum of shares must equal total amount
- No floating point errors

Update balances table accordingly.

Provide full implementation with validation.

# PROMPT FOR SETTLEMENT ALGORITHM
Implement a settlement algorithm in Go that minimizes number of transactions between users in a group.

Steps:
1. Compute net balance for each user
   - Positive means user should receive
   - Negative means user owes

2. Separate into:
   - Creditors list
   - Debtors list

3. Apply greedy matching:
   - Match highest debtor with highest creditor
   - Transfer minimum amount
   - Update balances
   - Repeat until settled

4. Return optimized list of transactions.

Ensure:
- Time complexity explanation
- Proper money handling
- Clear comments explaining algorithm

# PROMPT FOR MONEY HANDLING STRATEGY
Explain and implement best practice for handling money in Go backend systems.

Requirements:
- Avoid float64
- Either:
   a) Use int64 and store money in smallest currency unit (paise)
   OR
   b) Use decimal library

Explain:
- Why float causes precision issues
- How rounding is handled
- How database schema should define money column (NUMERIC(12,2))

Provide example calculations.

# PROMPT FOR DATABASE SCHEMA
Design PostgreSQL schema for Expense Tracker API with:

Tables:
- users
- groups
- group_members
- expenses
- expense_splits
- balances

Use:
- UUID primary keys
- Foreign keys
- Proper decimal types (NUMERIC)
- Unique constraint on phone number
- Indexes for performance

Provide CREATE TABLE statements.




