package config

type Database struct {
	URL string `env:"URL" envDefault:"postgres://postgres:pass2word@localhost:5432/wongnok?sslmode=disable"`
}
