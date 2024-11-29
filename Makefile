include .envrc

.PHONY: run/api
run/api:
	@echo 'Running book club API...'
	@go run ./cmd/api -port=3000 -env=production -db-dsn=${BOOKCLUB_DB_DSN}

.PHONY: db/psql
db/psql:
	psql ${BOOKCLUB_DB_DSN}

.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path=./migrations -database ${BOOKCLUB_DB_DSN} up