package application

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/phelliperodrigues/ccb-report/internal/domain"
	"github.com/phelliperodrigues/ccb-report/internal/infrastructure"
	"github.com/phelliperodrigues/ccb-report/pkg/common"
	"github.com/phelliperodrigues/ccb-report/pkg/config"
)

// ReportService errors
var (
	ErrInvalidLocalidade = fmt.Errorf("localidade inválida")
	ErrInvalidCSVFormat  = fmt.Errorf("formato de CSV inválido")
	ErrFileNotFound      = fmt.Errorf("arquivo não encontrado")
)

type reportService struct {
	repo         domain.ReportRepository
	pdfGenerator *infrastructure.PDFGenerator
	config       *config.Config
}

// NewReportService creates a new instance of ReportService
func NewReportService(repo domain.ReportRepository, pdfGenerator *infrastructure.PDFGenerator, cfg *config.Config) domain.ReportService {
	return &reportService{
		repo:         repo,
		pdfGenerator: pdfGenerator,
		config:       cfg,
	}
}

// GenerateReport generates a report for a given localidade
func (s *reportService) GenerateReport(ctx context.Context, localidade string) (*domain.Report, error) {
	if localidade == "" {
		return nil, ErrInvalidLocalidade
	}

	// Read input files
	booksMap, err := s.readBooksFile(ctx, s.config.BooksFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read books file: %w", err)
	}

	// Read input data
	records, err := s.readInputFile(ctx, s.config.InputFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}

	// Validate CSV format
	if len(records) <= s.config.CSVHeaderRows {
		return nil, ErrInvalidCSVFormat
	}

	// Process data
	report := &domain.Report{
		Localidade:     localidade,
		Livros:         make(map[string]*domain.BookSummary),
		DataGeracao:    time.Now(),
		TotalTrabalhos: 0,
	}

	// Process records
	for _, record := range records[s.config.CSVHeaderRows:] {
		if len(record) < 3 {
			continue // Skip invalid records
		}

		recordLocalidade := common.ExtractMiddleName(common.RemoveAccents(record[0]))
		livro := strings.TrimSpace(record[2])

		if livro == "" || recordLocalidade != localidade {
			continue
		}

		// Validate if book exists in books map
		if booksMap != nil {
			if localBooks, exists := booksMap[recordLocalidade]; exists {
				if !localBooks[livro] {
					continue // Skip books not in the books map
				}
			}
		}

		if _, exists := report.Livros[livro]; !exists {
			report.Livros[livro] = &domain.BookSummary{
				Nome:           livro,
				TotalTrabalhos: 0,
			}
		}

		report.Livros[livro].TotalTrabalhos++
		report.TotalTrabalhos++
	}

	// Generate PDF
	if err := s.pdfGenerator.GeneratePDF(report); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Save report
	if err := s.repo.Save(report); err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}

	return report, nil
}

// GetReport retrieves a report for a given localidade
func (s *reportService) GetReport(ctx context.Context, localidade string) (*domain.Report, error) {
	if localidade == "" {
		return nil, ErrInvalidLocalidade
	}
	return s.repo.FindByLocalidade(localidade)
}

// GetAllReports retrieves all reports
func (s *reportService) GetAllReports(ctx context.Context) ([]*domain.Report, error) {
	return s.repo.FindAll()
}

// readBooksFile reads and validates the books CSV file
func (s *reportService) readBooksFile(ctx context.Context, filePath string) (map[string]map[string]bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) <= s.config.CSVHeaderRows {
		return nil, ErrInvalidCSVFormat
	}

	booksMap := make(map[string]map[string]bool)
	for _, record := range records[s.config.CSVHeaderRows:] {
		if len(record) < 3 {
			continue // Skip invalid records
		}

		livro := strings.TrimSpace(record[0])
		localidade := strings.TrimSpace(common.RemoveAccents(record[2]))

		if livro == "" || localidade == "" {
			continue
		}

		if _, exists := booksMap[localidade]; !exists {
			booksMap[localidade] = make(map[string]bool)
		}

		booksMap[localidade][livro] = true
	}

	return booksMap, nil
}

// readInputFile reads and validates the input CSV file
func (s *reportService) readInputFile(ctx context.Context, filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) <= s.config.CSVHeaderRows {
		return nil, ErrInvalidCSVFormat
	}

	return records, nil
}
