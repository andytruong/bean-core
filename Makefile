server: migrate
	CONFIG=config.yaml go run cmd/main.go http-server

gql:
	gqlgen generate

migrate: gql
	CONFIG=config.yaml go run cmd/main.go migrate

test:
	go test ./... -v

clean:
	go mod tidy

# ---------------------
# dev commands
# ---------------------
dev-migrate:
	CONFIG=config.yaml \
		DB_DRIVER=postgres \
		DB_MASTER_URL=postgres://postgres:and1bean@127.0.0.1/core?sslmode=disable  \
		go run cmd/main.go migrate

dev-server:
	CONFIG=config.yaml \
		DB_DRIVER=postgres \
		DB_MASTER_URL=postgres://postgres:and1bean@127.0.0.1/core?sslmode=disable  \
		go run cmd/main.go http-server
