package service

import (
	"errors"
	"sort"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/user/expense-tracker/internal/models"
	"github.com/user/expense-tracker/internal/repository"
)

type ExpenseService struct {
	repo *repository.Repository
}

func NewExpenseService(repo *repository.Repository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) CreateExpense(expense *models.Expense, splitData []struct {
	UserID uuid.UUID       `json:"user_id"`
	Amount decimal.Decimal `json:"amount"`
}) error {
	members, err := s.repo.GetGroupMembers(expense.GroupID)
	if err != nil {
		return err
	}

	if expense.SplitType == "EQUAL" {
		count := decimal.NewFromInt(int64(len(members)))
		if count.IsZero() {
			return errors.New("group has no members")
		}
		
		// Calculate amount per person, handling rounding
		amountPerPerson := expense.TotalAmount.DivRound(count, 2)
		
		// Check for remainder
		totalDistributed := amountPerPerson.Mul(count)
		diff := expense.TotalAmount.Sub(totalDistributed)

		var splits []models.ExpenseSplit
		for i, m := range members {
			amount := amountPerPerson
			if i == 0 {
				amount = amount.Add(diff) // Add remainder to first person
			}
			splits = append(splits, models.ExpenseSplit{
				UserID: m.ID,
				Amount: amount,
			})
		}
		expense.Splits = splits
	} else if expense.SplitType == "SHARE" {
		var splits []models.ExpenseSplit
		sum := decimal.Zero
		for _, sd := range splitData {
			splits = append(splits, models.ExpenseSplit{
				UserID: sd.UserID,
				Amount: sd.Amount,
			})
			sum = sum.Add(sd.Amount)
		}

		if !sum.Equal(expense.TotalAmount) {
			return errors.New("sum of shares does not equal total amount")
		}
		expense.Splits = splits
	} else {
		return errors.New("invalid split type")
	}

	return s.repo.CreateExpense(expense)
}

func (s *ExpenseService) GetSettlements(groupID uuid.UUID) ([]models.Transaction, error) {
	balances, err := s.repo.GetBalances(groupID)
	if err != nil {
		return nil, err
	}
	
	// Import the algorithm (internal/algorithm)
	// For simplicity in this structure, we'll call the algorithm directly.
	// In a real project, we might use an interface or a dedicated settlement service.
	return MinimizeTransactions(balances), nil
}

func (s *ExpenseService) SettleGroup(groupID uuid.UUID) error {
	balances, err := s.repo.GetBalances(groupID)
	if err != nil {
		return err
	}

	// For physical settlement, we update all balances to 0
	for _, b := range balances {
		err = s.repo.UpdateBalance(groupID, b.UserID, decimal.Zero)
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper to avoid circular dependency or local copy
func MinimizeTransactions(balances []models.UserBalance) []models.Transaction {
	type person struct {
		id      uuid.UUID
		balance decimal.Decimal
	}

	var debtors []person
	var creditors []person

	for _, b := range balances {
		if b.Balance.IsNegative() {
			debtors = append(debtors, person{id: b.UserID, balance: b.Balance.Abs()})
		} else if b.Balance.IsPositive() {
			creditors = append(creditors, person{id: b.UserID, balance: b.Balance})
		}
	}

	var transactions []models.Transaction

	sort.Slice(debtors, func(i, j int) bool {
		return debtors[i].balance.GreaterThan(debtors[j].balance)
	})
	sort.Slice(creditors, func(i, j int) bool {
		return creditors[i].balance.GreaterThan(creditors[j].balance)
	})

	i, j := 0, 0
	for i < len(debtors) && j < len(creditors) {
		debtor := &debtors[i]
		creditor := &creditors[j]
		amount := decimal.Min(debtor.balance, creditor.balance)
		if amount.IsPositive() {
			transactions = append(transactions, models.Transaction{From: debtor.id, To: creditor.id, Amount: amount})
		}
		debtor.balance = debtor.balance.Sub(amount)
		creditor.balance = creditor.balance.Sub(amount)
		if debtor.balance.IsZero() { i++ }
		if creditor.balance.IsZero() { j++ }
	}
	return transactions
}
