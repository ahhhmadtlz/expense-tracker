package param

import (
	"time"
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/entity"
)

type CreateTransactionRequest struct {
	CategoryID  uint    `json:"category_id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"` 
}

type CreateTransactionResponse struct {
	Transaction TransactionInfo `json:"transaction"`
}

type UpdateTransactionRequest struct {
	CategoryID  *uint    `json:"category_id"`
	Type        *string  `json:"type"`
	Amount      *float64 `json:"amount"`
	Description *string  `json:"description"`
	Date        *string  `json:"date"` 
}

type UpdateTransactionResponse struct {
	Transaction TransactionInfo `json:"transaction"`
}

type GetTransactionResponse struct {
	Transaction TransactionInfo `json:"transaction"`
}

type ListTransactionsRequest struct {
	Type       string `json:"type"`        
	CategoryID *uint  `json:"category_id"` 
	StartDate  string `json:"start_date"`  // Format: "2006-01-02"
	EndDate    string `json:"end_date"`    // Format: "2006-01-02"
}

type ListTransactionsResponse struct {
	Transactions []TransactionInfo `json:"transactions"`
	TotalIncome  float64           `json:"total_income"`
	TotalExpense float64           `json:"total_expense"`
	Balance      float64           `json:"balance"`
}

type DeleteTransactionResponse struct {
	Message string `json:"message"`
}

type TransactionInfo struct {
	ID          uint      `json:"id"`
	CategoryID  uint      `json:"category_id"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToTransactionInfo(tx entity.Transaction) TransactionInfo {
	return TransactionInfo{
		ID:          tx.ID,
		CategoryID:  tx.CategoryID,
		Type:        tx.Type.String(),
		Amount:      tx.Amount,
		Description: tx.Description,
		Date:        tx.Date,
		CreatedAt:   tx.CreatedAt,
		UpdatedAt:   tx.UpdatedAt,
	}
}