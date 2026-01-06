package param

import "github.com/ahhhmadtlz/expense-tracker/internal/domain/report/entity"

type CreateReportRequest struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type CreateReportResponse struct {
	ReportID string `json:"report_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type GetReportResponse struct {
	Report entity.Report `json:"report"`
}

