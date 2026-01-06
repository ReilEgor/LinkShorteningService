package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		DB   `yaml:"db"`
	}

	App struct {
		Name  string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Level string `env-required:"true" yaml:"level"   env:"LOG_LEVEL"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	DB struct {
		URL string `env-required:"true" env:"DB_URL" yaml:"url"`
	}
)

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	if err := cleanenv.ReadConfig("configs/config.yaml", cfg); err != nil {
		return nil, err
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
