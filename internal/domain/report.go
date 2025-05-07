package domain

import (
	"context"
	"time"
)

// Report represents a report entity in the domain
type Report struct {
	Localidade     string
	Livros         map[string]*BookSummary
	DataGeracao    time.Time
	TotalTrabalhos int
}

// BookSummary represents the summary of a book in a report
type BookSummary struct {
	Nome           string
	TotalTrabalhos int
}

// ReportRepository defines the interface for report persistence
type ReportRepository interface {
	Save(report *Report) error
	FindByLocalidade(localidade string) (*Report, error)
	FindAll() ([]*Report, error)
}

// ReportService defines the interface for report business logic
type ReportService interface {
	GenerateReport(ctx context.Context, localidade string) (*Report, error)
	GetReport(ctx context.Context, localidade string) (*Report, error)
	GetAllReports(ctx context.Context) ([]*Report, error)
}
