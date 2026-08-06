package config

type Config struct {
	Version string
	Name    string
}

func NewConfig() Config {
	return Config{
		Version: "0.1.0",
		Name:    "allium",
	}
}
