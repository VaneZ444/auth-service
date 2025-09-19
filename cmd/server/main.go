package main

import (
	"database/sql"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/VaneZ444/auth-service/internal/handler"
	"github.com/VaneZ444/auth-service/internal/jwt"
	pgRepo "github.com/VaneZ444/auth-service/internal/repository/postgres"
	"github.com/VaneZ444/auth-service/internal/usecase"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	logger "github.com/VaneZ444/forum-shared/logger"
	protos "github.com/VaneZ444/golang-forum-protos/gen/go/sso"
	"google.golang.org/grpc"
)

func main() {
	log := logger.New(slog.LevelDebug)

	// Подключение к PostgreSQL
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=3781 dbname=auth sslmode=disable")
	if err != nil {
		log.Error("DB connection failed", logger.Err(err))
		panic(err)
	}
	defer db.Close()

	// Применение миграций
	if err := applyMigrations(db); err != nil {
		log.Error("Migrations failed", logger.Err(err))
		panic(err)
	}

	// Инициализация репозиториев
	userRepo := pgRepo.NewUserRepository(db, log)
	appRepo := pgRepo.NewAppRepository(db, log)
	jwtService := jwt.NewService("your-secret-key", 24*time.Hour)
	// 5. Создание UseCase
	authUC := usecase.NewAuthUseCase(
		userRepo,
		appRepo,
		jwtService,
		log,
	)

	// 6. Создание gRPC обработчика
	authHandler := handler.NewAuthHandler(authUC, jwtService, log)

	// 7. Запуск gRPC сервера
	server := grpc.NewServer()
	protos.RegisterAuthServer(server, authHandler)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Error("Server failed", logger.Err(err))
		os.Exit(1)
	}

	log.Info("Auth service started on :50051")
	if err := server.Serve(lis); err != nil {
		log.Error("Server failed", logger.Err(err))
		os.Exit(1)
	}
}

func applyMigrations(db *sql.DB) error {
	driver, err := migratePostgres.WithInstance(db, &migratePostgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../internal/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
