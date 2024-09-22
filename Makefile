# Defina o nome do aplicativo e a URL do banco de dados
DATABASE_URL := postgres://postgres:postgres@localhost:5432/mydatabase?sslmode=disable

# Diretórios
MIGRATIONS_DIR := ./db/migrations

# Variáveis de ambiente
export DATABASE_URL


.PHONY: build
build:
	@go build -o bin/andressa_lanches ./cmd/andressa-lanches

.PHONY: run
run:
	@go run ./cmd/andressa-lanches

# Criar uma nova migração
.PHONY: migrate-new
migrate-new:
ifndef name
	$(error Você deve especificar o nome da migração. Exemplo: make migrate-new name=create_users_table)
endif
	@echo "==> Criando nova migração: $(name)"
	@touch $(MIGRATIONS_DIR)/`date +%Y%m%d%H%M%S`_$(name).up.sql
	@touch $(MIGRATIONS_DIR)/`date +%Y%m%d%H%M%S`_$(name).down.sql

# Executar migrações "up"
.PHONY: migrate-up
migrate-up:
	@echo "==> Aplicando migrações..."
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) up

# Executar migrações "down"
.PHONY: migrate-down
migrate-down:
	@echo "==> Revertendo a última migração..."
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) down 1

# Executar migrações até uma versão específica
.PHONY: migrate-goto
migrate-goto:
ifndef version
	$(error Você deve especificar a versão. Exemplo: make migrate-goto version=20210908120000)
endif
	@echo "==> Migrando para a versão $(version)..."
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) goto $(version)

# Verificar versão atual das migrações
.PHONY: migrate-version
migrate-version:
	@echo "==> Versão atual das migrações:"
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) version

# Testes
.PHONY: test
test:
	@go test ./... -v

# Testes com relatório de cobertura
.PHONY: test-cover
test-cover:
	@go test ./... -coverprofile=coverage.out

# Relatório de cobertura em HTML
.PHONY: coverage
coverage: test-cover
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Relatório de cobertura gerado em coverage.html"

.PHONY: swag
swag:
	@swag init -g ./cmd/andressa-lanches/main.go

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: setup
setup:
	@echo "Installing Lefthook..."
	@go install github.com/evilmartians/lefthook@latest
	@echo "Setting up git hooks..."
	@lefthook install

# Ajuda
.PHONY: help
help:
	@echo "Comandos disponíveis:"
	@echo "  make migrate-new name=nome_da_migracao - Cria uma nova migração"
	@echo "  make migrate-up         - Executa todas as migrações pendentes"
	@echo "  make migrate-down       - Reverte a última migração"
	@echo "  make migrate-goto version=versao - Migra para uma versão específica"
	@echo "  make migrate-version    - Mostra a versão atual das migrações"
	@echo "  make test               - Executa os testes"
	@echo "  make test-cover         - Executa os testes e gera um relatório de cobertura"
	@echo "  make coverage           - Gera um relatório de cobertura em HTML"
	@echo "  make swag               - Gera a documentação da API com Swag"
	@echo "  make build              - Compila o aplicativo"
	@echo "  make run                - Executa o aplicativo"
	@echo "  make lint               - Executa o linter"
	@echo "  make setup              - Instala as dependências e configura os ganchos do git"
