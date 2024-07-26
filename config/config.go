package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Configs struct {
	APP app
	DB  DB
}

type DB struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
}

type app struct {
	Port string
	Path string
}

func New() (cfg Configs, err error) {
	root, err := os.Getwd()
	if err != nil {
		return
	}

	err = godotenv.Load(filepath.Join(root, ".env"))
	if err != nil {
		cfg.APP.Port = os.Getenv("APP_PORT")

		cfg.DB.Username = os.Getenv("DB_USERNAME")
		cfg.DB.Password = os.Getenv("DB_PASSWORD")
		cfg.DB.Host = os.Getenv("DB_HOST")
		cfg.DB.Port = os.Getenv("DB_PORT")
		cfg.DB.Name = os.Getenv("DB_NAME")

		return cfg, nil
	}

	if err = envconfig.Process("DB", &cfg.DB); err != nil {
		return
	}

	if err = envconfig.Process("APP", &cfg.APP); err != nil {
		return
	}

	return
}
