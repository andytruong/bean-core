Local environment
====

```
cd /path/to/bean-core

export DB_DRIVER="postgres"
export DB_MASTER_URL="postgres://postgres:and1bean@127.0.0.1/bean-core?sslmode=disable"

docker run -d \
    --name=hi-pg \
    -p 5432:5432 \
    -e "POSTGRES_PASSWORD=and1bean" \
    -v `pwd`'/data/postgres:/var/lib/postgresql/data' \
    postgres:12-alpine

make dev-migrate
```
