package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	DataPath       string
	MinimapPath    string
	AllowedOrigins string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = "../data/player_data"
	}

	minimapPath := os.Getenv("MINIMAP_PATH")
	if minimapPath == "" {
		minimapPath = dataPath + "/minimaps"
	}

	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		origins = "*"
	}

	return Config{
		Port:           port,
		DataPath:       dataPath,
		MinimapPath:    minimapPath,
		AllowedOrigins: origins,
	}
}

func (c Config) ParquetGlob() string {
	return c.DataPath + "/February_*/*"
}

func (c Config) PortInt() int {
	p, err := strconv.Atoi(c.Port)
	if err != nil {
		return 8080
	}
	return p
}
