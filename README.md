# Postgres mTLS

1. `go run main.go`
2. `docker build -t postgres-mtls .`
3. `docker run --rm -p 5432:5432 -e POSTGRES_PASSWORD=password postgres-mtls`
4. `chmod 600 certs/client.key`
5. `psql "host=127.0.0.1 port=5432 user=postgres dbname=postgres sslmode=verify-full sslcert=./certs/client.crt sslkey=./certs/client.key sslrootcert=./certs/ca.crt"`
