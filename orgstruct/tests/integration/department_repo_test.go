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

	t.Run("GetWithDepthHierarchy", func(t *testing.T) {
		// Arrange
		tests := []struct {
			name     string
			expDepts []*domain.Department
			expErr   error
		}{
			{
				name: "chain formatting",
				expDepts: []*domain.Department{
					{
						Id:       0,
						Name:     "Test Departament",
						ParentId: nil,
					},
					{
						Id:       0,
						Name:     "Test Departament 2",
						ParentId: new(1),
					},
					{
						Id:       0,
						Name:     "Test Departament 3",
						ParentId: new(2),
					},
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
				for i := range tt.expDepts {
					err = dept.Add(ctx, new(*tt.expDepts[i]))
					require.NoError(err)
				}

				actResp, actErr := dept.GetWithDepthHierarchy(ctx, 1, 3)

				// Assert
				require.ErrorIs(actErr, tt.expErr)

				for i := 0; i < len(tt.expDepts); i++ {
					assert.Equal(i+1, actResp[i].Id, "Id")
					assert.Equal(tt.expDepts[i].Name, actResp[i].Name, "Name")
					assert.Equal(tt.expDepts[i].ParentId, actResp[i].ParentId, "ParentId")
				}
			})
		}
	})

	t.Run("Update", func(t *testing.T) {
		// Arrange
		tests := []struct {
			name     string
			updDept  *domain.UpdateDepartment
			addDepts []domain.Department
			expDept  domain.Department
			expErr   error
		}{
			{
				name: "name",
				updDept: &domain.UpdateDepartment{
					Id:        3,
					Name:      domain.UpdateField[*string]{Value: new("Test Department 3 updated"), Valid: true},
					ParentId:  domain.UpdateField[*int]{Valid: false},
					CreatedAt: domain.UpdateField[*time.Time]{Valid: false},
				},
				addDepts: []domain.Department{
					{Id: 1, Name: "Test Department 1", ParentId: nil},
					{Id: 2, Name: "Test Department 2", ParentId: new(1)},
					{Id: 3, Name: "Test Department 3", ParentId: new(2)},
				},
				expDept: domain.Department{Id: 3, Name: "Test Department 3 updated", ParentId: new(2)},
				expErr:  nil,
			},
			{
				name: "parentId",
				updDept: &domain.UpdateDepartment{
					Id:        3,
					Name:      domain.UpdateField[*string]{Valid: false},
					ParentId:  domain.UpdateField[*int]{Value: new(1), Valid: true},
					CreatedAt: domain.UpdateField[*time.Time]{Valid: false},
				},
				addDepts: []domain.Department{
					{Id: 1, Name: "Test Department 1", ParentId: nil},
					{Id: 2, Name: "Test Department 2", ParentId: new(1)},
					{Id: 3, Name: "Test Department 3", ParentId: new(2)},
				},
				expDept: domain.Department{Id: 3, Name: "Test Department 3", ParentId: new(1)},
				expErr:  nil,
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
				for i := range tt.addDepts {
					err = dept.Add(ctx, new(tt.addDepts[i]))
					require.NoError(err)
				}

				actResp, actErr := dept.Update(ctx, *tt.updDept)

				// Assert
				require.ErrorIs(actErr, tt.expErr)

				assert.Equal(actResp.Id, tt.expDept.Id, "Id")
				assert.Equal(actResp.Name, tt.expDept.Name, "Name")
				assert.Equal(actResp.ParentId, tt.expDept.ParentId, "ParentId")
			})
		}
	})
}
