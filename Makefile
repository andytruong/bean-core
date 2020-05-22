server:
	CONFIG=config.yaml go run cmd/main.go http-server

gql:
	gqlgen generate
	rm pkg/infra/resolvers.go

migrate: gql
	CONFIG=config.yaml go run cmd/main.go migrate

test:
	go test ./... -v

build:
	go build -o /tmp/go-bean cmd/main.go

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
