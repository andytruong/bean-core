run:
	CONFIG=config.yaml go run cmd/main.go http-server

migrate:
	CONFIG=config.yaml go run cmd/main.go migrate

generate-graphql:
	gqlgen generate

# cmd: migrate
# cmd: docker build
# cmd: deploy
# cmd: help
