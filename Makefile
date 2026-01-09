APP_NAME=app
GO=go
MIGRATE=sql-migrate
ENV=development

.PHONY: run migrate-authors migrate-articles migrate-all

dev:
	$(GO) run main.go

migrate-authors:
	$(MIGRATE) up -config=domain/authors/dbconfig.yml -env=$(ENV)

migrate-articles:
	$(MIGRATE) up -config=domain/articles/dbconfig.yml -env=$(ENV)

migrate-all: migrate-authors migrate-articles
