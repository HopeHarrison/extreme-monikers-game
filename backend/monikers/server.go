package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"monikers/migrations"
	pb "monikers/proto"

	// Import the MySQL driver with Cloud SQL connector
	petname "github.com/dustinkirkland/golang-petname"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMonikersServer
	db *sql.DB
}

func (s *server) NewGame(ctx context.Context, req *pb.NewGameRequest) (*pb.NewGameResponse, error) {
	// Generate a unique ID using 4 random words
	gameID := petname.Generate(4, "-")

	// Get the current time in UTC
	now := time.Now().UTC()

	// Insert the game row
	gameInsertSQL := `
		INSERT INTO games (id, timeCreated)
		VALUES (?, ?)`
	_, err := s.db.Exec(gameInsertSQL, gameID, now)
	if err != nil {
		log.Printf("Failed to insert game: %v", err)
		return nil, fmt.Errorf("failed to create game: %v", err)
	}

	// Insert the player row
	playerInsertSQL := `
		INSERT INTO players (gameId, playerName, team, timeJoined)
		VALUES (?, ?, ?, ?)`
	_, err = s.db.Exec(playerInsertSQL, gameID, req.PlayerName, 0, now)
	if err != nil {
		log.Printf("Failed to insert player: %v", err)
		return nil, fmt.Errorf("failed to add player: %v", err)
	}

	return &pb.NewGameResponse{GameId: gameID}, nil
}

func main() {
	ctx := context.Background()

	// Get Cloud SQL connection details from environment variables
	projectID := os.Getenv("PROJECT_ID")
	region := os.Getenv("REGION")
	instanceName := os.Getenv("INSTANCE_NAME")

	if projectID == "" || region == "" || instanceName == "" {
		log.Fatal("PROJECT_ID, REGION, and INSTANCE_NAME environment variables are required")
	}

	instanceConnName := fmt.Sprintf("%s:%s:%s", projectID, region, instanceName)

	// Check if running locally or in Cloud Run
	var dsn string
	if os.Getenv("K_SERVICE") == "" {
		// Local development - use TCP with Docker host
		dsn = fmt.Sprintf("%s@tcp(host.docker.internal:3306)/%s", "monikers-app", "monikers")
	} else {
		// Cloud Run - use Unix socket with IAM auth
		dsn = fmt.Sprintf("%s@unix(/cloudsql/%s)/%s?allowCleartextPasswords=1&auth_plugin=mysql_clear_password",
			"monikers-app",   // IAM user
			instanceConnName, // project:region:instance
			"monikers")       // database name
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Add connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database!")

	// Run migrations
	if err := migrations.ApplyMigrations(ctx, db, instanceName); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// Create and start the gRPC server
	s := grpc.NewServer()
	pb.RegisterMonikersServer(s, &server{db: db})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Setting up gRPC server on %s", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
