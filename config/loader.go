package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"POSTGRES_HOST"`
	DBUsername string `mapstructure:"POSTGRES_USER"`
	DBPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName     string `mapstructure:"POSTGRES_DB"`
	DBPort     string `mapstructure:"POSTGRES_PORT"`

	ServerPort string `mapstructure:"PORT"`

	AccessTokenSecret    string        `mapstructure:"ACCESS_TOKEN_SECRET"`
	AccessTokenExpiresIn time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge    int           `mapstructure:"ACCESS_TOKEN_MANAGE"`

	RefreshTokenSecret    string        `mapstructure:"REFRESH_TOKEN_SECRET"`
	RefreshTokenExpiresIn time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	RefreshTokenMaxAge    int           `mapstructure:"REFRESH_TOKEN_MANAGE"`

	EmailFrom string `mapstructure:"EMAIL_FROM"`
	SMTPHost  string `mapstructure:"SMTP_HOST"`
	SMTPPass  string `mapstructure:"SMTP_PASS"`
	SMTPPort  int    `mapstructure:"SMTP_PORT"`
	SMTPUser  string `mapstructure:"SMTP_USER"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
