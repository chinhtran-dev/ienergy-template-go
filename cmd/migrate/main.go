package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func main() {
	var (
		action string
		steps  int
	)

	flag.StringVar(&action, "action", "up", "Migration action (up, down, force, version)")
	flag.IntVar(&steps, "steps", 0, "Number of steps to migrate (for up/down)")
	flag.Parse()

	// Load configuration
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Get database connection string
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")
	dbSSLMode := viper.GetString("SSL_MODE")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	// Create migration instance
	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		log.Fatalf("Error creating migration instance: %v", err)
	}
	defer m.Close()

	// Execute migration action
	switch action {
	case "up":
		if steps > 0 {
			err = m.Steps(steps)
		} else {
			err = m.Up()
		}
	case "down":
		if steps > 0 {
			err = m.Steps(-steps)
		} else {
			err = m.Down()
		}
	case "force":
		version := flag.Arg(0)
		if version == "" {
			log.Fatal("Version is required for force action")
		}
		err = m.Force(cast.ToInt(version))
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Error getting version: %v", err)
		}
		fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)
		os.Exit(0)
	default:
		log.Fatalf("Unknown action: %s", action)
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No changes to apply")
			os.Exit(0)
		}
		log.Fatalf("Error executing migration: %v", err)
	}

	fmt.Printf("Migration %s completed successfully\n", action)
}
