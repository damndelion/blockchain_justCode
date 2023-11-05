package auth

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App   `yaml:"app"`
		HTTP  `yaml:"http"`
		Log   `yaml:"logger"`
		PG    `yaml:"postgres"`
		JWT   `yaml:"jwt"`
		Kafka `yaml:"kafka"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true"                 env:"PG_URL"`
	}
	// Jwt
	JWT struct {
		SecretKey       string `mapstructure:"secret_key" yaml:"secret_key"`
		AccessTokenTTL  int64  `mapstructure:"access_token_ttl" yaml:"access_token_ttl"`
		RefreshTokenTTL int64  `mapstructure:"refresh_token_ttl" yaml:"refresh_token_ttl"`
	}
	Kafka struct {
		Brokers  []string `yaml:"brokers"`
		Producer Producer `yaml:"producer"`
		Consumer Consumer `yaml:"consumer"`
	}
	Producer struct {
		Topic string `yaml:"topic"`
	}
	Consumer struct {
		Topics []string `yaml:"topics"`
	}
)

// NewConfig returns user config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("config/auth/config.yml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
