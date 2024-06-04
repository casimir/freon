package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// copy of gorm.Model with desc tags
type Model struct {
	ID        uint           `gorm:"primarykey" desc:"readonly"`
	CreatedAt time.Time      `desc:"readonly"`
	UpdatedAt time.Time      `desc:"readonly"`
	DeletedAt gorm.DeletedAt `gorm:"index" desc:"hidden"`
}

// same than Model but with UUID
type ModelUUID struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" desc:"readonly"`
	CreatedAt time.Time      `desc:"readonly"`
	UpdatedAt time.Time      `desc:"readonly"`
	DeletedAt gorm.DeletedAt `gorm:"index" desc:"hidden"`
}

func (m *ModelUUID) BeforeCreate(tx *gorm.DB) (err error) {
	ID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	m.ID = ID
	return nil
}
