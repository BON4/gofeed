include .env
export

.PHONY: openapi
openapi: openapi_http

.PHONY: openapi_http
openapi_http:
	@./scripts/openapi-http.sh internal/accounts
	@./scripts/openapi-http.sh internal/posts

.PHONY: migrateup
migrateup:
	migrate -path sql/migrations -database "$(DB_URL)" -verbose up && migrate -path sql/migrations -database "$(TEST_DB_URL)" -verbose up

.PHONY: createdb
createdb:
	docker exec -it postgres createdb --username=root --owner=root gofeed && docker exec -it postgres createdb --username=root --owner=root gofeed_test

.PHONY: dropdb
dropdb:
	docker exec -it postgres dropdb gofeed && docker exec -it postgres dropdb gofeed_test

.PHONY: postgres
postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

.PHONY: sqlc
sqlc:
	@./scripts/sqlc.sh sql/accounts_sqlc.yaml
	@./scripts/sqlc.sh sql/posts_sqlc.yaml
