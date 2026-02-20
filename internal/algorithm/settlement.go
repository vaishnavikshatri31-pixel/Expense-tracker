package algorithm

import (
	"sort"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/user/expense-tracker/internal/models"
)

type SettlementResult struct {
	Transactions []models.Transaction `json:"transactions"`
}

// MinimizeTransactions uses a greedy approach to find the minimum number of transactions
// to settle all debts.
// 1. Separate users into Debtors (balance < 0) and Creditors (balance > 0).
// 2. Sort both lists by absolute balance value descending.
// 3. Match the largest debtor with the largest creditor.
// 4. Transfer the minimum of (debt, credit) amount.
// 5. Update balances and repeat until all are settled.
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

	// Sort to prioritize large amounts (Greedy)
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

		// Find the minimum amount to transfer
		amount := decimal.Min(debtor.balance, creditor.balance)

		if amount.IsPositive() {
			transactions = append(transactions, models.Transaction{
				From:   debtor.id,
				To:     creditor.id,
				Amount: amount,
			})
		}

		debtor.balance = debtor.balance.Sub(amount)
		creditor.balance = creditor.balance.Sub(amount)

		if debtor.balance.IsZero() {
			i++
		}
		if creditor.balance.IsZero() {
			j++
		}
	}

	return transactions
}
