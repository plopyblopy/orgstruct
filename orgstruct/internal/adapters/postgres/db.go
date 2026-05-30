package postgres

// DbConfig конфигурация для базы данных
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
