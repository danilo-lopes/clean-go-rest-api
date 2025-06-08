package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort int
	TimeZone   string
	DB         DatabaseConfig
}

type DatabaseConfig struct {
	Host                 string
	User                 string
	Password             string
	Name                 string
	Port                 string
	MigrationsFolderPath string
	Parameters           string
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func LoadConfig() *Config {
	return &Config{
		ServerPort: getPort(),
		TimeZone:   getEnv("TIME_ZONE", "America/Sao_Paulo"),
		DB: DatabaseConfig{
			Host:                 getEnv("DB_HOST", "localhost"),
			User:                 getEnv("DB_USER", "postgres"),
			Password:             getEnv("DB_PASSWORD", "postgres"),
			Name:                 getEnv("DB_NAME", "go_rest_api"),
			Port:                 getEnv("DB_PORT", "5432"),
			MigrationsFolderPath: getEnv("DB_MIGRATIONS_FOLDER", "file://migrations"),
			Parameters:           getEnv("DB_PARAMETERS", ""),
		},
	}
}

func (c *Config) DBConnectionString() string {
	return "host=" + c.DB.Host +
		" user=" + c.DB.User +
		" password=" + c.DB.Password +
		" dbname=" + c.DB.Name +
		" port=" + c.DB.Port +
		" " + c.DB.Parameters
}

func getPort() int {
	port := getEnv("PORT", "8080")
	return atoiOrDefault(port, 8080)
}

func atoiOrDefault(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
