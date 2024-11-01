package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Migration represents a database change
type Migration struct {
	Version int
	SQL     string
}

// All migrations in order
var Migrations = []Migration{
	{
		Version: 1,
		SQL: `
		CREATE TABLE IF NOT EXISTS monikers.games (
			id VARCHAR(255) NOT NULL PRIMARY KEY,
			timeCreated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS monikers.players (
			gameId VARCHAR(255),
			playerName VARCHAR(255) NOT NULL UNIQUE,
			team INT,
			timeJoined TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (gameId) REFERENCES monikers.games(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS monikers.cards (
			id VARCHAR(255) NOT NULL PRIMARY KEY,
			title VARCHAR(255) NOT NULL UNIQUE,
			description TEXT,
			playerNameCreatedBy VARCHAR(255),
			FOREIGN KEY (playerNameCreatedBy) REFERENCES monikers.players(playerName) ON DELETE SET NULL
		);

		CREATE TABLE IF NOT EXISTS monikers.turns (
			id VARCHAR(255) NOT NULL PRIMARY KEY,
			gameId VARCHAR(255),
			roundNumber INT,
			turnNumber INT,
			timeStarted TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			playerName VARCHAR(255),
			FOREIGN KEY (gameId) REFERENCES monikers.games(id) ON DELETE CASCADE,
			FOREIGN KEY (playerName) REFERENCES monikers.players(playerName) ON DELETE SET NULL
		);

		CREATE TABLE IF NOT EXISTS monikers.wonCards (
			turnId VARCHAR(255),
			cardId VARCHAR(255),
			FOREIGN KEY (turnId) REFERENCES monikers.turns(id) ON DELETE CASCADE,
			FOREIGN KEY (cardId) REFERENCES monikers.cards(id) ON DELETE CASCADE,
			PRIMARY KEY (turnId, cardId)
		);

		CREATE TABLE IF NOT EXISTS monikers.skippedCards (
			turnId VARCHAR(255),
			cardId VARCHAR(255),
			FOREIGN KEY (turnId) REFERENCES monikers.turns(id) ON DELETE CASCADE,
			FOREIGN KEY (cardId) REFERENCES monikers.cards(id) ON DELETE CASCADE,
			PRIMARY KEY (turnId, cardId)
		);
		`,
	},
	// Add new migrations here
}

func ApplyMigrations(ctx context.Context, db *sql.DB, dataset string) error {
	// Ensure schema_migrations table exists
	migrationsTable := fmt.Sprintf("%s.schema_migrations", dataset)

	// Check if table exists
	if _, err := db.QueryContext(ctx, fmt.Sprintf("SHOW TABLES LIKE '%s'", migrationsTable)); err != nil {
		// Create schema_migrations table if it doesn't exist
		schema := `
		CREATE TABLE IF NOT EXISTS %s (
			version INT NOT NULL,
			applied_at TIMESTAMP NOT NULL
		);
		`

		if _, err := db.ExecContext(ctx, fmt.Sprintf(schema, migrationsTable)); err != nil {
			return fmt.Errorf("failed to create migration table: %v", err)
		}
	}

	// Check which migrations have been applied
	for _, migration := range Migrations {
		// Use the SQL directly, no formatting needed
		if _, err := db.ExecContext(ctx, migration.SQL); err != nil {
			return fmt.Errorf("failed to apply migration %d: %v", migration.Version, err)
		}

		// Record that migration was applied
		if _, err := db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s (version, applied_at) VALUES (%d, CURRENT_TIMESTAMP)", migrationsTable, migration.Version)); err != nil {
			return fmt.Errorf("failed to record migration: %v", err)
		}

		log.Printf("Migration %d applied successfully", migration.Version)
	}
	return nil
}
