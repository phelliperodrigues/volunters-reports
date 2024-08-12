package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/jung-kurt/gofpdf/v2"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var revisores = map[string][]string{
	"Daniel": {
		"JARDIM DAS LARANJEIRAS",
		"CASA GRANDE",
		"JARDIM DOS VELEIROS",
		"JARDIM DOS ALAMOS",
	},
	"Phellipe": {
		"VILA ESPERANCA",
		"VILA SAO JOSE",
		"PARQUE FLORESTAL",
		"JARDIM IPORANGA",
		"FAZENDA DO SCHUNK",
	},
	"Fabio": {
		"RECANTO DOS NOBRES",
		"JARDIM LALO",
		"JARDIM GUANHEMBU",
		"INTERLAGOS",
	},
}

type Summary struct {
	TotalTrabalhos int
}

func removeAccents(input string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isNonSpacingMark), norm.NFC)
	result, _, _ := transform.String(t, input)
	return result
}

func isNonSpacingMark(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Unicode non-spacing marks
}

func extractMiddleName(localidade string) string {
	parts := strings.Split(localidade, " - ")
	if len(parts) > 1 {
		return parts[1]
	}
	return localidade
}

func readBooksFile(filePath string) (map[string]map[string]bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // Delimitador é vírgula
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	booksMap := make(map[string]map[string]bool)
	newsRecords := records[12:]
	for _, record := range newsRecords[12:] {
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

func main() {
	filePath := "./files/imput.csv"
	booksFilePath := "./files/books.csv"

	// Lê a lista de livros e localizações
	booksMap, err := readBooksFile(booksFilePath)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo de livros:", err)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // Delimitador é vírgula
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Erro ao ler o arquivo CSV:", err)
		return
	}

	// Mapa para armazenar os dados agrupados
	data := make(map[string]map[string]*Summary)
	livrosEncontrados := make(map[string]bool)

	// Percorre os registros a partir da segunda linha (ignorando o cabeçalho)
	for _, record := range records[1:] {
		localidade := extractMiddleName(removeAccents(record[0]))
		livro := strings.TrimSpace(record[2])

		if _, exists := data[localidade]; !exists {
			data[localidade] = make(map[string]*Summary)
		}

		if _, exists := data[localidade][livro]; !exists {
			data[localidade][livro] = &Summary{}
		}

		data[localidade][livro].TotalTrabalhos++
		livrosEncontrados[livro] = true
	}

	for localidade, livros := range booksMap {

		// Cria o PDF
		pdf := gofpdf.New("P", "mm", "A4", "")
		tr := pdf.UnicodeTranslatorFromDescriptor("")
		pdf.SetFont("Arial", "", 12)
		pdf.SetTitle("Relatório de Localidades", true)
		pdf.SetAuthor("Phellipe Rodrigues", true)

		pageWidth, _ := pdf.GetPageSize()
		margin := 10.0
		usableWidth := pageWidth - 2*margin

		// Verifica se a localidade existe em data; caso contrário, inicializa com todos os valores zerados
		if _, exists := data[localidade]; !exists {
			data[localidade] = make(map[string]*Summary)
			for livro := range livros {
				data[localidade][livro] = &Summary{TotalTrabalhos: 0}
			}
		}

		pdf.AddPage()

		// Título da localidade
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(40, 10, tr("Relatório de Localidade: ")+localidade)
		pdf.Ln(12)

		// Determina a largura das colunas
		larguraLivro := 40.0
		for livro := range livrosEncontrados {
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
		pdf.CellFormat(larguraTotalTrabalhos, 7, tr("Total Lancados"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(larguraSomaHoras, 7, "Apontamentos", "1", 1, "C", false, 0, "")

		// Dados da tabela
		pdf.SetFont("Arial", "", 12)
		for livro := range livrosEncontrados {
			summary, exists := data[localidade][livro]
			totalTrabalhos := 0
			if exists {
				totalTrabalhos = summary.TotalTrabalhos
			}
			if livro == "" {
				continue
			}
			bookName := livro[4:]
			for book := range booksMap[localidade] {
				if strings.Contains(bookName, book) {
					if totalTrabalhos < 1 {
						pdf.SetTextColor(255, 0, 0)
					}
					pdf.CellFormat(larguraLivro, 7, tr(bookName), "1", 0, "", false, 0, "")
					pdf.CellFormat(larguraTotalTrabalhos, 7, fmt.Sprintf("%d", totalTrabalhos), "1", 0, "C", false, 0, "")
					pdf.CellFormat(larguraSomaHoras, 7, "", "1", 1, "R", false, 0, "")
					pdf.SetTextColor(0, 0, 0) // Reseta a cor do texto
				}
			}
		}

		pdf.Ln(10)
		// Adiciona o resumo
		pdf.SetFont("Arial", "", 12)
		now := time.Now()

		pdf.MultiCell(0, 6, tr(fmt.Sprintf("O relatório acima apresenta um resumo dos trabalhos realizados.\nA tabela lista cada livro e o total de trabalhos lançados no SIGA até %s.", now.Format("02/01/2006 15:04"))), "", "", false)

		// Adiciona alertas
		pdf.SetFont("Arial", "B", 12)
		_, adminExists := data[localidade]["4 - ADMINISTRAÇÃO"]
		mp, manutencaoExists := data[localidade]["2 - MANUTENÇÃO PREVENTIVA"]
		_, brigadaExists := data[localidade]["4 - BRIGADA DE INCÊNDIO"]

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
		pdf.CellFormat(190, 7, tr("OBSERVACÕES"), "1", 1, "C", false, 0, "")
		for i := 0; i < 10; i++ {
			pdf.CellFormat(190, 7, "", "1", 1, "C", false, 0, "")
		}

		// Verifica se a localidade está associada a algum revisor
		diretorioDestino := "./files/output/outros"
		for revisor, localidades := range revisores {
			for _, loc := range localidades {
				if localidade == loc {
					diretorioDestino = fmt.Sprintf("./files/output/%s", revisor)
					break
				}
			}
			if diretorioDestino != "./files/output/outros" {
				break
			}
		}

		// Cria o diretório se não existir
		if _, err := os.Stat(diretorioDestino); os.IsNotExist(err) {
			err = os.MkdirAll(diretorioDestino, os.ModePerm)
			if err != nil {
				fmt.Println("Erro ao criar diretório:", err)
				return
			}
		}

		// Salva o arquivo PDF no diretório apropriado
		err = pdf.OutputFileAndClose(fmt.Sprintf("%s/relatorio-%s.pdf", diretorioDestino, strings.ToLower(localidade)))
		if err != nil {
			fmt.Println("Erro ao salvar o arquivo PDF:", err)
			return
		}
	}

	// Salva o arquivo PDF

	fmt.Println("Relatório gerado com sucesso!")

}

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
