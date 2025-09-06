# Kumparan API

Go API boilerplate menggunakan CloudWeGo Hertz, PostgreSQL (sqlx/pgx), Elasticsearch, dan Swagger/OpenAPI. Termasuk Docker dan docker-compose untuk pengembangan lokal.

## ✨ Features

- High-performance HTTP server dengan CloudWeGo Hertz
- PostgreSQL via sqlx (driver pgx) dengan health check koneksi dan Replica
- Integrasi Elasticsearch v7 untuk pencarian
- Swagger UI tersedia di `/swagger/index.html`
- Arsitektur berlapis: router, handler, service, domain, infra
- Konfigurasi via environment (`envconfig` + `.env`)
- Dockerfile dan docker-compose untuk stack lokal (App + Postgres + Elasticsearch + Kibana)

## 🌠 Tech Stack

- Framework: CloudWeGo Hertz (`github.com/cloudwego/hertz`)
- DB: PostgreSQL 15, `sqlx` dengan driver `pgx`
- Search: Elasticsearch 7.x (`github.com/olivere/elastic/v7`)
- Auth: JWT (library tersedia)
- Config: `envconfig`, `godotenv`
- Docs: Swagger/OpenAPI (via `hertz-contrib/swagger` dan `swaggo/files`)

## 📜 Swagger / OpenAPI

Swagger didaftarkan di main.go:
- Base Path: `/`
- Host (dev): `localhost:8080`
- UI: http://localhost:8080/swagger/index.html
- File definisi: docs/swagger.yaml dan docs/swagger.json  

Endpoint utama:

- **Author**
  - `POST /author/create`
  - `PUT  /author/update/{id}`
- **Article**
  - `POST /article/create`
  - `POST /article/create-bulk`
  - `PUT  /article/update/{id}`
  - `GET  /article/search?keyword=...`
  - `GET  /article/author/{id}`
  - `GET  /article/author-name?name=...`

## 📦 Requirements

- Go 1.24+
- Docker (opsional tapi direkomendasikan)
- PostgreSQL 13+ (jika jalan lokal tanpa Docker)
- Elasticsearch 7.x (jika jalan lokal tanpa Docker)

## ⚙️ Configuration

Environment variables didefinisikan di config/app_config.go:

```text
APP_NAME=HertzApp
PORT=8080
DB_HOST=postgres
DB_PORT=5432
DB_USER=hertz_user
DB_PASSWORD=hertz_pass
DB_NAME=hertz_db
DB_REPLICA_HOST=postgres-replica
DB_REPLICA_PORT=5432
JWT_SECRET=supersecret
ELASTIC_URL=http://elasticsearch:9200
```

Salin file contoh lalu sesuaikan:

```bash
cp .env.example .env
```

## 🚀 Menjalankan Aplikasi

### Opsi A — Docker Compose

```bash
# Menjalankan API + Postgres + Elasticsearch + Kibana
docker-compose up -d --build

# Layanan:
# API: http://localhost:8080
# Swagger UI: http://localhost:8080/swagger/index.html
# Postgres: localhost:5432 (user: hertz_user, db: hertz_db)
# Elasticsearch: http://localhost:9200
# Kibana: http://localhost:5601
```

### Opsi B — Lokal tanpa Docker

```bash
# Pastikan Postgres & Elasticsearch berjalan dan bisa diakses
# Atur .env sesuai environment Anda
go run main.go
```

## 🐂 Struktur Proyek

```text
.
├── api/
│   ├── handler/           # HTTP handlers
│   ├── router/            # Route definitions (SetupRouter)
│   └── service/           # Application services / use cases
├── config/                # Config loading dan struct
├── docs/                  # Swagger docs (swagger.yaml, swagger.json)
├── domain/
│   ├── articles/          # Domain Articles + migrations
│   ├── authors/           # Domain Authors + migrations
│   └── infra/             # Postgres, Elasticsearch, logger
├── internal/middleware/   # Middlewares (mis. auth)
├── k8s/                   # Kubernetes manifests
├── pkg/                   # Utilities / shared packages
├── Dockerfile
├── docker-compose.yaml
├── main.go
└── go.mod
```

## 🐃 Database & Migrations

```bash
# Apply migration manual ke Postgres (Docker)
docker cp domain/authors/db/migrations/20250904221348-create_authors_new.sql hertz-postgres:/tmp/
docker exec -it hertz-postgres psql -U hertz_user -d hertz_db -f /tmp/20250904221348-create_authors_new.sql

docker cp domain/articles/db/migrations/20250904221231-create_articles_new.sql hertz-postgres:/tmp/
docker exec -it hertz-postgres psql -U hertz_user -d hertz_db -f /tmp/20250904221231-create_articles_new.sql
```

> Saran: integrasikan tool migrasi (mis. golang-migrate) untuk otomatisasi.

## 🍷 Build Image Production

```bash
docker build -t hertz-boilerplate:latest .
docker run --rm -p 8080:8080 --env-file .env hertz-boilerplate:latest
```

> Entrypoint menggunakan `wait-for.sh` agar dependencies siap sebelum menjalankan app.

## 💡 Tips Development

- Loading config: config/config.go (load `.env` lalu map dengan `envconfig`)  
- Koneksi Postgres: domain/infra/postgres.go (DSN, ping timeout)  
- Client Elasticsearch: domain/infra/elasticsearch.go (ELASTIC_URL)  
- Routing: api/router/router.go — tambahkan endpoint baru dengan menambah handler/service
