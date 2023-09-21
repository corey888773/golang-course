postgres:
	docker run --name postgres15.4 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=simple_bank -d postgres:15.4-alpine

dbstart:
	docker start postgres15.4

createdb:
	docker exec -it postgres15.4 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15.4 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
	
sqlc:
	sqlc generate

psql:
	docker exec -it postgres15.4 psql -U root simple_bank

test:
	go test -v -cover -short ./... -mod=mod
 
server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store_mock.go github.com/corey888773/golang-course/db/sqlc Store

.PHONY:
	postgres createdb dropdb migrateup migratedown sqlc dbstart test psql server mock