# Development commands
init_migrations:
	migrate create -ext sql -dir database/migrations -seq schema
generate_sqlc:
	sqlc generate

# Database commands
flugo-db:
	docker run --name flugo-db --network flugo-net -p "5432:5432" -e POSTGRES_USER=abc_valera -e POSTGRES_PASSWORD=abc_valera -e POSTGRES_DB=flugo -d postgres
create_db:
	docker exec -it flugo-db createdb --username=abc_valera --owner=abc_valera flugo
drop_db:
	docker exec -it flugo-db dropdb flugo
migrate_up:
	migrate -path database/migration -database "postgresql://abc_valera:abc_valera@localhost:5432/flugo?sslmode=disable" -verbose up
migrate_down:
	migrate -path database/migration -database "postgresql://abc_valera:abc_valera@localhost:5432/flugo?sslmode=disable" -verbose down
# API commands
build_flugo-api_image:
	docker build -t flugo:latest .
flugo-api:
	docker run --rm --name flugo --network flugo-net -p 3000:3000 -e DATABASE_URL="postgresql://abc_valera:abc_valera@flugo-db:5432/flugo?sslmode=disable" flugo:latest

# Go commands
test:
	go test -cover -v -cover ./...