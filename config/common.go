package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

type ServiceConfig struct {
	DatabaseConfig DatabaseConfig `envconfig:"DB"`
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (ServiceConfig, error) {
	var cfg ServiceConfig

	// load from .env if exists
	if _, err := os.Stat(".env"); err == nil {
		if err := gotenv.Load(); err != nil {
			return cfg, err
		}
	}

	// parse environment variable to config struct
	err := envconfig.Process("service", &cfg)
	return cfg, err
}

type DatabaseConfig struct {
	Host     string `envconfig:"HOST"`
	Port     int    `envconfig:"PORT"`
	Username string `envconfig:"USERNAME"`
	Password string `envconfig:"PASSWORD"`
	Database string `envconfig:"NAME"`
}

// BuildDSN constructs the MySQL DSN in URL format
func (db DatabaseConfig) DatabaseDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.Database,
	)
}

func initDatabase(cfg ServiceConfig) *sqlx.DB {
	fmt.Println("DEBUG", cfg.DatabaseConfig.DatabaseDSN())
	db := sqlx.MustOpen("mysql", cfg.DatabaseConfig.DatabaseDSN())

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	return db
}
