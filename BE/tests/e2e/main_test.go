package e2e_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
)

// infra holds the E2E infrastructure references (set in TestMain).
var infra *e2eInfra

// recipeIDs holds the IDs of seeded test recipes (set in TestMain).
var recipeIDs []uuid.UUID

func TestMain(m *testing.M) {
	ctx := context.Background()

	log.Println("e2e: setting up...")

	// Connect to already-running docker-compose services,
	// run migrations + seeds, start fake auth.
	infra = setupInfra(ctx)

	// Seed test recipes.
	recipeIDs = seedTestRecipes(ctx, infra.pgConnStr)
	log.Printf("e2e: seeded %d test recipes", len(recipeIDs))

	code := m.Run()

	log.Println("e2e: tearing down...")
	infra.teardown()

	os.Exit(code)
}

// requireE2E skips the test if the E2E infrastructure is not available.
func requireE2E(t *testing.T) {
	t.Helper()
	if infra == nil {
		t.Skip("skipping: E2E infrastructure not available")
	}
}
