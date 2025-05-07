package infrastructure

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/phelliperodrigues/ccb-report/internal/domain"
)

type fileRepository struct {
	basePath string
}

// NewFileRepository creates a new instance of file-based repository
func NewFileRepository(basePath string) domain.ReportRepository {
	return &fileRepository{
		basePath: basePath,
	}
}

func (r *fileRepository) Save(report *domain.Report) error {
	// Create directory if it doesn't exist
	dir := filepath.Join(r.basePath, "reports")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Save report data to CSV
	filePath := filepath.Join(dir, fmt.Sprintf("%s.csv", report.Localidade))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Livro", "TotalTrabalhos"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data
	for livro, summary := range report.Livros {
		if err := writer.Write([]string{
			livro,
			fmt.Sprintf("%d", summary.TotalTrabalhos),
		}); err != nil {
			return fmt.Errorf("failed to write data: %w", err)
		}
	}

	return nil
}

func (r *fileRepository) FindByLocalidade(localidade string) (*domain.Report, error) {
	filePath := filepath.Join(r.basePath, "reports", fmt.Sprintf("%s.csv", localidade))
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("report not found for localidade: %s", localidade)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("invalid report file format")
	}

	report := &domain.Report{
		Localidade: localidade,
		Livros:     make(map[string]*domain.BookSummary),
	}

	// Skip header row
	for _, record := range records[1:] {
		if len(record) < 2 {
			continue
		}

		livro := strings.TrimSpace(record[0])
		totalTrabalhos := 0
		fmt.Sscanf(record[1], "%d", &totalTrabalhos)

		report.Livros[livro] = &domain.BookSummary{
			Nome:           livro,
			TotalTrabalhos: totalTrabalhos,
		}
		report.TotalTrabalhos += totalTrabalhos
	}

	return report, nil
}

func (r *fileRepository) FindAll() ([]*domain.Report, error) {
	dir := filepath.Join(r.basePath, "reports")
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var reports []*domain.Report
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			localidade := strings.TrimSuffix(file.Name(), ".csv")
			report, err := r.FindByLocalidade(localidade)
			if err != nil {
				continue // Skip files that can't be read
			}
			reports = append(reports, report)
		}
	}

	return reports, nil
}
