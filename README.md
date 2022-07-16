# Client
This part contains a cli command tool to test out some of the functionality of the game. It was quickly built to provide a simple interface to test out some of the game functionalities.
## Run
- `go run ./cmd/client/main.go`

# Server
## Run
- load .env file in the root directory of the project. An example is provided at `.env.example`
- run `docker-compose up -d`
- open http://localhost:8080

## Migrations
- download [go-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- run the below command
```bash
EXPORT DATABASE_URL=postgres://user:pass@host:port/db?sslmode=false # your database url
make migrateup      # to migrate up
make migratedown    # to migrate down
```

## Notes
* When connecting using docker, the POSTGRES_HOST is set to `pg-dev` 
