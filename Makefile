APP_NAME=app
GO=go
MIGRATE=sql-migrate
ENV=development

.PHONY: run migrate-authors migrate-articles migrate-all

dev:
	$(GO) run main.go

swagger:
	swag init -g main.go -o ./docs

migrate-authors:
	$(MIGRATE) up -config=domain/authors/dbconfig.yml -env=$(ENV)

migrate-articles:
	$(MIGRATE) up -config=domain/articles/dbconfig.yml -env=$(ENV)

migrate-all: migrate-authors migrate-articles

rollback-authors:
	$(MIGRATE) down -config=domain/authors/dbconfig.yml -env=$(ENV)

rollback-articles:
	$(MIGRATE) down -config=domain/articles/dbconfig.yml -env=$(ENV)
