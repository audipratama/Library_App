package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Name   string
	Writer string
	Year   int
	Owned  int
	User []*User `gorm:"many2many:borrows;"`
}
