package transaction

import (
	"time"

)

type TransactionType string

const (
	TypeIncome  TransactionType = "income"
	TypeExpense TransactionType = "expense"
)

func (t TransactionType) String() string {
	return string(t)
}

func MapToTransactionType(txType string) TransactionType {
	switch txType {
	case "income":
		return TypeIncome
	case "expense":
		return TypeExpense
	default:
		return TypeExpense
	}
}

type Transaction struct {
	ID         uint
	UserID     uint
	CategoryID uint
	Type       TransactionType
	Amount     float64
	Description string
	Date time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}