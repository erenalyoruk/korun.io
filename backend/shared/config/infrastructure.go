package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type InfrastructureConfig struct {
	Kafka  KafkaConfig
	Redis  RedisConfig
	Logger LoggerConfig
}

type KafkaConfig struct {
	Brokers       []string
	RetryAttempts int           `mapstructure:"retry_attempts"`
	RetryBackoff  time.Duration `mapstructure:"retry_backoff"`
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type LoggerConfig struct {
	Level  string
	Format string
}

func LoadInfrastructureConfig() (*InfrastructureConfig, error) {
	viper.SetEnvPrefix("KORUN")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// set default values
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.retry_attempts", 3)
	viper.SetDefault("kafka.retry_backoff", "1s")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "json")

	var cfg InfrastructureConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal infrastructure config: %w", err)
	}

	return &cfg, nil
}
