package testdata

import (
	"context"
	"errors"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// testEnv хранит указатель на testdata.TestSuite.
type Env struct {
	*TestSuite
}

// TestEnv глобальная переменная пакета. Нужна для согласованности тестов.
var testEnv *Env

// InitTestEnv для инициализации глобальной переменной пакета - testEnv, если не была создана ранее.
func InitTestEnv(ctx context.Context, pgContainer *postgres.PostgresContainer, c *TestConfig) error {
	testSuite, err := NewTestSuite(ctx, pgContainer, *c)
	if err != nil {
		return err
	}

	testEnv = &Env{
		TestSuite: testSuite,
	}
	return nil
}

// для получения
func GetTestEnv() (*Env, error) {
	if testEnv != nil {
		return testEnv, nil
	} else {
		return nil, errors.New("TestEnv is empty. First you need to perform initialization.")
	}
}
