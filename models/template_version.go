package models

import "gorm.io/gorm"

type TemplateVersion struct {
	gorm.Model
	Selector   string
	Template   Template
	TemplateID int
}
