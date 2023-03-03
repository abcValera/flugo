# Docker commands
rm_containers:
	docker rm $(docker ps -aq) -f

# Development commands
init_migrations:
	migrate create -ext sql -dir database/migrations -seq schema
generate_sqlc:
	(cd ./internal/database/config;sqlc generate)
build_flugo-api_image:
	docker build -f ./deploy/Dockerfile -t flugo:latest .

# Go commands
build:
	go build -o flugo cmd/api/main.go
test:
	go test -cover -v -cover ./...
flugo-api_local:
	go build -o flugo cmd/api/main.go;./flugo

# Database commands
run_flugo-db:
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
run_flugo-api:
	docker run --rm --name flugo --network flugo-net -p 3000:3000 -e DATABASE_URL="postgresql://abc_valera:abc_valera@flugo-db:5432/flugo?sslmode=disable" flugo:latest
run_flugo-all:
	docker compose -f ./deploy/docker-compose.yml up