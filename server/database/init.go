package database

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	dbPath := os.Getenv("FREON_DB_PATH")
	if dbPath == "" {
		dbPath = "freon.db"
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	if gin.Mode() == gin.DebugMode {
		db = db.Debug()
	}
	DB = db
}
