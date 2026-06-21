package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"ecommerce/packages/config"
	"ecommerce/packages/migrate"
)

// Run: go run ./apps/api/cmd/migrate
// Scan modules/*/migrations and apply all unapplied goose .sql files.
func main() {
	cfg := config.Load()

	conn, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "open db:", err)
		os.Exit(1)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		fmt.Fprintln(os.Stderr, "ping db:", err)
		os.Exit(1)
	}

	// Ensure schemas exist before running migrations
	if err := ensureSchemas(conn); err != nil {
		fmt.Fprintln(os.Stderr, "ensure schemas:", err)
		os.Exit(1)
	}

	migs, err := migrate.LoadAll(cfg.ModulesDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "load migrations:", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d migrations\n", len(migs))
	if err := migrate.Up(conn, migs); err != nil {
		fmt.Fprintln(os.Stderr, "migrate up:", err)
		os.Exit(1)
	}
	fmt.Println("Migrations complete.")
}

var schemas = []string{
	"auth", "identity", "customer", "product", "pricing",
	"inventory", "cart", "ordering", "payment", "fulfillment",
	"promotion", "region", "notification", "brand",
}

func ensureSchemas(conn *sql.DB) error {
	_, _ = conn.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	for _, s := range schemas {
		if _, err := conn.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %q`, s)); err != nil {
			return err
		}
	}
	return nil
}
