package models

import "gorm.io/gorm"

type Template struct {
	gorm.Model
	Name       string
	NomadJobID string
	Content    string
}
