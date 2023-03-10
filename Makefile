# Docker commands
rm_containers:
	docker rm $(docker ps -aq) -f

# Development commands
init_migrations:
	migrate create -ext sql -dir database/migrations -seq schema
generate_sqlc:
	(cd ./internal/database/config;sqlc generate)

# Database commands
create_db:
	docker exec -it flugo-db createdb --username=abc_valera --owner=abc_valera flugo
drop_db:
	docker exec -it flugo-db dropdb flugo
migrate_up:
	migrate -path database/migration -database "postgresql://abc_valera:abc_valera@localhost:5432/flugo?sslmode=disable" -verbose up
migrate_down:
	migrate -path database/migration -database "postgresql://abc_valera:abc_valera@localhost:5432/flugo?sslmode=disable" -verbose down

# Docker commands
build_flugo-api_image:
	docker build -t flugo:latest .

# Go commands
build_flugo-api:
	go build -o build/flugo cmd/api/main.go
test:
	go test -cover -v -cover ./...

# Run commands
run_flugo-db:
	docker run --name flugo-db --network flugo-net -p "5432:5432" -e POSTGRES_USER=abc_valera -e POSTGRES_PASSWORD=abc_valera -e POSTGRES_DB=flugo -d postgres
run_flugo-api_local:
	go build -o flugo cmd/api/main.go;./flugo
run_flugo-api:
	docker run --rm --name flugo --network flugo-net -p 3000:3000 -e DATABASE_URL="postgresql://abc_valera:abc_valera@flugo-db:5432/flugo?sslmode=disable" flugo:latest
run_flugo-all:
	docker compose -f ./deploy/docker-compose.yml up

# Deploy commands
fly_deploy:
	fly deploy -a flugo-api -c ./api.env --dockerfile ./deploy/Dockerfile