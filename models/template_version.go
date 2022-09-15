package models

import (
	"gorm.io/gorm"
	"time"
)

type TemplateVersion struct {
	gorm.Model
	Selector       string
	Template       Template
	TemplateID     int
	LastDeployedAt time.Time
	LastVersion    string
}
