# CCB Report Generator

Este projeto é um gerador de relatórios para a Congregação Cristã no Brasil (CCB). Ele processa arquivos CSV contendo informações sobre trabalhos realizados em diferentes localidades e gera relatórios em PDF organizados por setores.

## Estrutura do Projeto

O projeto segue os princípios da Clean Architecture:

```
.
├── cmd/
│   └── main.go           # Ponto de entrada da aplicação
├── internal/
│   ├── domain/           # Entidades e interfaces do domínio
│   ├── usecase/          # Casos de uso da aplicação
│   ├── infrastructure/   # Implementações concretas (repositórios, serviços)
│   └── interfaces/       # Adaptadores de interface (HTTP, CLI)
├── files/
│   ├── input.csv         # Arquivo de entrada com os trabalhos
│   └── books.csv         # Arquivo com a lista de livros por localidade
└── go.mod                # Dependências do projeto
```

## Funcionalidades

- Processamento de arquivos CSV
- Geração de relatórios individuais por localidade
- Geração de relatório resumo
- Organização por setores (9.1, 9.2, 9.3)
- Alertas para trabalhos faltantes ou insuficientes
- Seção de observações em cada relatório

## Como Usar

1. Coloque os arquivos CSV na pasta `files/`:
   - `input.csv`: Dados dos trabalhos
   - `books.csv`: Lista de livros por localidade

2. Execute o programa:
   ```bash
   go run cmd/main.go
   ```

3. Os relatórios serão gerados na pasta `files/output/`, organizados por setor.

## Dependências

- github.com/jung-kurt/gofpdf/v2: Geração de PDFs
- golang.org/x/text: Manipulação de texto e caracteres especiais

## Estrutura dos Arquivos CSV

### input.csv
```csv
Data,Localidade,Livro,Trabalho
2024-01-01,JARDIM DAS LARANJEIRAS,4 - ADMINISTRAÇÃO,Trabalho 1
...
```

### books.csv
```csv
Livro,Código,Localidade
4 - ADMINISTRAÇÃO,001,JARDIM DAS LARANJEIRAS
...
```

## Contribuindo

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Faça commit das suas alterações (`git commit -am 'Adiciona nova feature'`)
4. Faça push para a branch (`git push origin feature/nova-feature`)
5. Crie um Pull Request 