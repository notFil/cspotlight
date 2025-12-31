package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	CORS    CORS
	Server  Server
	Session Session
	Store   Store
}

type Server struct {
	Port        int
	Environment string
	BaseURL     string
	StaticPath  string
}

type CORS struct {
	Origins       []string
	AllowMethods  []string
	AllowHeaders  []string
	ExposeHeaders []string
}

type Session struct {
	ExpiryInSeconds int
	SecretKey       []byte
	UseCookieStore  bool
}

type Store struct {
	DSN         string
	MigratorDSN string
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
		StaticPath:  viper.GetString("APP_STATIC_PATH"),
	}

	session := Session{
		ExpiryInSeconds: viper.GetInt("SESSION_EXPIRY_IN_SECONDS"),
		SecretKey:       []byte(viper.GetString("SESSION_SECRET_KEY")),
		UseCookieStore:  viper.GetBool("SESSION_USE_COOKIE_STORE"),
	}

	store := Store{
		DSN:         viper.GetString("DB_ADDR"),
		MigratorDSN: viper.GetString("MIGRATOR_DB_ADDR"),
	}

	cors := CORS{
		Origins:       viper.GetStringSlice("CORS_ORIGINS"),
		AllowMethods:  viper.GetStringSlice("CORS_ALLOW_METHODS"),
		AllowHeaders:  viper.GetStringSlice("CORS_ALLOW_HEADERS"),
		ExposeHeaders: viper.GetStringSlice("CORS_EXPOSE_HEADERS"),
	}

	cfg := Config{
		Server:  server,
		Session: session,
		Store:   store,
		CORS:    cors,
	}
	return cfg
}

// IsProduction returns true if the environment is production
func (cfg *Config) IsProduction() bool {
	return cfg.Server.Environment == "production"
}
