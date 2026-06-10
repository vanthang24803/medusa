package migrate

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Migration represents 1 .sql file in goose format.
type Migration struct {
	Version string
	Name    string
	Path    string
	Up      string
	Down    string
}

// LoadAll scans modules/*/migrations and parses goose .sql files.
// This is a lightweight runner compatible with goose format (-- +goose Up / Down,
// StatementBegin/End). The .sql files can still run using the real goose CLI.
func LoadAll(modulesDir string) ([]Migration, error) {
	var migs []Migration
	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		migDir := filepath.Join(modulesDir, e.Name(), "migrations")
		files, err := os.ReadDir(migDir)
		if err != nil {
			continue // module has no migrations yet
		}
		for _, f := range files {
			if !strings.HasSuffix(f.Name(), ".sql") {
				continue
			}
			content, err := os.ReadFile(filepath.Join(migDir, f.Name()))
			if err != nil {
				return nil, err
			}
			up, down := parseGoose(string(content))
			// version = module + filename to be unique and sort stably
			version := e.Name() + "/" + f.Name()
			migs = append(migs, Migration{
				Version: version,
				Name:    f.Name(),
				Path:    filepath.Join(migDir, f.Name()),
				Up:      up,
				Down:    down,
			})
		}
	}
	sort.Slice(migs, func(i, j int) bool {
		return migs[i].Version < migs[j].Version
	})
	return migs, nil
}

// parseGoose separates Up and Down, removing comment annotations.
func parseGoose(content string) (up, down string) {
	lines := strings.Split(content, "\n")
	var section string
	var upB, downB strings.Builder
	for _, ln := range lines {
		trimmed := strings.TrimSpace(ln)
		switch {
		case strings.HasPrefix(trimmed, "-- +goose Up"):
			section = "up"
			continue
		case strings.HasPrefix(trimmed, "-- +goose Down"):
			section = "down"
			continue
		case strings.HasPrefix(trimmed, "-- +goose"):
			continue // StatementBegin/End markers
		}
		if section == "up" {
			upB.WriteString(ln + "\n")
		} else if section == "down" {
			downB.WriteString(ln + "\n")
		}
	}
	return strings.TrimSpace(upB.String()), strings.TrimSpace(downB.String())
}

// EnsureTable creates a table to track applied migrations.
func EnsureTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS public.schema_migrations (
			version    VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`)
	return err
}

// Applied returns the set of applied migration versions.
func Applied(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query(`SELECT version FROM public.schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	set := map[string]bool{}
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		set[v] = true
	}
	return set, rows.Err()
}

// Up runs all unapplied migrations in order, each within a transaction.
func Up(db *sql.DB, migs []Migration) error {
	if err := EnsureTable(db); err != nil {
		return err
	}
	applied, err := Applied(db)
	if err != nil {
		return err
	}
	for _, m := range migs {
		if applied[m.Version] {
			continue
		}
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(m.Up); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migration %s failed: %w", m.Version, err)
		}
		if _, err := tx.Exec(
			`INSERT INTO public.schema_migrations (version) VALUES ($1)`, m.Version); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		fmt.Printf("  ✓ applied %s\n", m.Version)
	}
	return nil
}
