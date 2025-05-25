package usecase

import (
	"fmt"
	"time"

	"report/internal/domain"
)

// ReportGenerator define o caso de uso para geração de relatórios
type ReportGenerator struct {
	localidadeRepo domain.LocalidadeRepository
	setorRepo      domain.SetorRepository
	livroRepo      domain.LivroRepository
	pdfService     PDFService
}

// NewReportGenerator cria uma nova instância de ReportGenerator
func NewReportGenerator(
	localidadeRepo domain.LocalidadeRepository,
	setorRepo domain.SetorRepository,
	livroRepo domain.LivroRepository,
	pdfService PDFService,
) *ReportGenerator {
	return &ReportGenerator{
		localidadeRepo: localidadeRepo,
		setorRepo:      setorRepo,
		livroRepo:      livroRepo,
		pdfService:     pdfService,
	}
}

// GenerateReports gera todos os relatórios
func (g *ReportGenerator) GenerateReports() error {
	localidades, err := g.localidadeRepo.GetAll()
	if err != nil {
		return fmt.Errorf("erro ao obter localidades: %v", err)
	}

	livros, err := g.livroRepo.GetAll()
	if err != nil {
		return fmt.Errorf("erro ao obter livros: %v", err)
	}

	// Gera relatórios individuais
	for localidade, dadosLocalidade := range localidades {
		setor, err := g.setorRepo.GetByLocalidade(localidade)
		if err != nil {
			return fmt.Errorf("erro ao obter setor para localidade %s: %v", localidade, err)
		}

		err = g.generateLocalidadeReport(localidade, dadosLocalidade, setor)
		if err != nil {
			return fmt.Errorf("erro ao gerar relatório para localidade %s: %v", localidade, err)
		}
	}

	// Gera relatório resumo
	err = g.generateSummaryReport(localidades, livros)
	if err != nil {
		return fmt.Errorf("erro ao gerar relatório resumo: %v", err)
	}

	return nil
}

func (g *ReportGenerator) generateLocalidadeReport(
	localidade string,
	dados map[string]*domain.Summary,
	setor *domain.Setor,
) error {
	config := &domain.RelatorioConfig{
		LarguraLivro:          100.0,
		LarguraTotalTrabalhos: 45.0,
		LarguraApontamentos:   45.0,
		MargemPagina:          10.0,
	}

	reportData := &ReportData{
		Titulo:     fmt.Sprintf("Relatório - %s", localidade),
		Data:       time.Now(),
		Localidade: localidade,
		Livros:     dados,
		Config:     config,
	}

	outputPath := g.getOutputPath(setor, localidade)
	return g.pdfService.GenerateLocalidadeReport(reportData, outputPath)
}

func (g *ReportGenerator) generateSummaryReport(
	localidades map[string]map[string]*domain.Summary,
	livros map[string]map[string]bool,
) error {
	reportData := &ReportData{
		Titulo:      "Resumo de Todas as Localidades",
		Data:        time.Now(),
		Localidades: localidades,
		LivrosMap:   livros,
	}

	return g.pdfService.GenerateSummaryReport(reportData, "./files/output/resumo_localidades.pdf")
}

func (g *ReportGenerator) getOutputPath(setor *domain.Setor, localidade string) string {
	diretorio := "./files/output/outros"
	if setor != nil {
		diretorio = fmt.Sprintf("./files/output/%s", setor.Responsavel)
	}
	return fmt.Sprintf("%s/relatorio-%s.pdf", diretorio, localidade)
}

// ReportData contém os dados necessários para gerar um relatório
type ReportData struct {
	Titulo      string
	Data        time.Time
	Localidade  string
	Livros      map[string]*domain.Summary
	Localidades map[string]map[string]*domain.Summary
	LivrosMap   map[string]map[string]bool
	Config      *domain.RelatorioConfig
}
