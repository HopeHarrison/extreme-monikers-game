package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"monikers/migrations"
	pb "monikers/proto"

	"cloud.google.com/go/bigquery"
	petname "github.com/dustinkirkland/golang-petname"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

const (
	bigqueryScope      = "https://www.googleapis.com/auth/bigquery"
	cloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"
)

type server struct {
	pb.UnimplementedMonikersServer
	bqClient *bigquery.Client
}

func (s *server) NewGame(ctx context.Context, req *pb.NewGameRequest) (*pb.NewGameResponse, error) {
	// Generate a unique ID using 4 random words
	gameID := petname.Generate(4, "-")

	// Create an insert query
	inserter := s.bqClient.Dataset("monikers").Table("games").Inserter()

	// Define the row structure
	type GameRow struct {
		ID      string    `bigquery:"id"`
		Creator string    `bigquery:"creator"`
		Created time.Time `bigquery:"created"`
	}

	// Create the row
	items := []*GameRow{
		{
			ID:      gameID,
			Creator: req.PlayerName,
			Created: time.Now(),
		},
	}

	// Insert the row
	if err := inserter.Put(ctx, items); err != nil {
		log.Printf("Failed to insert game: %v", err)
		return nil, fmt.Errorf("failed to create game: %v", err)
	}

	return &pb.NewGameResponse{GameId: gameID}, nil
}

func main() {
	ctx := context.Background()

	// Load credentials with explicit scope
	credsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	credsBytes, err := os.ReadFile(credsFile)
	if err != nil {
		log.Fatalf("Failed to read credentials file: %v", err)
	}

	creds, err := google.CredentialsFromJSON(ctx, credsBytes,
		bigqueryScope,
		cloudPlatformScope,
	)
	if err != nil {
		log.Fatalf("Failed to parse credentials: %v", err)
	}

	// Initialize BigQuery client
	bqClient, err := bigquery.NewClient(ctx, "original-storm-432806-p1", option.WithCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to create BigQuery client: %v", err)
	}
	defer bqClient.Close()

	// Ensure dataset exists and apply migrations
	dataset := bqClient.Dataset("monikers")
	if _, err := dataset.Metadata(ctx); err != nil {
		if err := dataset.Create(ctx, nil); err != nil {
			log.Fatalf("Failed to create dataset: %v", err)
		}
	}

	if err := migrations.ApplyMigrations(ctx, bqClient, "monikers"); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// Start gRPC server
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

	s := grpc.NewServer()
	pb.RegisterMonikersServer(s, &server{bqClient: bqClient})

	log.Printf("Starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
