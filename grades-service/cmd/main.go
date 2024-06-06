package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"newgolang/auth-service/pkg/logger"
	"newgolang/grades-service/internal"

	"newgolang/proto/pb"
	"os"
	"sync"

	_ "github.com/lib/pq"
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

	gradesRepo := internal.NewGradeRepository(dbConn)

	gradesHandler := internal.NewGradeHandler(*gradesRepo)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterGradeServiceServer(grpcServer, gradesHandler)

	reflection.Register(grpcServer)

	list, err := net.Listen("tcp", "localhost:50057")
	if err != nil {
		return err
	}

	log.Printf("gRPC server listening at %v", "localhost:50057")

	return grpcServer.Serve(list)
}

func startHttpServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	err1 := pb.RegisterGradeServiceHandlerFromEndpoint(ctx, mux, "localhost:50057", opts)
	if err1 != nil {
		return err1
	}

	log.Printf("HTTP server listening at %v", "localhost:8087")

	return http.ListenAndServe("localhost:8087", mux)
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
