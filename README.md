# CCB Report Generator

Um sistema web para geração de relatórios construído com Go usando arquitetura hexagonal.

## Estrutura do Projeto

```
.
├── cmd/
│   └── api/            # Ponto de entrada da aplicação
├── internal/
│   ├── domain/         # Modelos e interfaces de domínio
│   ├── application/    # Lógica de negócio
│   ├── infrastructure/ # Implementações externas
│   └── interfaces/     # Adaptadores de interface
├── pkg/                # Pacotes compartilhados
└── data/              # Armazenamento de dados
```

## Funcionalidades

- Geração de relatórios para diferentes localidades
- Visualização de relatórios existentes
- Interface API RESTful
- Armazenamento baseado em arquivos

## Endpoints da API

- `POST /api/reports/:localidade` - Gerar um novo relatório
- `GET /api/reports/:localidade` - Obter um relatório específico
- `GET /api/reports` - Listar todos os relatórios

## Configuração

O projeto usa variáveis de ambiente para configuração. Você pode configurar as seguintes variáveis:

- `BASE_PATH` - Caminho base para os arquivos (padrão: ".")
- `SERVER_PORT` - Porta do servidor (padrão: "8080")

## Estrutura dos Arquivos CSV

### books.csv
O arquivo deve estar em `files/books.csv` e ter o seguinte formato:
```csv
[12 linhas de cabeçalho]
livro,data,localidade
```

### input.csv
O arquivo deve estar em `files/input.csv` e ter o seguinte formato:
```csv
[12 linhas de cabeçalho]
localidade,data,livro
```

## Como Executar

### Usando Go localmente

1. Instale as dependências:
   ```bash
   make tidy
   ```

2. Configure o ambiente:
   ```bash
   make setup
   ```

3. Execute a aplicação:
   ```bash
   make run
   ```

### Usando Docker

1. Construa a imagem:
   ```bash
   make docker-build
   ```

2. Execute o container:
   ```bash
   make docker-run
   ```

3. Para parar o container:
   ```bash
   make docker-stop
   ```

O servidor estará disponível em `http://localhost:8080`

## Comandos Make

- `make build` - Compila o projeto
- `make test` - Executa os testes
- `make lint` - Executa o linter
- `make clean` - Remove arquivos gerados
- `make run` - Executa a aplicação
- `make setup` - Cria arquivos de exemplo
- `make tidy` - Atualiza dependências
- `make docker-build` - Constrói a imagem Docker
- `make docker-run` - Executa o container Docker
- `make docker-stop` - Para o container Docker

## Desenvolvimento

O projeto segue os princípios da arquitetura hexagonal:

- Camada de Domínio: Contém a lógica de negócio e interfaces
- Camada de Aplicação: Implementa os casos de uso
- Camada de Infraestrutura: Fornece implementações externas
- Camada de Interface: Lida com a comunicação externa

## Qualidade de Código

### Testes

Para executar os testes:
```bash
make test
```

### Linting

O projeto usa o golangci-lint para garantir a qualidade do código. Para executar o linter:
```bash
make lint
```

O arquivo `.golangci.yml` contém a configuração do linter com as seguintes verificações:
- gofmt - Formatação de código
- golint - Convenções de código
- govet - Análise estática
- errcheck - Verificação de erros
- staticcheck - Análise estática avançada
- gosimple - Simplificação de código
- ineffassign - Atribuições ineficientes
- unconvert - Conversões desnecessárias
- misspell - Erros de digitação
- gosec - Segurança

### Configuração do Editor

O projeto inclui configurações para facilitar o desenvolvimento:

#### VSCode
O arquivo `.vscode/settings.json` configura:
- Formatação automática ao salvar
- Organização automática de imports
- Configurações do linter
- Configurações de teste
- Decoradores de cobertura de código

#### EditorConfig
O arquivo `.editorconfig` padroniza:
- Estilo de indentação
- Fim de linha
- Charset
- Remoção de espaços em branco
- Configurações específicas para cada tipo de arquivo

## Dependências

- Gin Web Framework - Framework web
- gofpdf - Geração de PDF
- testify - Framework de testes
- golang.org/x/text - Processamento de texto

## Contribuindo

Veja [CONTRIBUTING.md](CONTRIBUTING.md) para detalhes sobre como contribuir com o projeto.

## Changelog

Veja [CHANGELOG.md](CHANGELOG.md) para um histórico de alterações.

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes. 