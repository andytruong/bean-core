migrate:
	CONFIG=config.yaml go run cmd/main.go migrate

run:
	CONFIG=config.yaml go run cmd/main.go http-server

generate-graphql:
	gqlgen generate

# ---------------------
# go mod commands
# ---------------------
tidy:
	go mod tidy

# ---------------------
# dev commands
# ---------------------
dev-migrate:
	CONFIG=config.yaml \
		DB_DRIVER=postgres \
		DB_MASTER_URL=postgres://postgres:and1bean@127.0.0.1/core?sslmode=disable  \
		go run cmd/main.go migrate

dev-run:
	CONFIG=config.yaml \
		DB_DRIVER=postgres \
		DB_MASTER_URL=postgres://postgres:and1bean@127.0.0.1/core?sslmode=disable  \
		go run cmd/main.go http-server

# ---------------------
# todo
# ---------------------
# cmd: migrate
# cmd: docker build
# cmd: deploy
# cmd: help
