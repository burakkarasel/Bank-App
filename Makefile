postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root bank_app
dropdb:
	docker exec -it postgres12 dropdb bank_app
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
up:
	soda migrate -p ./db
down:
	soda migrate down -p ./db
server :
	go run cmd/main/main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/burakkarasel/Bank-App/db/sqlc Store
.PHONY: postgres createdb dropdb sqlc test up down server mock