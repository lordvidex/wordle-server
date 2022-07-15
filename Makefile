migrateup: |
	migrate -path internal/db/pg/migration -database "$(DATABASE_URL)" -verbose up
migratedown: |
	migrate -path internal/db/pg/migration -database "$(DATABASE_URL)" -verbose down

	