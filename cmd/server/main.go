package main

import (
	"database/sql"
	"log"
	"net"
	"time"

	"github.com/VaneZ444/golang-forum/auth-service/internal/handler"
	"github.com/VaneZ444/golang-forum/auth-service/internal/repository/postgres"
	"github.com/VaneZ444/golang-forum/auth-service/internal/usecase"

	ssov1 "github.com/VaneZ444/golang-forum/tree/main/shared/protos/gen/go/sso"
	"google.golang.org/grpc"
)

func main() {
	// 1. Подключение к PostgreSQL
	db, err := sql.Open("postgres", "your_connection_string")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. Инициализация слоев
	repo := postgres.NewPostgresRepo(db)
	authUC := usecase.NewAuthUseCase(repo, "your-secret-key", 24*time.Hour)
	authHandler := handler.NewAuthHandler(authUC)

	// 3. gRPC-сервер
	server := grpc.NewServer()
	protos.RegisterAuthServer(server, authHandler)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Auth service started on :50051")
	server.Serve(lis)
}
