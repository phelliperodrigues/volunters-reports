package infrastructure

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
	"github.com/phelliperodrigues/ccb-report/internal/domain"
	"github.com/phelliperodrigues/ccb-report/pkg/config"
)

type PDFGenerator struct {
	basePath string
}

func NewPDFGenerator(basePath string) *PDFGenerator {
	return &PDFGenerator{
		basePath: basePath,
	}
}

func (g *PDFGenerator) GeneratePDF(report *domain.Report) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetFont("Arial", "", 12)
	pdf.SetTitle("Relatório de Localidades", true)
	pdf.SetAuthor("Phellipe Rodrigues", true)

	pageWidth, _ := pdf.GetPageSize()
	margin := 10.0
	usableWidth := pageWidth - 2*margin

	pdf.AddPage()

	// Título da localidade
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, tr("Relatório de Localidade: ")+report.Localidade)
	pdf.Ln(12)

	// Determina a largura das colunas
	larguraLivro := 40.0
	for livro := range report.Livros {
		larguraLivro = maxFloat64(larguraLivro, pdf.GetStringWidth(livro)+10)
	}

	larguraTotalTrabalhos := 38.0
	larguraSomaHoras := 38.0

	// Ajusta as larguras para caber na página
	totalTableWidth := larguraLivro + larguraTotalTrabalhos + larguraSomaHoras
	if totalTableWidth > usableWidth {
		scaleFactor := usableWidth / totalTableWidth
		larguraLivro *= scaleFactor
		larguraTotalTrabalhos *= scaleFactor
		larguraSomaHoras *= scaleFactor
	}

	// Cabeçalhos da tabela
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(larguraLivro, 7, "Livro", "1", 0, "C", false, 0, "")
	pdf.CellFormat(larguraTotalTrabalhos, 7, tr("Total Lançados"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(larguraSomaHoras, 7, "Apontamentos", "1", 1, "C", false, 0, "")

	// Dados da tabela
	pdf.SetFont("Arial", "", 12)
	for livro, summary := range report.Livros {
		if summary.TotalTrabalhos < 1 {
			pdf.SetTextColor(255, 0, 0)
		}
		pdf.CellFormat(larguraLivro, 7, tr(livro), "1", 0, "", false, 0, "")
		pdf.CellFormat(larguraTotalTrabalhos, 7, fmt.Sprintf("%d", summary.TotalTrabalhos), "1", 0, "C", false, 0, "")
		pdf.CellFormat(larguraSomaHoras, 7, "", "1", 1, "R", false, 0, "")
		pdf.SetTextColor(0, 0, 0) // Reseta a cor do texto
	}

	pdf.Ln(10)
	// Adiciona o resumo
	pdf.SetFont("Arial", "", 12)
	now := time.Now()

	pdf.MultiCell(0, 6, tr(fmt.Sprintf("O relatório acima apresenta um resumo dos trabalhos realizados.\nA tabela lista cada livro e o total de trabalhos lançados no SIGA até %s.", now.Format("02/01/2006 15:04"))), "", "", false)

	// Adiciona alertas
	pdf.SetFont("Arial", "B", 12)
	_, adminExists := report.Livros["4 - ADMINISTRAÇÃO"]
	mp, manutencaoExists := report.Livros["2 - MANUTENÇÃO PREVENTIVA"]
	_, brigadaExists := report.Livros["4 - BRIGADA DE INCÊNDIO"]

	if !adminExists || (!manutencaoExists || mp.TotalTrabalhos < 8) || !brigadaExists {
		pdf.SetTextColor(255, 0, 0)
		pdf.MultiCell(0, 8, tr("PONTOS DE ATENÇÃO:"), "", "", false)
	}
	if !adminExists {
		pdf.SetTextColor(237, 81, 14)
		pdf.MultiCell(0, 8, tr("> Não há apontamentos de ADMINISTRACAO."), "", "", false)
	}
	if !manutencaoExists || mp.TotalTrabalhos < 8 {
		pdf.SetTextColor(237, 81, 14)
		pdf.MultiCell(0, 8, tr("> Menos de 8 apontamentos de MANUTENCAO."), "", "", false)
	}
	if !brigadaExists {
		pdf.SetTextColor(237, 81, 14)
		pdf.MultiCell(0, 8, tr("> Não há apontamentos de BRIGADA DE INCENDIO."), "", "", false)
	}
	pdf.SetTextColor(0, 0, 0) // Reseta a cor do texto
	pdf.Ln(30)

	// Adicionar linhas de observacões
	pdf.CellFormat(190, 7, tr("OBSERVAÇÕES"), "1", 1, "C", false, 0, "")
	for i := 0; i < 10; i++ {
		pdf.CellFormat(190, 7, "", "1", 1, "C", false, 0, "")
	}

	// Determina o diretório de destino baseado no revisor
	diretorioDestino := "./files/output/outros"
	for revisor, localidades := range config.Revisores {
		for _, loc := range localidades {
			if report.Localidade == loc {
				diretorioDestino = fmt.Sprintf("./files/output/%s", revisor)
				break
			}
		}
		if diretorioDestino != "./files/output/outros" {
			break
		}
	}

	// Cria o diretório se não existir
	if err := os.MkdirAll(diretorioDestino, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Salva o arquivo PDF
	return pdf.OutputFileAndClose(fmt.Sprintf("%s/relatorio-%s.pdf", diretorioDestino, strings.ToLower(report.Localidade)))
}

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
