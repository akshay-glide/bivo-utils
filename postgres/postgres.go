package postgres

import (
	"fmt"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ConfigPostgres struct {
	Host             *string `json:"host" validate:"required"`
	Port             *int    `json:"port" validate:"required,numeric"`
	Database         *string `json:"database" validate:"required"`
	User             *string `json:"user" validate:"required"`
	Password         *string `json:"password" validate:"required"`
	DSN              *string `json:"dsn"`
	SSLMode          *string `json:"sslmode"`
	ConnMaxOpen      *int    `json:"connmaxopen" validate:"required,numeric,gt=0"`
	ConnMaxIdleTime  *int64  `json:"connmaxidletime" validate:"required,numeric,gt=30000"`
	ConnMaxIdleConns *int    `json:"connmaxidleconns" validate:"required,numeric,gt=0"`
}

func optionalOrDefault(val *string, def string) string {
	if val != nil {
		return *val
	}
	return def
}

func getPostgresDSN(cfg *ConfigPostgres) string {
	// If DSN is already provided, use it directly
	if cfg.DSN != nil && *cfg.DSN != "" {
		return *cfg.DSN
	}

	host := *cfg.Host
	port := strconv.Itoa(*cfg.Port)
	user := *cfg.User
	password := *cfg.Password
	dbname := *cfg.Database
	sslmode := optionalOrDefault(cfg.SSLMode, "disable")

	// Construct the DSN string
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=5",
		host, port, user, password, dbname, sslmode,
	)
}

// NewPostgresDB initializes a new GORM DB connection using the provided configuration
func NewPostgresDB(config *ConfigPostgres) (*gorm.DB, error) {
	dsn := getPostgresDSN(config)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Migrate runs GORM's automigrate for the given models
func Migrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}
