package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type InfrastructureConfig struct {
	Logger LoggerConfig `mapstructure:"logger"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func LoadInfrastructureConfig() (*InfrastructureConfig, error) {
	viper.SetEnvPrefix("KORUN")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "json")

	var cfg InfrastructureConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}

	return &cfg, nil
}
