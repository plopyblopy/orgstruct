package integration

import (
	"context"
	"testing"

	"github.com/plopyblopy/orgstruct/internal/adapters/postgres"
	"github.com/plopyblopy/orgstruct/internal/domain"
	"github.com/plopyblopy/orgstruct/tests/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// В этом тесте используется Testcontainer где в начале каждого теста вызывается testSuite.SetupTestPg - он закрывает все подключения, делает Rollback до начального состояния
func TestDepartmentRepository(t *testing.T) {
	// setup
	TestSetup(t)
	env, err := testdata.GetTestEnv()
	if err != nil {
		t.Fatal(err)
	}

	// tests
	t.Run("Add", func(t *testing.T) {
		// Arrange
		tests := []struct {
			name    string
			expDept *domain.Department
			expErr  error
		}{
			{
				name: "normal add",
				expDept: &domain.Department{
					Id:       1,
					Name:     "Test Departament",
					ParentId: nil,
				},
				expErr: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)
				ctx := context.Background()

				err := env.SetupTestPg(ctx)
				require.NoError(err)

				dept := postgres.NewDepartmentRepository(env.Db)

				// Act
				parentIdPtr := testdata.CopyInt(tt.expDept.ParentId)

				actDept, err := domain.NewDepartment(tt.expDept.Name, parentIdPtr)
				require.NoError(err)

				actErr := dept.Add(ctx, actDept)

				// Assert
				require.ErrorIs(actErr, tt.expErr)

				assert.Equal(tt.expDept.Id, actDept.Id, "Id")
				assert.Equal(tt.expDept.Name, actDept.Name, "Name")
				assert.Equal(tt.expDept.ParentId, actDept.ParentId, "ParentId")
			})
		}
	})
}
