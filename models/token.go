package models

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	Name  string
	Token string
}
