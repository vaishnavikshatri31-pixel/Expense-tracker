package algorithm

import (
	"testing"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/user/expense-tracker/internal/models"
)

func TestMinimizeTransactions(t *testing.T) {
	uA := uuid.New()
	uB := uuid.New()
	uC := uuid.New()

	// Scenario:
	// A owes B 100
	// B owes C 100
	// Expected: A owes C 100
	balances := []models.UserBalance{
		{UserID: uA, Balance: decimal.NewFromInt(-100)},
		{UserID: uB, Balance: decimal.NewFromInt(0)},
		{UserID: uC, Balance: decimal.NewFromInt(100)},
	}

	transactions := MinimizeTransactions(balances)

	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(transactions))
	}

	if transactions[0].From != uA || transactions[0].To != uC || !transactions[0].Amount.Equal(decimal.NewFromInt(100)) {
		t.Errorf("Unexpected transaction: %+v", transactions[0])
	}
}

func TestComplexSettlement(t *testing.T) {
	u1 := uuid.New()
	u2 := uuid.New()
	u3 := uuid.New()
	u4 := uuid.New()

	// Net Balances:
	// U1: -50 (owes 50)
	// U2: -30 (owes 30)
	// U3: +20 (owed 20)
	// U4: +60 (owed 60)
	balances := []models.UserBalance{
		{UserID: u1, Balance: decimal.NewFromInt(-50)},
		{UserID: u2, Balance: decimal.NewFromInt(-30)},
		{UserID: u3, Balance: decimal.NewFromInt(20)},
		{UserID: u4, Balance: decimal.NewFromInt(60)},
	}

	transactions := MinimizeTransactions(balances)

	// Since it's greedy:
	// Max debtor (U1: 50) matches Max creditor (U4: 60) -> U1 pays U4 50
	// U4 now needs 10.
	// Next Max debtor (U2: 30) matches Max creditor (U4: 10) -> U2 pays U4 10
	// U4 settled.
	// Next Max debtor (U2: 20) matches Max creditor (U3: 20) -> U2 pays U3 20
	// All settled.

	if len(transactions) > 3 {
		t.Errorf("Transaction count %d is higher than expected", len(transactions))
	}

	totalSettled := decimal.Zero
	for _, tx := range transactions {
		totalSettled = totalSettled.Add(tx.Amount)
	}

	if !totalSettled.Equal(decimal.NewFromInt(80)) {
		t.Errorf("Total settled amount %s doesn't match total debt 80", totalSettled.String())
	}
}
