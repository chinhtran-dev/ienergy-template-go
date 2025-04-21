package config

import (
	"log"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Config is the top-level configuration struct
type Config struct {
	DB     DBConfig
	JWT    JWTConfig
	Server ServerCfg
}

// DBConfig holds the database-related configuration values
type DBConfig struct {
	Host               string `envconfig:"DB_HOST" default:"localhost"`      // Database host
	Port               string `envconfig:"DB_PORT" default:"5432"`           // Database port
	User               string `envconfig:"DB_USER" default:"postgres"`       // Database user
	Password           string `envconfig:"DB_PASSWORD" default:"postgres"`   // Database password
	DBName             string `envconfig:"DB_NAME" default:"postgres"`       // Database name
	SSLMode            string `envconfig:"SSL_MODE" default:"disable"`       // SSL mode for database connection
	SetMaxIdleConns    string `envconfig:"SET_MAX_IDLE_CONNS" default:""`    // Max idle connections
	SetMaxOpenConns    string `envconfig:"SET_MAX_OPEN_CONNS" default:""`    // Max open connections
	SetConnMaxLifetime string `envconfig:"SET_CONN_MAX_LIFETIME" default:""` // Connection max lifetime
}

// JWTConfig holds the JWT-related configuration values
type JWTConfig struct {
	Secret                string `envconfig:"JWT_SECRET"`                  // JWT secret key
	ExpirationTime        string `envconfig:"JWT_EXPIRATION_TIME"`         // JWT expiration time
	RefreshSecret         string `envconfig:"JWT_REFRESH_SECRET"`          // JWT refresh token secret key
	RefreshExpirationTime string `envconfig:"JWT_REFRESH_EXPIRATION_TIME"` // JWT refresh token expiration time
}

// ServerCfg holds the server-related configuration values
type ServerCfg struct {
	ServerURL  string `envconfig:"SERVER_URL" default:"localhost"`    // Server URL
	Port       string `envconfig:"PORT" default:"8080"`               // Server port
	Env        string `envconfig:"ENVIRONMENT" default:"development"` // Environment (e.g., development, production)
	GINMode    string `envconfig:"GIN_MODE" default:"debug"`          // Gin framework mode
	Production bool   `envconfig:"PRODUCTION" default:"false"`        // Is production environment
}

func NewConfig() (*Config, error) {
	LoadConfig()

	var cfg Config

	if err := envconfig.Process("", &cfg.DB); err != nil {
		log.Fatalf("Failed to process DB config: %v", err)
	}
	if err := envconfig.Process("", &cfg.JWT); err != nil {
		log.Fatalf("Failed to process JWT config: %v", err)
	}
	if err := envconfig.Process("", &cfg.Server); err != nil {
		log.Fatalf("Failed to process Server config: %v", err)
	}

	return &cfg, nil
}

// LoadConfig loads configuration values from environment variables
func LoadConfig() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	for _, env := range viper.AllKeys() {
		if viper.GetString(env) != "" {
			_ = os.Setenv(env, viper.GetString(env))
			_ = os.Setenv(strings.ToUpper(env), viper.GetString(env))
		}
	}
}

func ServerConfig() *Config {
	return &Config{}
}

var Module = fx.Options(
	fx.Provide(NewConfig),
)
