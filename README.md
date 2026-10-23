# Clean Architecture - Order System

Implementação do desafio de **Clean Architecture** em Go, expondo o use case de listagem de orders através de três interfaces simultaneamente: **REST**, **gRPC** e **GraphQL**.

## Tecnologias

- **Go 1.21**
- **MySQL 8.0**
- **gRPC** com protobuf
- **GraphQL** com gqlgen
- **REST** com chi router
- **Docker** e **Docker Compose**
- **golang-migrate** para migrations automáticas

## Como executar

> Único comando necessário:

```bash
docker compose up
```

O Docker Compose irá:
1. Subir o banco de dados MySQL
2. Aguardar o banco ficar saudável (healthcheck)
3. Aplicar as migrations automaticamente
4. Iniciar a aplicação nas três portas

## Portas

| Serviço  | Porta |
|----------|-------|
| REST     | 8000  |
| gRPC     | 50051 |
| GraphQL  | 8080  |

## Endpoints

### REST

**Criar order:**
```
POST http://localhost:8000/order
Content-Type: application/json

{
  "price": 100.0,
  "tax": 10.0
}
```

**Listar orders:**
```
GET http://localhost:8000/order
```

### GraphQL

Playground disponível em: `http://localhost:8080`

**Criar order (mutation):**
```graphql
mutation {
  createOrder(input: { price: 100.0, tax: 10.0 }) {
    id
    price
    tax
    finalPrice
  }
}
```

**Listar orders (query):**
```graphql
query {
  listOrders {
    id
    price
    tax
    finalPrice
  }
}
```

### gRPC

Usando [Evans](https://github.com/ktr0731/evans) ou qualquer cliente gRPC:

```bash
# Com evans
evans -r repl --host localhost --port 50051

# Listar orders
call ListOrders

# Criar order
call CreateOrder
```

## Arquitetura

```
.
├── cmd/ordersystem/        # Entrypoint da aplicação
├── configs/                # Configuração via env vars
├── internal/
│   ├── domain/
│   │   ├── entity/         # Order entity
│   │   ├── event/          # Eventos de domínio
│   │   └── repository/     # Interface do repositório
│   ├── infra/
│   │   ├── database/       # Implementação MySQL
│   │   ├── grpc/           # Server e protobuf gerado
│   │   ├── graphql/        # Resolvers e schema
│   │   └── web/            # Handlers REST e webserver
│   └── usecase/            # CreateOrder e ListOrders
├── pkg/events/             # Dispatcher de eventos
├── sql/migrations/         # Migrations SQL
├── docker-compose.yaml
├── Dockerfile
└── api.http                # Requisições prontas para teste
```

## Variáveis de ambiente

| Variável            | Padrão    |
|---------------------|-----------|
| DB_HOST             | localhost |
| DB_PORT             | 3306      |
| DB_USER             | root      |
| DB_PASSWORD         | root      |
| DB_NAME             | orders    |
| WEB_SERVER_PORT     | :8000     |
| GRPC_SERVER_PORT    | 50051     |
| GRAPHQL_SERVER_PORT | :8080     |
