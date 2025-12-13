package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server  Server
	Session Session
	Store   Store
}

type Server struct {
	Port        int
	Environment string
	BaseURL     string
}

type Session struct {
	ExpiryInMinutes int
	SecretKey       []byte
	UseCookieStore  bool
}

type Store struct {
	DSN   string
	Redis Redis
}

type Redis struct {
	Addr      string
	Username  string
	Password  string
	IdleConns int
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

	session := Session{
		ExpiryInMinutes: viper.GetInt("SESSION_EXPIRY_IN_MINUTES"),
		SecretKey:       []byte(viper.GetString("SESSION_SECRET_KEY")),
		UseCookieStore:  viper.GetBool("SESSION_USE_COOKIE_STORE"),
	}

	store := Store{
		DSN: viper.GetString("DB_ADDR"),
		Redis: Redis{
			Addr:      viper.GetString("REDIS_ADDR"),
			Username:  viper.GetString("REDIS_USERNAME"),
			Password:  viper.GetString("REDIS_PASSWORD"),
			IdleConns: viper.GetInt("REDIS_IDLE_CONNS"),
		},
	}
	cfg := Config{
		Server:  server,
		Session: session,
		Store:   store,
	}
	return cfg
}
