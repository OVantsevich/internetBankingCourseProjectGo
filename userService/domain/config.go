package domain

type Config struct {
	DatabaseUrl string `env:"POSTGRES_DB_URL"`
	JwtKey      string `env:"JWT_KEY" `
}
