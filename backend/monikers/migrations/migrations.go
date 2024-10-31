package migrations

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
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
			id STRING,
			creator STRING,
			created TIMESTAMP
		)`,
	},
	// Add new migrations here
}

func ApplyMigrations(ctx context.Context, client *bigquery.Client, dataset string) error {
	// Ensure schema_migrations table exists
	migrationsTable := client.Dataset(dataset).Table("schema_migrations")

	// Check if table exists
	if _, err := migrationsTable.Metadata(ctx); err != nil {
		// Create schema_migrations table if it doesn't exist
		schema := bigquery.Schema{
			{Name: "version", Type: bigquery.IntegerFieldType, Required: true},
			{Name: "applied_at", Type: bigquery.TimestampFieldType, Required: true},
		}

		tableMetadata := &bigquery.TableMetadata{
			Schema: schema,
		}

		if err := migrationsTable.Create(ctx, tableMetadata); err != nil {
			return fmt.Errorf("failed to create migration table: %v", err)
		}
	}

	// Check which migrations have been applied
	for _, migration := range Migrations {
		// Use the SQL directly, no formatting needed
		q := client.Query(migration.SQL)
		job, err := q.Run(ctx)
		if err != nil {
			return fmt.Errorf("failed to apply migration %d: %v", migration.Version, err)
		}
		// Wait for the query to complete
		status, err := job.Wait(ctx)
		if err != nil {
			return fmt.Errorf("failed to wait for migration %d: %v", migration.Version, err)
		}
		if err := status.Err(); err != nil {
			return fmt.Errorf("migration %d failed: %v", migration.Version, err)
		}

		// Record that migration was applied
		inserter := client.Dataset(dataset).Table("schema_migrations").Inserter()
		row := struct {
			Version   int64     `bigquery:"version"`
			AppliedAt time.Time `bigquery:"applied_at"`
		}{
			Version:   int64(migration.Version),
			AppliedAt: time.Now(),
		}
		if err := inserter.Put(ctx, []interface{}{row}); err != nil {
			return fmt.Errorf("failed to record migration: %v", err)
		}

		log.Printf("Migration %d applied successfully", migration.Version)
	}
	return nil
}
