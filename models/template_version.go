package models

import (
	"github.com/dustin/go-humanize"
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

func (t TemplateVersion) GetPrettyDate() string {
	return humanize.Time(t.LastDeployedAt)
}
