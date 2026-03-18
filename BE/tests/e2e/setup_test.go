package e2e_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultAppURL = "http://localhost:8080"
	defaultPGConn = "postgresql://postgres:postgres@localhost:5432/whatdish?sslmode=disable"
	fakeAuthAddr  = ":9999"
)

// e2eInfra holds references needed by the E2E test suite.
// Unlike testcontainers, the services are already running via docker-compose.
type e2eInfra struct {
	fakeAuth   *fakeAuthServer
	appBaseURL string
	pgConnStr  string
}

// setupInfra connects to the already-running docker-compose services,
// runs migrations + seeds, and starts the fake Supabase auth server.
func setupInfra(ctx context.Context) *e2eInfra {
	infra := &e2eInfra{}

	infra.appBaseURL = envOr("APP_BASE_URL", defaultAppURL)
	infra.pgConnStr = envOr("DATABASE_URL", defaultPGConn)

	// 1. Run migrations (idempotent).
	runMigrations(infra.pgConnStr)

	// 2. Run seed data (idempotent — uses ON CONFLICT DO NOTHING).
	runSeeds(ctx, infra.pgConnStr)

	// 3. Start fake Supabase auth on a fixed port so the app container
	//    can reach it via host.docker.internal:9999.
	infra.fakeAuth = startFakeAuthServerOnAddr(fakeAuthAddr)
	log.Printf("e2e: fake auth running on %s", fakeAuthAddr)

	// 4. Wait for the app to be healthy.
	waitForHealth(infra.appBaseURL, 60*time.Second)

	log.Printf("e2e: app ready at %s", infra.appBaseURL)
	return infra
}

func (infra *e2eInfra) teardown() {
	if infra.fakeAuth != nil {
		infra.fakeAuth.Close()
	}
}

func runMigrations(connStr string) {
	migrationsPath := filepath.Join(projectRootDir(), "migrations")
	m, err := migrate.New("file://"+migrationsPath, connStr)
	if err != nil {
		log.Fatalf("e2e: create migrator: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("e2e: run migrations: %v", err)
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		log.Fatalf("e2e: close migrator source: %v", srcErr)
	}
	if dbErr != nil {
		log.Fatalf("e2e: close migrator db: %v", dbErr)
	}
}

func runSeeds(ctx context.Context, connStr string) {
	seedDir := filepath.Join(projectRootDir(), "seeds")
	seedFiles := []string{
		"01_categories.sql",
		"02_tags.sql",
		"03_nutrition_goals.sql",
		"04_suggestion_configs.sql",
	}

	pool, err := pgxPoolFromConnStr(ctx, connStr)
	if err != nil {
		log.Fatalf("e2e: seed pool: %v", err)
	}
	defer pool.Close()

	for _, f := range seedFiles {
		path := filepath.Join(seedDir, f)
		sql, err := readFileContent(path)
		if err != nil {
			log.Printf("e2e: skip seed %s: %v", f, err)
			continue
		}
		if _, err := pool.Exec(ctx, sql); err != nil {
			log.Fatalf("e2e: run seed %s: %v", f, err)
		}
		log.Printf("e2e: seeded %s", f)
	}
}

func waitForHealth(baseURL string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 2 * time.Second}
	for time.Now().Before(deadline) {
		resp, err := client.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(500 * time.Millisecond)
	}
	log.Fatalf("e2e: app did not become healthy within %v", timeout)
}

func projectRootDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
