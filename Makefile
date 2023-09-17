postgres:
	docker run --name postgres15.4 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15.4-alpine

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
 
.PHONY:
	postgres createdb dropdb migrateup migratedown sqlc dbstart test psql