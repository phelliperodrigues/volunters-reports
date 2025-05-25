package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"report/internal/domain"
	"report/internal/usecase"

	"github.com/jung-kurt/gofpdf/v2"
)

// GofpdfService implementa o PDFService usando a biblioteca gofpdf
type GofpdfService struct{}

// NewGofpdfService cria uma nova instância de GofpdfService
func NewGofpdfService() *GofpdfService {
	return &GofpdfService{}
}

// GenerateLocalidadeReport gera o relatório de uma localidade
func (s *GofpdfService) GenerateLocalidadeReport(data *usecase.ReportData, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetFont("Arial", "", 12)
	pdf.SetTitle(data.Titulo, true)
	pdf.SetAuthor("Phellipe Rodrigues", true)
	pdf.AddPage()

	// Cabeçalho
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, tr("Relatório de Trabalhos"))
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, tr("Localidade: "+data.Localidade))
	pdf.Ln(15)

	// Cabeçalhos da tabela
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(data.Config.LarguraLivro, 7, "Livro", "1", 0, "C", false, 0, "")
	pdf.CellFormat(data.Config.LarguraTotalTrabalhos, 7, tr("Total Lançados"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(data.Config.LarguraApontamentos, 7, "Apontamentos", "1", 1, "C", false, 0, "")

	// Dados da tabela
	pdf.SetFont("Arial", "", 12)
	for livro, summary := range data.Livros {
		if summary.TotalTrabalhos < 1 {
			pdf.SetTextColor(255, 0, 0)
		}
		pdf.CellFormat(data.Config.LarguraLivro, 7, tr(livro), "1", 0, "", false, 0, "")
		pdf.CellFormat(data.Config.LarguraTotalTrabalhos, 7, fmt.Sprintf("%d", summary.TotalTrabalhos), "1", 0, "C", false, 0, "")
		pdf.CellFormat(data.Config.LarguraApontamentos, 7, "", "1", 1, "R", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
	}

	// Adiciona alertas
	s.addAlerts(pdf, tr, data.Livros)

	// Adiciona data do relatório
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 6, tr(fmt.Sprintf("Relatório gerado em %s", data.Data.Format("02/01/2006 15:04"))), "", "", false)

	// Adiciona observações
	s.addObservacoes(pdf, tr)

	// Cria o diretório se não existir
	dir := strings.TrimSuffix(outputPath, filepath.Base(outputPath))
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("erro ao criar diretório: %v", err)
	}

	return pdf.OutputFileAndClose(outputPath)
}

// GenerateSummaryReport gera o relatório resumo
func (s *GofpdfService) GenerateSummaryReport(data *usecase.ReportData, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetFont("Arial", "", 8)
	pdf.SetTitle(data.Titulo, true)
	pdf.SetAuthor("Phellipe Rodrigues", true)
	pdf.AddPage()

	// Cabeçalhos
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(50, 50, "Localidade", "1", 0, "C", false, 0, "")

	// Cabeçalhos das colunas para os livros
	livrosEncontrados := make(map[string]bool)
	for _, livros := range data.LivrosMap {
		for livro := range livros {
			livrosEncontrados[livro] = true
		}
	}

	livroOrdem := make([]string, 0, len(livrosEncontrados))
	for livro := range livrosEncontrados {
		livroOrdem = append(livroOrdem, livro)
		x, y := pdf.GetXY()
		pdf.TransformBegin()
		pdf.TransformRotate(90, x+0.5, y+22.5)
		pdf.CellFormat(12, 50, tr(livro), "", 0, "C", false, 0, "")
		pdf.TransformEnd()
		pdf.SetXY(x+7, y)
	}
	pdf.Ln(-1)

	// Dados
	pdf.SetFont("Arial", "", 8)
	for localidade, livros := range data.Localidades {
		pdf.CellFormat(50, 5, localidade, "1", 0, "", false, 0, "")
		for _, livro := range livroOrdem {
			if summary, exists := livros[livro]; exists {
				if summary.TotalTrabalhos == 0 {
					pdf.SetTextColor(255, 0, 0)
				}
				pdf.CellFormat(7, 5, fmt.Sprintf("%d", summary.TotalTrabalhos), "1", 0, "C", false, 0, "")
				pdf.SetTextColor(0, 0, 0)
			} else {
				pdf.CellFormat(7, 5, "X", "1", 0, "C", false, 0, "")
			}
		}
		pdf.Ln(-1)
	}

	return pdf.OutputFileAndClose(outputPath)
}

func (s *GofpdfService) addAlerts(pdf *gofpdf.Fpdf, tr func(string) string, livros map[string]*domain.Summary) {
	_, adminExists := livros["4 - ADMINISTRAÇÃO"]
	mp, manutencaoExists := livros["2 - MANUTENÇÃO PREVENTIVA"]
	_, brigadaExists := livros["4 - BRIGADA DE INCÊNDIO"]

	if !adminExists || (!manutencaoExists || (manutencaoExists && mp.TotalTrabalhos < 8)) || !brigadaExists {
		pdf.SetFont("Arial", "B", 12)
		pdf.SetTextColor(255, 0, 0)
		pdf.MultiCell(0, 8, tr("PONTOS DE ATENÇÃO:"), "", "", false)

		if !adminExists {
			pdf.SetTextColor(237, 81, 14)
			pdf.MultiCell(0, 8, tr("> Não há apontamentos de ADMINISTRAÇÃO."), "", "", false)
		}
		if !manutencaoExists || (manutencaoExists && mp.TotalTrabalhos < 8) {
			pdf.SetTextColor(237, 81, 14)
			pdf.MultiCell(0, 8, tr("> Menos de 8 apontamentos de MANUTENÇÃO."), "", "", false)
		}
		if !brigadaExists {
			pdf.SetTextColor(237, 81, 14)
			pdf.MultiCell(0, 8, tr("> Não há apontamentos de BRIGADA DE INCÊNDIO."), "", "", false)
		}
	}
}

func (s *GofpdfService) addObservacoes(pdf *gofpdf.Fpdf, tr func(string) string) {
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 7, tr("OBSERVAÇÕES"), "1", 1, "C", false, 0, "")
	for i := 0; i < 10; i++ {
		pdf.CellFormat(190, 7, "", "1", 1, "C", false, 0, "")
	}
}
