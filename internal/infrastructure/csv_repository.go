package infrastructure

import (
	"encoding/csv"
	"os"
	"strings"
	"unicode"

	"report/internal/domain"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// CSVLocalidadeRepository implementa LocalidadeRepository
type CSVLocalidadeRepository struct {
	inputPath string
}

// CSVSetorRepository implementa SetorRepository
type CSVSetorRepository struct {
	setoresMap map[string]*domain.Setor
}

// CSVLivroRepository implementa LivroRepository
type CSVLivroRepository struct {
	booksPath string
}

// NewCSVRepositories cria novas instâncias dos repositórios
func NewCSVRepositories(inputPath, booksPath string) (*CSVLocalidadeRepository, *CSVSetorRepository, *CSVLivroRepository) {
	return &CSVLocalidadeRepository{inputPath: inputPath},
		&CSVSetorRepository{setoresMap: initSetoresMap()},
		&CSVLivroRepository{booksPath: booksPath}
}

func initSetoresMap() map[string]*domain.Setor {
	setores := make(map[string]*domain.Setor)

	// Setor 9.1
	setor91 := &domain.Setor{
		Nome:        "Setor 9.1",
		Responsavel: "Setor 9.1",
		Localidades: []string{
			"JARDIM DAS LARANJEIRAS",
			"CASA GRANDE",
			"JARDIM DOS VELEIROS",
			"JARDIM DOS ALAMOS",
			"VILA ESPERANCA",
			"VILA SAO JOSE",
			"PARQUE FLORESTAL",
			"JARDIM IPORANGA",
			"FAZENDA DO SCHUNK",
			"RECANTO DOS NOBRES",
			"JARDIM LALO",
			"JARDIM GUANHEMBU",
			"INTERLAGOS",
		},
	}
	for _, loc := range setor91.Localidades {
		setores[loc] = setor91
	}

	// Setor 9.2
	setor92 := &domain.Setor{
		Nome:        "Setor 9.2",
		Responsavel: "Setor 9.2",
		Localidades: []string{
			"CHACARA MARIETA",
			"CHACARAS SANTO AMARO",
			"ILHA DO BORORE",
			"ITAIM",
			"JARDIM ELIANE",
			"JARDIM LUCELIA",
			"JARDIM MARILDA",
			"JARDIM SANTA BARBARA",
			"JARDIM SAO BERNARDO",
			"JARDIM SAO RAFAEL",
			"JARDIM SETE DE SETEMBRO",
			"JARDIM TRES CORACOES",
			"PARQUE GRAJAU",
			"PARQUE RESIDENCIAL COCAIA",
		},
	}
	for _, loc := range setor92.Localidades {
		setores[loc] = setor92
	}

	// Setor 9.3
	setor93 := &domain.Setor{
		Nome:        "Setor 9.3",
		Responsavel: "Setor 9.3",
		Localidades: []string{
			"BARRAGEM",
			"CIDADE NOVA AMERICA",
			"COLONIA PAULISTA",
			"EMBURA",
			"ESTACAO EVANGELISTA DE SOUZA",
			"JARDIM DAS FONTES",
			"ENGENHEIRO MARSILAC",
			"PARELHEIROS",
			"PONTE SECA",
			"RECANTO ANA MARIA",
			"JARDIM SAO NORBERTO",
			"VARGEM GRANDE",
			"VILA ROSCHEL",
		},
	}
	for _, loc := range setor93.Localidades {
		setores[loc] = setor93
	}

	return setores
}

// GetAll retorna todas as localidades
func (r *CSVLocalidadeRepository) GetAll() (map[string]map[string]*domain.Summary, error) {
	file, err := os.Open(r.inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	data := make(map[string]map[string]*domain.Summary)
	for _, record := range records[12:] {
		localidade := extractMiddleName(removeAccents(record[0]))
		livro := strings.TrimSpace(record[2])
		if livro == "" {
			continue
		}

		if _, exists := data[localidade]; !exists {
			data[localidade] = make(map[string]*domain.Summary)
		}

		if _, exists := data[localidade][livro]; !exists {
			data[localidade][livro] = &domain.Summary{}
		}

		data[localidade][livro].TotalTrabalhos++
	}

	return data, nil
}

// Save implementa a interface LocalidadeRepository
func (r *CSVLocalidadeRepository) Save(localidade *domain.Localidade) error {
	return nil // Sistema somente leitura
}

// GetAll retorna todos os setores
func (r *CSVSetorRepository) GetAll() (map[string]*domain.Setor, error) {
	return r.setoresMap, nil
}

// GetByLocalidade retorna o setor de uma localidade
func (r *CSVSetorRepository) GetByLocalidade(localidade string) (*domain.Setor, error) {
	if setor, exists := r.setoresMap[localidade]; exists {
		return setor, nil
	}
	return nil, nil
}

// GetAll retorna todos os livros
func (r *CSVLivroRepository) GetAll() (map[string]map[string]bool, error) {
	file, err := os.Open(r.booksPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	booksMap := make(map[string]map[string]bool)
	for _, record := range records[12:] {
		livro := strings.TrimSpace(record[0])
		localidade := strings.TrimSpace(removeAccents(record[2]))

		if livro == "" {
			continue
		}

		if _, exists := booksMap[localidade]; !exists {
			booksMap[localidade] = make(map[string]bool)
		}

		booksMap[localidade][livro] = true
	}

	return booksMap, nil
}

// GetByLocalidade retorna os livros de uma localidade
func (r *CSVLivroRepository) GetByLocalidade(localidade string) (map[string]bool, error) {
	allBooks, err := r.GetAll()
	if err != nil {
		return nil, err
	}
	return allBooks[localidade], nil
}

func removeAccents(input string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isNonSpacingMark), norm.NFC)
	result, _, _ := transform.String(t, input)
	return result
}

func isNonSpacingMark(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

func extractMiddleName(localidade string) string {
	parts := strings.Split(localidade, " - ")
	if len(parts) > 1 {
		return parts[1]
	}
	return localidade
}
