package usecase

// PDFService define as operações para geração de PDFs
type PDFService interface {
	GenerateLocalidadeReport(data *ReportData, outputPath string) error
	GenerateSummaryReport(data *ReportData, outputPath string) error
}
