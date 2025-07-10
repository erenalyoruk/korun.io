package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"korun.io/shared/config"
)

type Config struct {
	App    AppConfig                    `mapstructure:"app"`
	Server ServerConfig                 `mapstructure:"server"`
	Token  TokenConfig                  `mapstructure:"token"`
	Infra  *config.InfrastructureConfig `mapstructure:"infra"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type TokenConfig struct {
	AccessTokenSecret     string        `mapstructure:"access_token_secret"`
	RefreshTokenSecret    string        `mapstructure:"refresh_token_secret"`
	AccessTokenExpiresIn  time.Duration `mapstructure:"access_token_expires_in"`
	RefreshTokenExpiresIn time.Duration `mapstructure:"refresh_token_expires_in"`
}

func LoadConfig(configPath, configName string) (*Config, error) {
	viper.SetEnvPrefix("KORUN")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	viper.SetDefault("app.name", "token-service")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.write_timeout", "15s")
	viper.SetDefault("server.idle_timeout", "60s")
	viper.SetDefault("token.access_token_secret", "12my_secret_key_34")
	viper.SetDefault("token.refresh_token_secret", "56my_refresh_secret_78")
	viper.SetDefault("token.access_token_expires_in", "15m")
	viper.SetDefault("token.refresh_token_expires_in", "720h")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}

	infraConfig, err := config.LoadInfrastructureConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load infrastructure config: %w", err)
	}
	cfg.Infra = infraConfig

	return &cfg, nil
}
