package integration

import (
	"context"
	"testing"
	"time"

	"github.com/plopyblopy/orgstruct/internal/adapters/postgres"
	"github.com/plopyblopy/orgstruct/internal/domain"
	"github.com/plopyblopy/orgstruct/tests/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// // В этом тесте используется Testcontainer где в начале каждого теста вызывается testSuite.SetupTestPg - он закрывает все подключения, делает Rollback до начального состояния
func TestEmployeeRepository(t *testing.T) {
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
			name   string
			depts  []*domain.Department
			expEmp *domain.Employee
			expErr error
		}{
			{
				name: "normal add",
				depts: []*domain.Department{
					{
						Id:        1,
						Name:      "Dept 1",
						ParentId:  nil,
						CreatedAt: time.Now(),
					},
				},
				expEmp: &domain.Employee{
					Id:           1,
					FullName:     "Full Name Employee",
					DepartmentId: 1,
					Position:     "Developer",
					HiredAt:      nil,
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
				emp := postgres.NewEmployeeRepository(env.Db)

				// Act
				for i := range tt.depts {
					err = dept.Add(ctx, tt.depts[i])
					require.NoError(err)
				}

				hiredAtCopy := testdata.CopyDateTime(tt.expEmp.HiredAt)

				actEmp, err := domain.NewEmployee(tt.expEmp.DepartmentId, tt.expEmp.FullName, tt.expEmp.Position, hiredAtCopy)
				require.NoError(err)

				actErr := emp.Add(ctx, actEmp)

				// Assert
				require.ErrorIs(actErr, tt.expErr)

				assert.Equal(tt.expEmp.Id, actEmp.Id, "Id")
				assert.Equal(tt.expEmp.FullName, actEmp.FullName, "FullName")
				assert.Equal(tt.expEmp.DepartmentId, actEmp.DepartmentId, "DepartmentId")
				assert.Equal(tt.expEmp.Position, actEmp.Position, "Position")
				assert.Equal(tt.expEmp.HiredAt, actEmp.HiredAt, "HiredAt")
			})
		}
	})
}
