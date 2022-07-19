package app

type Config struct {
	DB    *DBConfig
	Token *TokenConfig
}

func NewConfig() *Config {
	return &Config{
		DB:    NewDBConfig(),
		Token: &TokenConfig{},
	}
}

type DBConfig struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     int
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	DBName   string `mapstructure:"POSTGRES_DB"`
	Url      string `mapstructure:"POSTGRES_URL"`
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
