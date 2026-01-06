package entity

import "time"

type ReportType string
type ReportStatus string

const (
	TypeMonthly ReportType = "monthly"
	TypeYearly  ReportType = "yearly"
	TypeCustom  ReportType = "custom"
)

const (
	StatusPending    ReportStatus = "pending"
	StatusProcessing ReportStatus = "processing"
	StatusCompleted  ReportStatus = "completed"
	StatusFailed     ReportStatus = "fail"
)

type Report struct {
	ID           string
	UserID       uint
	Type         ReportType
	Year         int
	Month        int
	Status       ReportStatus
	Result       *ReportStatus
	ErrorMessage string
	CreatedAt    time.Time
	CompletedAt *time.Time
}


type ReportResult struct{
	Summary Summary `json:"summary"`
	Categories []CategoryBreakdown `json:"categories"`
	Transactions []TransactionSummary
}



type Summary struct{
 TotalIncome float64 `json:"total_income"`
 TotalExpenses float64 `json:"total_expenses"`
 NetSavings float64  `json:"net_savings"`
 Month string `json:"month"`
 Year int `json:"year"`
}

type CategoryBreakdown struct {
	CategoryID uint `json:"category_id"`
	CategoryName string `json:"category_name"`
	Type string `json:"type"`
	Amount float64 `json:"amount"`
	Count int `json:"count"`
}

type TransactionSummary struct {
	ID uint `json:"id"`
	Date string `json:"date"`
	Category string `json:"category"`
	Type string `json:"type"`
	Amount float64 `json:"amount"`
	Description string `json:"description"`
}