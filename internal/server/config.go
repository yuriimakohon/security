package server

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
	return Config{
		Port:         "9000",
		SessionKey:   []byte("secret-key"),
		SessionName:  "general",
		CSRFKey:      []byte("Yvj-ZSpifRE_rtGnu9tlUQlTvwmm9JoZ"),
		TemplatesDir: "./templates",
		DatabaseURL:  "postgres://postgres:postgres@localhost:5439/postgres?sslmode=disable",
	}
}
