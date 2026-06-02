package testdata

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/pressly/goose/v3"
)

func InitMigrate(folderPath string, c TestConfig) error {
	rawPath := filepath.ToSlash(folderPath)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.DbConfig.Host,
		c.DbConfig.Username,
		c.DbConfig.Password,
		c.DbConfig.Database,
		c.DbConfig.Port,
		c.DbConfig.Sslmode,
	)

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	goose.SetDialect("postgres")
	return goose.Up(sqlDB, rawPath)
}
