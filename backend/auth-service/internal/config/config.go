package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	sharedConfig "korun.io/shared/config"
)

type Config struct {
	Server   ServerConfig                      `mapstructure:"server"`
	Database DatabaseConfig                    `mapstructure:"database"`
	JWT      JWTConfig                         `mapstructure:"jwt"`
	Infra    sharedConfig.InfrastructureConfig `mapstructure:"infra"`
}

type ServerConfig struct {
	Port               int           `mapstructure:"port"`
	ReadTimeout        time.Duration `mapstructure:"read_timeout"`
	WriteTimeout       time.Duration `mapstructure:"write_timeout"`
	IdleTimeout        time.Duration `mapstructure:"idle_timeout"`
	CORSAllowedOrigins []string      `mapstructure:"cors_allowed_origins"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

func (dc *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dc.Host, dc.Port, dc.User, dc.Password, dc.DBName, dc.SSLMode)
}

type JWTConfig struct {
	Secret          string        `mapstructure:"secret"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
}

func LoadConfig(configPath string) (*Config, error) {
	vp := viper.New()
	vp.AddConfigPath(configPath)
	vp.SetConfigName("config.dev")
	vp.SetConfigType("yaml")

	vp.SetEnvPrefix("KORUN")
	vp.AutomaticEnv()
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := vp.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
