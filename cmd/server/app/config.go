package app

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	DB    *DBConfig    `mapstructure:",squash"`
	Token *TokenConfig `mapstructure:",squash"`
}

func NewConfig() (*Config, error) {
	c := &Config{
		DB:    NewDBConfig(),
		Token: &TokenConfig{},
	}

	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal configs")
	}
	return c, err
}

type DBConfig struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     int
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	DBName   string `mapstructure:"POSTGRES_DB"`
	URL      string `mapstructure:"POSTGRES_URL"`
}

type TokenConfig struct {
	PASETOSecret string `mapstructure:"PASETO_SECRET"`
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		Port: 5432,
		Host: "localhost",
	}
}
