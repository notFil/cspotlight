package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server Server
	Store  Store
	Token  Token
}

type Server struct {
	Port        int
	Environment string
	BaseURL     string
}

type Token struct {
	SecretKey              string
	RefreshSecretKey       string
	ExpiryInMinutes        int
	RefreshExpiryInMinutes int
}

type Store struct {
	DSN   string
	Redis string
}

// init sets Viper to read .env and environment variables
func init() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf(".env not loaded: %v", err)
	}

	viper.AutomaticEnv()
}

// LoadConfig returns string for key from Viper
func LoadConfig() Config {
	server := Server{
		Port:        viper.GetInt("APP_PORT"),
		Environment: viper.GetString("APP_ENVIRONMENT"),
		BaseURL:     viper.GetString("APP_BASE_URL"),
	}

	token := Token{
		SecretKey:              viper.GetString("JWT_SECRET_KEY"),
		ExpiryInMinutes:        viper.GetInt("JWT_ACCESS_TOKEN_EXPIRY_MINUTES"),
		RefreshSecretKey:       viper.GetString("JWT_REFRESH_SECRET_KEY"),
		RefreshExpiryInMinutes: viper.GetInt("JWT_REFRESH_TOKEN_EXPIRY_MINUTES"),
	}

	store := Store{
		DSN:   viper.GetString("DB_ADDR"),
		Redis: viper.GetString("REDIS_ADDR"),
	}
	cfg := Config{
		Server: server,
		Token:  token,
		Store:  store,
	}
	return cfg
}
