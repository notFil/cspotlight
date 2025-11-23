package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig
	Store   StoreConfig
	JWTAuth JWTConfig
}

type ServerConfig struct {
	Port        int
	Environment string
}

type JWTConfig struct {
	SecretKey       string
	ExpiryInMinutes int
}

type StoreConfig struct {
	DatabaseURL string
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
	server := ServerConfig{
		Port: viper.GetInt("PORT"),
	}

	jwtAuth := JWTConfig{
		SecretKey:       viper.GetString("JWT_SECRET_KEY"),
		ExpiryInMinutes: viper.GetInt("JWT_SECRET_EXPIRY_MINUTES"),
	}

	store := StoreConfig{
		DatabaseURL: viper.GetString("DB_URL"),
	}
	cfg := Config{
		Server:  server,
		JWTAuth: jwtAuth,
		Store:   store,
	}
	return cfg
}
