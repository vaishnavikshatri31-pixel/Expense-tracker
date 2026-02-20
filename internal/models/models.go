package models

import (
	"time"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type User struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	Name      string          `json:"name" db:"name"`
	Phone     string          `json:"phone" db:"phone"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

type Group struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	Name      string          `json:"name" db:"name"`
	CreatedBy uuid.UUID       `json:"created_by" db:"created_by"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	Members   []User          `json:"members,omitempty"`
}

type Expense struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	GroupID     uuid.UUID       `json:"group_id" db:"group_id"`
	PaidBy      uuid.UUID       `json:"paid_by" db:"paid_by"`
	Description string          `json:"description" db:"description"`
	TotalAmount decimal.Decimal `json:"total_amount" db:"total_amount"`
	SplitType   string          `json:"split_type" db:"split_type"` // "EQUAL", "SHARE"
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	Splits      []ExpenseSplit  `json:"splits,omitempty"`
}

type ExpenseSplit struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	ExpenseID uuid.UUID       `json:"expense_id" db:"expense_id"`
	UserID    uuid.UUID       `json:"user_id" db:"user_id"`
	Amount    decimal.Decimal `json:"amount" db:"amount"`
}

type UserBalance struct {
	GroupID uuid.UUID       `json:"group_id" db:"group_id"`
	UserID  uuid.UUID       `json:"user_id" db:"user_id"`
	Balance decimal.Decimal `json:"balance" db:"balance"`
}

type Transaction struct {
	From   uuid.UUID       `json:"from"`
	To     uuid.UUID       `json:"to"`
	Amount decimal.Decimal `json:"amount"`
}
