package suggestion_test

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/lehoanglong/whatdish/internal/shared/testutil"
)

var testDB *testutil.TestDB

func TestMain(m *testing.M) {
	flag.Parse()

	if isIntegrationRun() {
		testDB = testutil.SetupTestDB()
		code := m.Run()
		testDB.Cleanup()
		os.Exit(code)
	}

	os.Exit(m.Run())
}

func requireIntegrationDB(t *testing.T) {
	t.Helper()
	if testDB == nil {
		t.Skip("skipping: no test database (run with -run Integration)")
	}
	testDB.TruncateAll(t)
}

func isIntegrationRun() bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, "Integration") {
			return true
		}
	}
	return os.Getenv("TEST_INTEGRATION") == "1"
}
