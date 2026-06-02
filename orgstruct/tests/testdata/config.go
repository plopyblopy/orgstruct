package testdata

import (
	"path/filepath"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/plopyblopy/orgstruct/internal/adapters/postgres"
	"github.com/plopyblopy/orgstruct/internal/adapters/rest"
)

type TestConfig struct {
	// HTTP
	HttpConfig rest.HttpConfig

	// DB
	Image        string `env:"DB_IMAGE"`
	ShapshotName string `env:"DB_SNAPSHOTNAME"`
	DbConfig     postgres.DbConfig
}

// load .env from current folder
func NewTestConfig() (*TestConfig, error) {
	c := &TestConfig{}

	curdir := GetCurDirPath()

	err := godotenv.Load(filepath.Join(curdir, ".env"))
	if err != nil {
		return nil, err
	}

	err = env.Parse(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
