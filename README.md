## Run
- load .env file
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
