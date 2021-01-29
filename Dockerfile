FROM golang:1.15-alpine

WORKDIR   /code
COPY    . /code

RUN apk add --no-cache git gcc g++ sqlite
RUN go build -mod=vendor -ldflags "-w" -o /app /code/cmd/main.go

FROM alpine:3.12
RUN apk add --no-cache ca-certificates
COPY --from=0 /app        /app
COPY          config.yaml /config.yaml

ENTRYPOINT ["/app"]
