package application

import (
	"context"
	"testing"

	"github.com/phelliperodrigues/ccb-report/internal/domain"
	"github.com/phelliperodrigues/ccb-report/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockReportRepository is a mock implementation of domain.ReportRepository
type MockReportRepository struct {
	mock.Mock
}

func (m *MockReportRepository) Save(report *domain.Report) error {
	args := m.Called(report)
	return args.Error(0)
}

func (m *MockReportRepository) FindByLocalidade(localidade string) (*domain.Report, error) {
	args := m.Called(localidade)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Report), args.Error(1)
}

func (m *MockReportRepository) FindAll() ([]*domain.Report, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Report), args.Error(1)
}

// MockPDFGenerator is a mock implementation of infrastructure.PDFGenerator
type MockPDFGenerator struct {
	mock.Mock
}

func (m *MockPDFGenerator) GeneratePDF(report *domain.Report) error {
	args := m.Called(report)
	return args.Error(0)
}

func TestGenerateReport(t *testing.T) {
	tests := []struct {
		name           string
		localidade     string
		expectedError  error
		setupMocks     func(*MockReportRepository, *MockPDFGenerator)
		expectedReport *domain.Report
	}{
		{
			name:          "empty localidade",
			localidade:    "",
			expectedError: ErrInvalidLocalidade,
		},
		{
			name:       "successful report generation",
			localidade: "TestLocalidade",
			setupMocks: func(repo *MockReportRepository, pdfGen *MockPDFGenerator) {
				repo.On("Save", mock.Anything).Return(nil)
				pdfGen.On("GeneratePDF", mock.Anything).Return(nil)
			},
			expectedReport: &domain.Report{
				Localidade:     "TestLocalidade",
				Livros:         make(map[string]*domain.BookSummary),
				TotalTrabalhos: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			repo := new(MockReportRepository)
			pdfGen := new(MockPDFGenerator)
			cfg := config.New()
			service := NewReportService(repo, pdfGen, cfg)

			if tt.setupMocks != nil {
				tt.setupMocks(repo, pdfGen)
			}

			// Execute
			report, err := service.GenerateReport(context.Background(), tt.localidade)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, report)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, report)
				assert.Equal(t, tt.expectedReport.Localidade, report.Localidade)
			}

			// Verify mocks
			repo.AssertExpectations(t)
			pdfGen.AssertExpectations(t)
		})
	}
}

func TestGetReport(t *testing.T) {
	tests := []struct {
		name           string
		localidade     string
		expectedError  error
		setupMocks     func(*MockReportRepository)
		expectedReport *domain.Report
	}{
		{
			name:          "empty localidade",
			localidade:    "",
			expectedError: ErrInvalidLocalidade,
		},
		{
			name:       "report found",
			localidade: "TestLocalidade",
			setupMocks: func(repo *MockReportRepository) {
				repo.On("FindByLocalidade", "TestLocalidade").Return(&domain.Report{
					Localidade: "TestLocalidade",
				}, nil)
			},
			expectedReport: &domain.Report{
				Localidade: "TestLocalidade",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			repo := new(MockReportRepository)
			pdfGen := new(MockPDFGenerator)
			cfg := config.New()
			service := NewReportService(repo, pdfGen, cfg)

			if tt.setupMocks != nil {
				tt.setupMocks(repo)
			}

			// Execute
			report, err := service.GetReport(context.Background(), tt.localidade)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, report)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, report)
				assert.Equal(t, tt.expectedReport.Localidade, report.Localidade)
			}

			// Verify mocks
			repo.AssertExpectations(t)
		})
	}
}

func TestGetAllReports(t *testing.T) {
	tests := []struct {
		name            string
		expectedError   error
		setupMocks      func(*MockReportRepository)
		expectedReports []*domain.Report
	}{
		{
			name: "reports found",
			setupMocks: func(repo *MockReportRepository) {
				repo.On("FindAll").Return([]*domain.Report{
					{Localidade: "TestLocalidade1"},
					{Localidade: "TestLocalidade2"},
				}, nil)
			},
			expectedReports: []*domain.Report{
				{Localidade: "TestLocalidade1"},
				{Localidade: "TestLocalidade2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			repo := new(MockReportRepository)
			pdfGen := new(MockPDFGenerator)
			cfg := config.New()
			service := NewReportService(repo, pdfGen, cfg)

			if tt.setupMocks != nil {
				tt.setupMocks(repo)
			}

			// Execute
			reports, err := service.GetAllReports(context.Background())

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, reports)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, reports)
				assert.Equal(t, len(tt.expectedReports), len(reports))
				for i, expected := range tt.expectedReports {
					assert.Equal(t, expected.Localidade, reports[i].Localidade)
				}
			}

			// Verify mocks
			repo.AssertExpectations(t)
		})
	}
}
