package models

import (
	"fmt"
	"gorm.io/gorm"
)

type Template struct {
	gorm.Model
	Name       string
	NomadJobID string
	Content    string
	Server     Server
	ServerID   int
	Versions   []TemplateVersion
}

func (t Template) GetNomadJobUrl() string {
	return fmt.Sprintf("%s/ui/jobs//%s", t.Server.Address, t.NomadJobID)
}
