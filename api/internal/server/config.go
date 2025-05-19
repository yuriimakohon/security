package server

import "github.com/spf13/viper"

type Config struct {
	Port         string
	SessionKey   []byte
	SessionName  string
	CSRFKey      []byte
	TemplatesDir string
	DatabaseURL  string
}

// In real project you should use some config file, like YAML, JSON or use ENV.

func NewConfig() Config {
	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", "9000")
	viper.SetDefault("DB_URL", "postgres://postgres:postgres@localhost:5439/postgres?sslmode=disable")

	return Config{
		Port:         viper.GetString("APP_PORT"),
		SessionKey:   []byte("secret-key"),
		SessionName:  "general",
		CSRFKey:      []byte("Yvj-ZSpifRE_rtGnu9tlUQlTvwmm9JoZ"),
		TemplatesDir: "./templates",
		DatabaseURL:  viper.GetString("DB_URL"),
	}
}
