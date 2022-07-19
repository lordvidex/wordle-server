migrateup: |
	migrate -path internal/db/pg/migration -database "$(DATABASE_URL)" -verbose up
migratedown: |
	migrate -path internal/db/pg/migration -database "$(DATABASE_URL)" -verbose down
mockgen: |
	mockgen -source=./internal/auth/ports.go -destination=./internal/auth/mock_ports.go -package=auth;
	
	