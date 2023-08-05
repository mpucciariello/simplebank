postgres:
	docker run --name postgres12 --network bank_network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

format:
	go fmt ./...

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go github.com/micaelapucciariello/simplebank/db/sqlc Store

proto:
	rm -f pb/*.go
	rm -f docs/swagger/*.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true \
		proto/*.proto
		statik -src=./docs/swagger -dest=./docs

evans:
	evans --host localhost --port 9091 -r repl

linter:
	golangci-lint run ./...


.PHONY: postgres createdb dropdb migrateup migratedown format sqlc test server mock proto evans
