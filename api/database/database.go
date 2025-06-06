package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the database instance
var DB *gorm.DB

// Connect establishes a connection to the database
func Connect(dsn string) error {
	var err error

	// Custom logger configuration for GORM
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // Log slow queries
			LogLevel:                  logger.Info, // Log all SQL
			IgnoreRecordNotFoundError: true,        // Ignore not found errors
			Colorful:                  true,        // Enable color
		},
	)

	// Connect to the database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return err
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Connected to the database")
	return nil
}

// Migrate automatically migrates all models
func Migrate(models ...interface{}) error {
	log.Println("Running database migrations")
	return DB.AutoMigrate(models...)
}
