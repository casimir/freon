package database

import (
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getEnvAsBool(name string) bool {
	v := os.Getenv(name)
	if val, err := strconv.ParseBool(v); err == nil {
		return val
	}

	return false
}

var DB *gorm.DB

func init() {
	dbPath := os.Getenv("FREON_DB_PATH")
	config := gorm.Config{}
	if testing.Testing() {
		dbPath = "file::memory:"
	} else if dbPath == "" {
		dbPath = "freon.db"
	}

	if !getEnvAsBool("LOG_SQL") {
		config.Logger = logger.Discard
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
