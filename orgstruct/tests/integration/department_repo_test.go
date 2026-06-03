package integration

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/plopyblopy/orgstruct/internal/adapters/postgres"
	"github.com/plopyblopy/orgstruct/internal/domain"
	"github.com/plopyblopy/orgstruct/internal/shared"
	"github.com/plopyblopy/orgstruct/tests/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// В этом тесте используется Testcontainer где в начале каждого теста вызывается testSuite.SetupTestPg - он закрывает все подключения, делает Rollback до начального состояния
func TestDepartmentRepository(t *testing.T) {
	pgContainerCtx := context.Background()
	// config
	c, err := testdata.NewTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	// container
	pgContainer, err := testdata.NewPostgresTestcontainer(pgContainerCtx, *c)
	if err != nil {
		t.Fatal(err)
	}

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

	// project root from sys root
	root, err := shared.FindProjectRoot()
	if err != nil {
		t.Fatal(err)
	}

	// path to migrate folder
	migrationsPath := filepath.Join(root, "migrations")

	// initialize migration on db
	err = testdata.InitMigrate(migrationsPath, *c)
	if err != nil {
		t.Fatal(err)
	}

	// wrap db conns in TestSuite
	testSuite, err := testdata.NewTestSuite(pgContainerCtx, pgContainer, *c)
	if err != nil {
		t.Fatal(err)
	}
	// remove testcontainer from docker
	t.Cleanup(func() {
		if err := testSuite.PgContainer.Terminate(pgContainerCtx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	t.Run("Add", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		err := testSuite.SetupTestPg(ctx)
		require.NoError(t, err)

		assert := assert.New(t)
		// require := require.New(t)

		departamentRepo := postgres.NewDepartmentRepository(testSuite.Db)

		tests := []struct {
			name    string
			expDept *domain.Department
			expErr  error
		}{
			{"normal add", &domain.Department{Name: "Test Departament", ParentId: nil}, nil},
		}

		for _, tt := range tests {
			// Act
			actErr := departamentRepo.Add(ctx, tt.expDept)

			// Assert
			assert.Equal(tt.expDept.Name, "Test Departament")
			assert.ErrorIs(actErr, tt.expErr)
		}
	})
}
