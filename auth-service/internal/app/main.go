package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"newgolang/auth-service/internal/controller"
	"newgolang/auth-service/internal/repository"
	"newgolang/auth-service/pkg/logger"
	"newgolang/proto/pb"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	_ "google.golang.org/protobuf/types/known/timestamppb"

	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	grpcAddress = "localhost:50053"
	httpAddress = "localhost:8083"
)

func main() {
	logger.InitLogger()

	ctx := context.Background()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		if err := startGrpcServer(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()

		if err := startHttpServer(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}

func startGrpcServer() error {
	log := logger.GetLogger()

	dbConn, err := ConnectDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer dbConn.Close()

	userRepo := repository.NewUserRepository(dbConn)

	userHandler := controller.NewUserHandler(*userRepo)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, userHandler)

	reflection.Register(grpcServer)

	list, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return err
	}

	log.Printf("gRPC server listening at %v", grpcAddress)

	return grpcServer.Serve(list)
}

func startHttpServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	err1 := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err1 != nil {
		return err1
	}

	log.Printf("HTTP server listening at %v", httpAddress)

	return http.ListenAndServe(httpAddress, mux)
}

func ConnectDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		logger.GetLogger().Printf("Error loading .env file: %v", err)
	}

	connectionString := os.Getenv("PG_URL")
	if connectionString == "" {
		return nil, fmt.Errorf("DB_CONNECTION_STRING environment variable not set")
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open the database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	return db, nil
}
