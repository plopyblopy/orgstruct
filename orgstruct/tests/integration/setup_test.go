package integration

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/plopyblopy/orgstruct/internal/shared"
	"github.com/plopyblopy/orgstruct/tests/testdata"
)

func TestSetup(t *testing.T) {
	pgContainerCtx := context.Background()
	// config.
	c, err := testdata.NewTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	// container.
	pgContainer, err := testdata.NewPostgresTestcontainer(pgContainerCtx, *c)
	if err != nil {
		t.Fatal(err)
	}

	// dynamic retrieval of the actual host, port, and assignment in the config.
	host, err := pgContainer.Host(pgContainerCtx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := pgContainer.MappedPort(pgContainerCtx, "5432")
	if err != nil {
		t.Fatal(err)
	}

	c.DbConfig.Host = host
	c.DbConfig.Port = port.Port()

	// project root from sys root.
	root, err := shared.FindProjectRoot()
	if err != nil {
		t.Fatal(err)
	}

	// path to migrate folder.
	migrationsPath := filepath.Join(root, "migrations")

	// initialize migration on db.
	err = testdata.InitMigrate(migrationsPath, *c)
	if err != nil {
		t.Fatal(err)
	}

	// wrap db conns in TestSuite.
	err = testdata.InitTestEnv(pgContainerCtx, pgContainer, c)
	if err != nil {
		t.Fatal(err)
	}

	testSuite, err := testdata.GetTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	// remove testcontainer from docker.
	t.Cleanup(func() {
		if err := testSuite.PgContainer.Terminate(pgContainerCtx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})
}
