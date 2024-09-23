# Andressa Lanches

Bem-vindo ao projeto **Andressa Lanches**! Esta é uma aplicação escrita em Go que fornece uma API para gerenciamento de vendas, produtos e acréscimos para a lanchonete da minha mãe. O projeto inclui autenticação JWT, validações, testes de integração, CI/CD com GitHub Actions e análise de qualidade de código com SonarQube.

## Sumário

- [Pré-requisitos](#pré-requisitos)
- [Instalação](#instalação)
  - [Configuração do Ambiente](#configuração-do-ambiente)
    - [Variáveis de Ambiente](#variáveis-de-ambiente)
    - [Banco de Dados](#banco-de-dados)
- [Executando a Aplicação](#executando-a-aplicação)
- [Executando os Testes](#executando-os-testes)
- [Gerando a Documentação Swagger](#gerando-a-documentação-swagger)
- [Configuração do SonarQube](#configuração-do-sonarqube)
- [Configuração do Lefthook (Githooks)](#configuração-do-lefthook-githooks)
- [Comandos Úteis do Makefile](#comandos-úteis-do-makefile)
- [Contribuindo](#contribuindo)
- [Licença](#licença)

---

## Pré-requisitos

Certifique-se de ter os seguintes softwares instalados em sua máquina:

- **Go** (versão 1.20 ou superior): [Instalação do Go](https://go.dev/dl/)
- **Docker**: [Instalação do Docker](https://docs.docker.com/get-docker/)
- **Docker Compose**: [Instalação do Docker Compose](https://docs.docker.com/compose/install/)
- **Git**: [Instalação do Git](https://git-scm.com/downloads)
- **Make** (Opcional, mas recomendado): Geralmente já instalado em sistemas Unix-like. Para Windows, pode ser necessário instalar manualmente.

---

## Instalação

1. **Clone o Repositório**

   ```bash
   git clone https://github.com/seu-usuario/andressa-lanches.git
   cd andressa-lanches

2. **Configurar o Ambiente**

   Execute o comando de setup para instalar as dependências e configurar o ambiente:

   ```bash
   make setup

Este comando irá:

- Instalar as dependências Go.
- Instalar o Lefthook para gerenciar os githooks.
- Instalar o golangci-lint para linting.
- Instalar o migrate para migrações de banco de dados.
- Instalar o swag para gerar a documentação Swagger.
- Configurar os githooks.

### Configuração do Ambiente

#### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto e configure as seguintes variáveis:

  ```env
  # Configurações da Aplicação
  SERVER_ADDRESS=:3333
  JWT_SECRET=sua_chave_secreta
  
  # Configurações do Banco de Dados
  DATABASE_URL=postgres://postgres:postgres@localhost:5432/mydatabase?sslmode=disable
  
  # Configuração do SonarQube
  SONAR_LOGIN=seu_token_sonar
  
  # Login
  AUTH_USER=user
  AUTH_PASSWORD=admin
  ```

#### Banco de Dados

O projeto utiliza o PostgreSQL como banco de dados. Você pode executá-lo usando o Docker Compose.

1. **Inicie o Banco de Dados**

```bash
docker-compose up -d db
```

2. **Execute as Migrações**

```bash
make migrate-up
```

Este comando executa as migrações de banco de dados, criando as tabelas necessárias.

---

## Executando a Aplicação

Após configurar o ambiente e o banco de dados, você pode executar a aplicação.

```bash
make run
```

A aplicação estará disponível em [http://localhost:3333](http://localhost:3333).

---

## Executando os Testes

Para executar todos os testes, incluindo os de integração:

```bash
make test
```

---

## Gerando a Documentação Swagger

A aplicação utiliza o Swag para gerar a documentação Swagger.

1. **Instalar o Swag (caso não tenha sido instalado no setup)**

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. **Gerar a Documentação**

```bash
make swagger
```

A documentação estará disponível em http://localhost:3333/swagger/index.html após iniciar a aplicação.

---

## Configuração do SonarQube

Para realizar a análise de qualidade do código com o SonarQube:

1. **Inicie o SonarQube**

```bash
docker-compose up -d sonarqube
```

Acesse o SonarQube em http://localhost:9000.

2. **Crie um Token no SonarQube**

- Faça login com as credenciais padrão (admin/admin).
- Altere a senha quando solicitado.
- Navegue até "My Account" > "Security".
- Gere um novo token e copie-o.

3. **Configure o Token**

Adicione o token ao arquivo .env:

```env
SONAR_LOGIN=seu_token_sonar
```

4. **Executar a Análise**

```bash
make sonar-scan
```

---

## Comandos Úteis do Makefile

O Makefile inclui vários comandos para facilitar o desenvolvimento. Para visualizar eles, execute:

```bash
make help
```

---

Obrigado por utilizar o Andressa Lanches!
