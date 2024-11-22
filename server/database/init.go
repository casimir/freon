package database

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	dbPath := os.Getenv("FREON_DB_PATH")
	config := gorm.Config{}
	if testing.Testing() {
		dbPath = "file::memory:"
		config.Logger = logger.Discard
	} else if dbPath == "" {
		dbPath = "freon.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &config)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	if gin.Mode() == gin.DebugMode {
		db = db.Debug()
	}
	DB = db
}
