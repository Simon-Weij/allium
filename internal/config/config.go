package config

type Config struct {
	Version string
	Name    string
}

func NewConfig() Config {
	return Config{
		Name:    "allium",
		Version: "0.1.0",
	}
}
