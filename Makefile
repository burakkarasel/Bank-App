DB_URL=postgresql://root:password@localhost:5432/bank_app?sslmode=disable

postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root bank_app
dropdb:
	docker exec -it postgres12 dropdb bank_app
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
up:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 
down:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
server :
	go run cmd/main/main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/burakkarasel/Bank-App/db/sqlc Store
db_docs:
	dbdocs build doc/db.dbml
dc_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml
proto:
	rm -f doc/swagger/*.swagger.json
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=bank_app \
    proto/*.proto
evans:
	evans --host localhost --port 9090 -r repl
.PHONY: postgres createdb dropdb sqlc test up down server mock db_docs db_schema proto evans