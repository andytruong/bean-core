test:
	go test -mod=vendor ./pkg/... -race -v

server:
	CONFIG=config.yaml go run -mod=vendor cmd/main.go http-server

gql:
	gqlgen generate
	rm pkg/infra/__tmp__resolvers.go

migrate:
	CONFIG=config.yaml go run -mod=vendor cmd/main.go migrate

clean:
	go mod vendor
	go fmt ./...
	go mod tidy
	git fetch --prune origin
	rm -rf ./pkg/infra/gql/__tmp__*






















# ---------------------
# dev commands
# ---------------------
check-size:
	go build -mod=vendor -o /tmp/go-bean cmd/main.go
	du -h /tmp/go-bean
	rm /tmp/go-bean

dev-migrate:
	CONFIG=config.yaml DB_DRIVER=postgres DB_MASTER_URL=postgres://postgres:and1bean@127.0.0.1/core?sslmode=disable \
		go run cmd/main.go migrate

docker-db-start:
	docker run --rm -d --name=hi-pg -p 5432:5432 -e "POSTGRES_PASSWORD=and1bean" postgres:12-alpine

docker-db-stop:
	docker stop hi-pg

dev-server:
	CONFIG=config.yaml DB_DRIVER=postgres DB_MASTER_URL=postgres://postgres:and1bean@127.0.0.1/core?sslmode=disable \
		go run cmd/main.go http-server

gen-key:
	CONFIG=config.yaml go run cmd/main.go gen-key

gen-ulid:
	go run cmd/tools/ulid/main.go
