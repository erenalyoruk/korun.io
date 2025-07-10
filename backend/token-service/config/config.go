package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"korun.io/shared/config"
)

type Config struct {
	App         AppConfig                   `mapstructure:"app"`
	Server      ServerConfig                `mapstructure:"server"`
	Token       TokenConfig                 `mapstructure:"token"`
	Logger      config.LoggerConfig         `mapstructure:"logger"`
	InfraConfig config.InfrastructureConfig `mapstructure:"infrastructure"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type TokenConfig struct {
	AccessTokenSecret     string        `mapstructure:"access_token_secret"`
	RefreshTokenSecret    string        `mapstructure:"refresh_token_secret"`
	AccessTokenExpiresIn  time.Duration `mapstructure:"access_token_expires_in"`
	RefreshTokenExpiresIn time.Duration `mapstructure:"refresh_token_expires_in"`
}

func LoadConfig() (*Config, error) {
	viper.SetEnvPrefix("KORUN")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("app.name", "token-service")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "json")
	viper.SetDefault("token.access_token_secret", "12my_secret_key_34")
	viper.SetDefault("token.refresh_token_secret", "56my_refresh_secret_78")
	viper.SetDefault("token.access_token_expires_in", "15m")
	viper.SetDefault("token.refresh_token_expires_in", "720h")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}

	return &cfg, nil
}
