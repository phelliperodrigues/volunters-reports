package main

import (
	"fmt"
	"log"
	"os"

	"report/internal/infrastructure"
	"report/internal/usecase"
)

func main() {
	// Configuração dos caminhos dos arquivos
	inputPath := "./files/input.csv"
	booksPath := "./files/books.csv"

	// Verifica se os arquivos existem
	if err := checkFiles(inputPath, booksPath); err != nil {
		log.Fatal(err)
	}

	// Inicializa os repositórios
	localidadeRepo, setorRepo, livroRepo := infrastructure.NewCSVRepositories(inputPath, booksPath)

	// Inicializa o serviço de PDF
	pdfService := infrastructure.NewGofpdfService()

	// Inicializa o gerador de relatórios
	reportGenerator := usecase.NewReportGenerator(localidadeRepo, setorRepo, livroRepo, pdfService)

	// Gera os relatórios
	if err := reportGenerator.GenerateReports(); err != nil {
		log.Fatal("Erro ao gerar relatórios:", err)
	}

	fmt.Println("Relatórios gerados com sucesso!")
}

func checkFiles(paths ...string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("arquivo não encontrado: %s", path)
		}
	}
	return nil
}
