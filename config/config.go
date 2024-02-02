package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      `yaml:"app"`
		Http     `yaml:"http"`
		Log      `yaml:"log"`
		Database `yaml:"database"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-requited:"true" yaml:"version" env:"APP_VERSION"`
	}

	Http struct {
		Port      string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		StaticDir string `env-requited:"true" yaml:"staticDir" env:"HTTP_STATIC_DIR"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	}

	Database struct {
		Port     string `yaml:"db_port" env:"DB_PORT" env-default:"7777"`
		Host     string `yaml:"db_host" env:"DB_HOST" env-default:"localhost"`
		Name     string `yaml:"db_name" env:"DB_NAME" env-default:"postgres"`
		User     string `yaml:"db_user" env:"DB_USER" env-default:"postgres"`
		Password string `yaml:"db_password" env:"DB_PASSWORD"`
		URL      string `env-required:"true" env:"DB_URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig("./config/config.yml", cfg); err != nil {
		return nil, fmt.Errorf("config file error: %v", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
