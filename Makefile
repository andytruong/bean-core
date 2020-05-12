run:
	CONFIG=config.yaml go run cmd/main.go http-server

generate-graphql:
	gqlgen generate
