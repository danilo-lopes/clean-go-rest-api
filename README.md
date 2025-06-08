# Clean Architecture Go REST API

Estrutura baseada em Clean Architecture para o projeto go-rest-api.

```
clean-go-rest-api/
.
│
├── cmd/
│   └── main.go                # Ponto de entrada da aplicação
│
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   └── user.go        # Entidade e interface User
│   │   └── dto/
│   │       └── user.go        # DTOs de entrada/saída
│   │
│   ├── usecase/
│   │   └── user.go            # Casos de uso de User
│   │
│   ├── adapter/
│   │   ├── handler/
│   │   │   └── user_handler.go # Handlers HTTP
│   │   └── repository/
│   │       └── postgres_user_repository.go # Repositório PostgreSQL
│   │
│   └── infrastructure/
│       ├── db/
│       │   └── migration.go   # Migrations
│       └── server/
│           └── server.go      # Inicialização do servidor HTTP
│
├── cmd/migrations/            # Arquivos SQL de migração
│
└── README.md
```

- O domínio (entidades, interfaces e DTOs) está isolado em `internal/domain`.
- Casos de uso em `internal/usecase`.
- Adaptadores (handlers HTTP e repositórios) em `internal/adapter`.
- Infraestrutura (migrations, servidor) em `internal/infrastructure`.
- O ponto de entrada (`main.go`) faz a orquestração das dependências em `cmd/`.
