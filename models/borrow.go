package models

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type Borrow struct {
	gorm.Model
	DeadlineDate time.Time
	ReturnDate   sql.NullTime
	BookID       uint `gorm:"primaryKey"`
	UserID       uint `gorm:"primaryKey"`
}
