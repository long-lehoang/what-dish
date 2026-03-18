package testutil

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// TestDB holds a testcontainers PostgreSQL instance and a connection pool.
type TestDB struct {
	Pool      *pgxpool.Pool
	container testcontainers.Container
}

// SetupTestDB starts a PostgreSQL container, runs migrations, and returns a TestDB.
// Intended for use in TestMain. Call Cleanup() in the defer.
func SetupTestDB() *TestDB {
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("whatdish_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("testutil.SetupTestDB: start container: %v", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("testutil.SetupTestDB: connection string: %v", err)
	}

	// Run migrations.
	migrationsPath := migrationDir()
	m, err := migrate.New("file://"+migrationsPath, connStr)
	if err != nil {
		log.Fatalf("testutil.SetupTestDB: create migrator: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("testutil.SetupTestDB: run migrations: %v", err)
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		log.Fatalf("testutil.SetupTestDB: close migrator source: %v", srcErr)
	}
	if dbErr != nil {
		log.Fatalf("testutil.SetupTestDB: close migrator db: %v", dbErr)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("testutil.SetupTestDB: create pool: %v", err)
	}

	return &TestDB{
		Pool:      pool,
		container: container,
	}
}

// Cleanup closes the pool and stops the container.
func (tdb *TestDB) Cleanup() {
	tdb.Pool.Close()
	if err := tdb.container.Terminate(context.Background()); err != nil {
		log.Printf("testutil.Cleanup: terminate container: %v", err)
	}
}

// TruncateAll truncates all application tables for test isolation.
func (tdb *TestDB) TruncateAll(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	tables := []string{
		"recipe_tags",
		"recipe_ingredients",
		"recipe_steps",
		"suggestion_exclusions",
		"suggestion_sessions",
		"suggestion_configs",
		"engagement_favorites",
		"engagement_views",
		"engagement_ratings",
		"nutrition_recipe",
		"nutrition_goals",
		"user_allergies",
		"user_profiles",
		"recipes",
		"tags",
		"categories",
	}
	for _, table := range tables {
		_, err := tdb.Pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Fatalf("testutil.TruncateAll: truncate %s: %v", table, err)
		}
	}
}

// RequireTestDB skips the test if tdb is nil (not in integration mode).
func RequireTestDB(t *testing.T, tdb *TestDB) {
	t.Helper()
	if tdb == nil {
		t.Skip("skipping integration test: no test database (run with -run Integration)")
	}
}

// migrationDir returns the absolute path to the migrations directory.
func migrationDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "..", "migrations")
}
