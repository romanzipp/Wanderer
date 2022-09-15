package models

import (
	"gorm.io/gorm"
	"time"
)

type Token struct {
	gorm.Model
	Name       string
	Token      string
	LastUsedAt time.Time
}
