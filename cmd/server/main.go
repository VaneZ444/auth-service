package main

import (
	"database/sql"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/VaneZ444/auth-service/internal/handler"
	"github.com/VaneZ444/auth-service/internal/repository/postgres"
	"github.com/VaneZ444/auth-service/internal/usecase"

	logger "github.com/VaneZ444/forum-shared/logger"
	protos "github.com/VaneZ444/golang-forum-protos/gen/go/sso"
	"google.golang.org/grpc"
)

func main() {
	// 1. Инициализация логгера
	log := logger.New(slog.LevelDebug) // или LevelProd для прода

	// 2. Подключение к PostgreSQL
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=auth sslmode=disable")
	if err != nil {
		log.Error("failed to connect to DB", logger.Err(err))
		panic(err)
	}
	defer db.Close()

	// 3. Инициализация репозиториев
	userRepo := postgres.NewUserRepository(db)
	appRepo := postgres.NewAppRepository(db) // Реализация AppRepository

	// 4. Создание UseCase
	authUC := usecase.NewAuthUseCase(
		userRepo,
		appRepo,           // Реализация интерфейса AppRepository
		"your-secret-key", // Секрет для JWT
		24*time.Hour,      // TTL токена
		log,
	)

	// 5. Создание gRPC обработчика
	authHandler := handler.NewAuthHandler(authUC, log)

	// 6. Запуск gRPC сервера
	srv := grpc.NewServer()
	protos.RegisterAuthServer(srv, authHandler)
	// 3. gRPC-сервер
	server := grpc.NewServer()
	protos.RegisterAuthServer(server, authHandler)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Error("Server failed", logger.Err(err))
		os.Exit(1)
	}
	log.Info("Auth service started on :50051")
	server.Serve(lis)
}
