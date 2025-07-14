package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common columns for all tables
type BaseModel struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
