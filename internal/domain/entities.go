package domain

// Summary representa o resumo de trabalhos de um livro
type Summary struct {
	TotalTrabalhos int
}

// Localidade representa uma casa de oração
type Localidade struct {
	Nome   string
	Livros map[string]*Summary
}

// Setor representa um setor administrativo
type Setor struct {
	Nome        string
	Localidades []string
	Responsavel string
}

// RelatorioConfig representa a configuração do relatório
type RelatorioConfig struct {
	LarguraLivro          float64
	LarguraTotalTrabalhos float64
	LarguraApontamentos   float64
	MargemPagina          float64
}
