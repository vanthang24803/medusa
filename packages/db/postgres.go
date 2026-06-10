package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config contains database connection configuration parameters.
type Config struct {
	Driver    string // "postgres", "mysql", "mssql" (default "postgres")
	MasterDSN string
	SlaveDSN  string // empty -> use MasterDSN
}

// Querier — common interface for both *sqlx.DB and *sqlx.Tx, for repository
// can be used both inside and outside transactions.
type Querier interface {
	sqlx.QueryerContext
	sqlx.ExecerContext
	sqlx.ExtContext
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
}

// Engine represents the underlying database driver (PostgreSQL, MySQL, MSSQL, etc.).
type Engine interface {
	Writer(ctx context.Context) Querier
	Reader(ctx context.Context) Querier
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	Close() error
}

// DB wraps the database Engine to provide a unified interface for Repositories.
type DB struct {
	engine Engine
}

// Writer returns the Querier for writing data (supports transaction).
func (d *DB) Writer(ctx context.Context) Querier {
	return d.engine.Writer(ctx)
}

// Reader returns the Querier for reading data (supports master/slave and transaction).
func (d *DB) Reader(ctx context.Context) Querier {
	return d.engine.Reader(ctx)
}

// WithTransaction executes logic inside a safe transaction.
func (d *DB) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.engine.WithTransaction(ctx, fn)
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.engine.Close()
}

// Connect initializes the database connection matching the Driver config.
func Connect(cfg Config) (*DB, error) {
	if cfg.Driver == "" {
		cfg.Driver = "postgres"
	}

	switch cfg.Driver {
	case "postgres":
		engine, err := newPostgresEngine(cfg)
		if err != nil {
			return nil, err
		}
		return &DB{engine: engine}, nil
	case "mysql":
		engine, err := newMySQLEngine(cfg)
		if err != nil {
			return nil, err
		}
		return &DB{engine: engine}, nil
	case "mssql", "sqlserver":
		engine, err := newMSSQLEngine(cfg)
		if err != nil {
			return nil, err
		}
		return &DB{engine: engine}, nil
	default:
		return nil, fmt.Errorf("db: unsupported driver: %s", cfg.Driver)
	}
}

// MustConnect initializes the database connection and panics if an error occurs.
func MustConnect(cfg Config) *DB {
	d, err := Connect(cfg)
	if err != nil {
		panic(err)
	}
	return d
}

// ── PostgreSQL Implementation ────────────────────────────────────────────────

type txKey struct{}

type postgresEngine struct {
	Write *sqlx.DB
	Read  *sqlx.DB
}

func newPostgresEngine(cfg Config) (*postgresEngine, error) {
	openDB := func(dsn string) (*sqlx.DB, error) {
		d, err := sqlx.Connect("postgres", dsn)
		if err != nil {
			return nil, err
		}
		d.SetMaxOpenConns(25)
		d.SetMaxIdleConns(10)
		d.SetConnMaxLifetime(5 * time.Minute)
		return d, nil
	}

	write, err := openDB(cfg.MasterDSN)
	if err != nil {
		return nil, fmt.Errorf("connect master: %w", err)
	}

	read := write
	if cfg.SlaveDSN != "" && cfg.SlaveDSN != cfg.MasterDSN {
		read, err = openDB(cfg.SlaveDSN)
		if err != nil {
			_ = write.Close()
			return nil, fmt.Errorf("connect slave: %w", err)
		}
	}

	if err := ensurePostgresSchemas(write); err != nil {
		_ = write.Close()
		if read != write {
			_ = read.Close()
		}
		return nil, fmt.Errorf("ensure schemas: %w", err)
	}

	return &postgresEngine{Write: write, Read: read}, nil
}

func (e *postgresEngine) Close() error {
	_ = e.Write.Close()
	if e.Read != e.Write {
		_ = e.Read.Close()
	}
	return nil
}

func (e *postgresEngine) Writer(ctx context.Context) Querier {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return e.Write
}

func (e *postgresEngine) Reader(ctx context.Context) Querier {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return e.Read
}

func (e *postgresEngine) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	tx, err := e.Write.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	ctx = context.WithValue(ctx, txKey{}, tx)
	return fn(ctx)
}

// Postgres Schemas
var postgresSchemas = []string{
	"auth", "identity", "customer", "product", "pricing",
	"inventory", "cart", "ordering", "payment", "fulfillment",
	"promotion", "region", "notification",
}

func ensurePostgresSchemas(d *sqlx.DB) error {
	_, _ = d.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	for _, s := range postgresSchemas {
		if _, err := d.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %q`, s)); err != nil {
			return err
		}
	}
	return nil
}

// ── MySQL Implementation (Stub) ──────────────────────────────────────────────

type mysqlEngine struct {
	// Declare MySQL connection here when implementing
}

func newMySQLEngine(cfg Config) (*mysqlEngine, error) {
	// TODO: Implement actual MySQL connection
	return nil, fmt.Errorf("db: mysql driver is not implemented yet")
}

func (e *mysqlEngine) Close() error {
	return nil
}

func (e *mysqlEngine) Writer(ctx context.Context) Querier {
	return nil
}

func (e *mysqlEngine) Reader(ctx context.Context) Querier {
	return nil
}

func (e *mysqlEngine) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return nil
}

// ── MSSQL Implementation (Stub) ──────────────────────────────────────────────

type mssqlEngine struct {
	// Declare MSSQL connection here when implementing
}

func newMSSQLEngine(cfg Config) (*mssqlEngine, error) {
	// TODO: Implement actual MSSQL connection
	return nil, fmt.Errorf("db: mssql driver is not implemented yet")
}

func (e *mssqlEngine) Close() error {
	return nil
}

func (e *mssqlEngine) Writer(ctx context.Context) Querier {
	return nil
}

func (e *mssqlEngine) Reader(ctx context.Context) Querier {
	return nil
}

func (e *mssqlEngine) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return nil
}
