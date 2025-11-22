package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Test     bool
}

const configPath = "./config.env"

func LoadConfig() *Config {
	var cfg Config

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("env")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	cfg.Server.Port = v.GetString("SERVER_PORT")

	cfg.Database.Host = v.GetString("DATABASE_HOST")
	cfg.Database.Port = v.GetInt("DATABASE_PORT")
	cfg.Database.User = v.GetString("DATABASE_USER")
	cfg.Database.Password = v.GetString("DATABASE_PASSWORD")
	cfg.Database.Name = v.GetString("DATABASE_NAME")
	cfg.Database.Test = v.GetBool("DATABASE_TEST")

	return &cfg
}
