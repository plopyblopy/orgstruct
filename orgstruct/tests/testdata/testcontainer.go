package testdata

import (
	"context"

	postgresDb "github.com/plopyblopy/orgstruct/internal/adapters/postgres"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Creates postgres container.
func NewPostgresTestcontainer(ctx context.Context, c TestConfig) (*postgres.PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx, c.Image,
		postgres.WithDatabase(c.DbConfig.Database),
		postgres.WithUsername(c.DbConfig.Username),
		postgres.WithPassword(c.DbConfig.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		return nil, err
	}

	return pgContainer, nil
}

// TestSuite используется для переиспользования Testcontainer
// приводя базу данных к начальному состоянию.
type TestSuite struct {
	PgContainer  *postgres.PostgresContainer
	Db           *postgresDb.Db
	c            TestConfig
	snapshotName string
}

func NewTestSuite(ctx context.Context, pgContainer *postgres.PostgresContainer, c TestConfig) (*TestSuite, error) {
	snapshotName := c.ShapshotName

	err := pgContainer.Snapshot(ctx, postgres.WithSnapshotName(snapshotName))
	if err != nil {
		return nil, err
	}

	db := postgresDb.NewDb(c.DbConfig)

	err = db.Open(c.DbConfig)
	if err != nil {
		return nil, err
	}

	return &TestSuite{
		PgContainer:  pgContainer,
		Db:           db,
		c:            c,
		snapshotName: snapshotName,
	}, nil
}

// Completes active connections, applies base snapshot, restoring the database to its original form.
func (ts *TestSuite) SetupTestPg(ctx context.Context) error {
	closeConns(ts.Db)

	err := ts.PgContainer.Restore(ctx, postgres.WithSnapshotName(ts.snapshotName))
	if err != nil {
		return err
	}

	db, err := newPostgresConnection(ctx, ts.PgContainer, ts.c)
	if err != nil {
		return err
	}

	ts.Db = db

	return nil
}

func newPostgresConnection(ctx context.Context, pgContainer *postgres.PostgresContainer, c TestConfig) (*postgresDb.Db, error) {
	db := postgresDb.NewDb(c.DbConfig)

	err := db.Open(c.DbConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func closeConns(db *postgresDb.Db) {
	if db != nil {
		db.Close()
	}
}
