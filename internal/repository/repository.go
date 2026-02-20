package repository

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/user/expense-tracker/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// User Methods
func (r *Repository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (name, phone) VALUES ($1, $2) RETURNING id, created_at`
	return r.db.QueryRow(query, user.Name, user.Phone).Scan(&user.ID, &user.CreatedAt)
}

func (r *Repository) GetUserByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, phone, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Phone, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *Repository) GetUserByPhone(phone string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, phone, created_at FROM users WHERE phone = $1`
	err := r.db.QueryRow(query, phone).Scan(&user.ID, &user.Name, &user.Phone, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

// Group Methods
func (r *Repository) CreateGroup(group *models.Group) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO groups (name, created_by) VALUES ($1, $2) RETURNING id, created_at`
	err = tx.QueryRow(query, group.Name, group.CreatedBy).Scan(&group.ID, &group.CreatedAt)
	if err != nil {
		return err
	}

	// Add creator as first member
	memberQuery := `INSERT INTO group_members (group_id, user_id) VALUES ($1, $2)`
	_, err = tx.Exec(memberQuery, group.ID, group.CreatedBy)
	if err != nil {
		return err
	}

	// Initialize balance
	balanceQuery := `INSERT INTO user_balances (group_id, user_id, balance) VALUES ($1, $2, 0)`
	_, err = tx.Exec(balanceQuery, group.ID, group.CreatedBy)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) AddMemberToGroup(groupID uuid.UUID, userID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO group_members (group_id, user_id) VALUES ($1, $2)`
	_, err = tx.Exec(query, groupID, userID)
	if err != nil {
		return err
	}

	balanceQuery := `INSERT INTO user_balances (group_id, user_id, balance) VALUES ($1, $2, 0)`
	_, err = tx.Exec(balanceQuery, groupID, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetGroupMembers(groupID uuid.UUID) ([]models.User, error) {
	query := `
		SELECT u.id, u.name, u.phone, u.created_at 
		FROM users u
		JOIN group_members gm ON u.id = gm.user_id
		WHERE gm.group_id = $1`
	
	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Phone, &u.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, u)
	}
	return members, nil
}

// Expense Methods
func (r *Repository) CreateExpense(expense *models.Expense) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO expenses (group_id, paid_by, description, total_amount, split_type) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	err = tx.QueryRow(query, expense.GroupID, expense.PaidBy, expense.Description, expense.TotalAmount, expense.SplitType).Scan(&expense.ID, &expense.CreatedAt)
	if err != nil {
		return err
	}

	for i := range expense.Splits {
		expense.Splits[i].ExpenseID = expense.ID
		splitQuery := `INSERT INTO expense_splits (expense_id, user_id, amount) VALUES ($1, $2, $3) RETURNING id`
		err = tx.QueryRow(splitQuery, expense.ID, expense.Splits[i].UserID, expense.Splits[i].Amount).Scan(&expense.Splits[i].ID)
		if err != nil {
			return err
		}

		// Update Balances
		// The person who paid gets a POSITIVE increase in balance (owed more)
		// The person who owes gets a NEGATIVE decrease in balance (owes more)
		
		// If the person who paid is the same as the person in split, they net out part of it locally
		// Logic:
		// 1. Decrease borrower's balance by split amount
		balanceUpdateQuery := `UPDATE user_balances SET balance = balance - $1 WHERE group_id = $2 AND user_id = $3`
		_, err = tx.Exec(balanceUpdateQuery, expense.Splits[i].Amount, expense.GroupID, expense.Splits[i].UserID)
		if err != nil {
			return err
		}
	}

	// 2. Increase payer's balance by total amount
	payerUpdateQuery := `UPDATE user_balances SET balance = balance + $1 WHERE group_id = $2 AND user_id = $3`
	_, err = tx.Exec(payerUpdateQuery, expense.TotalAmount, expense.GroupID, expense.PaidBy)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetBalances(groupID uuid.UUID) ([]models.UserBalance, error) {
	query := `SELECT group_id, user_id, balance FROM user_balances WHERE group_id = $1`
	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []models.UserBalance
	for rows.Next() {
		var b models.UserBalance
		if err := rows.Scan(&b.GroupID, &b.UserID, &b.Balance); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}
	return balances, nil
}

func (r *Repository) UpdateBalance(groupID, userID uuid.UUID, newBalance decimal.Decimal) error {
	query := `UPDATE user_balances SET balance = $1 WHERE group_id = $2 AND user_id = $3`
	_, err := r.db.Exec(query, newBalance, groupID, userID)
	return err
}
