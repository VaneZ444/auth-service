# auth-service
# Применить все миграции
migrate -path ./migrations -database "postgres://user:password@localhost:5432/auth_db?sslmode=disable" up

# Откатить последнюю миграцию
migrate -path ./migrations -database "postgres://..." down 1

# Проверить текущую версию
migrate -path ./migrations -database "postgres://..." version

migrate -path ./internal/migrations -database "postgres://login:password@localhost:5432/auth?sslmode=disable" up