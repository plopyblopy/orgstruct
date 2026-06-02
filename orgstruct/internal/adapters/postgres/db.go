package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DbConfig конфигурация для базы данных.
type DbConfig struct {
	Database string `env:"DB_DATABASE"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Sslmode  string `env:"DB_SSLMODE"`
	MinConns int    `env:"DB_MINCONNS"`
	MaxConns int    `env:"DB_MAXCONNS"`
}

// Db хранит указатель на gorm.DB и строку dsn.
type Db struct {
	db  *gorm.DB
	dsn string
}

// NewDb конструктор для Db.
func NewDb(c DbConfig) *Db {
	return &Db{dsn: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host,
		c.Username,
		c.Password,
		c.Database,
		c.Port,
		c.Sslmode,
	)}
}

// Open открывает подключение к базе данных.
func (db *Db) Open(c DbConfig) error {
	gormDb, err := gorm.Open(postgres.Open(db.dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	db.db = gormDb

	sqlDb, err := db.db.DB()
	if err != nil {
		return err
	}

	sqlDb.SetMaxOpenConns(c.MaxConns)

	err = sqlDb.Ping()
	if err != nil {
		return err
	}

	return nil
}

// Close закрывает подключение к базе данных.
func (db *Db) Close() error {
	sqlDb, err := db.db.DB()
	if err != nil {
		return err
	}

	err = sqlDb.Close()
	if err != nil {
		return err
	}

	return nil
}
