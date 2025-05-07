.PHONY: all build test clean run docker-build docker-run docker-stop lint

# Variáveis
BINARY_NAME=ccb-report
GO=go
DOCKER_COMPOSE=docker-compose
GOLANGCI_LINT=golangci-lint

all: clean build test

build:
	$(GO) build -o $(BINARY_NAME) cmd/api/main.go

test:
	$(GO) test ./... -v

lint:
	$(GOLANGCI_LINT) run

clean:
	rm -f $(BINARY_NAME)
	rm -f files/output/*.pdf
	rm -f data/reports/*.csv

run:
	$(GO) run cmd/api/main.go

setup:
	cp -n files/books.csv.example files/books.csv || true
	cp -n files/input.csv.example files/input.csv || true

tidy:
	$(GO) mod tidy
	$(GO) mod vendor

docker-build:
	$(DOCKER_COMPOSE) build

docker-run:
	$(DOCKER_COMPOSE) up -d

docker-stop:
	$(DOCKER_COMPOSE) down

help:
	@echo "Comandos disponíveis:"
	@echo "  make build         - Compila o projeto"
	@echo "  make test          - Executa os testes"
	@echo "  make lint          - Executa o linter"
	@echo "  make clean         - Remove arquivos gerados"
	@echo "  make run           - Executa a aplicação"
	@echo "  make setup         - Cria arquivos de exemplo"
	@echo "  make tidy          - Atualiza dependências"
	@echo "  make all           - Executa clean, build e test"
	@echo "  make docker-build  - Constrói a imagem Docker"
	@echo "  make docker-run    - Executa o container Docker"
	@echo "  make docker-stop   - Para o container Docker" 