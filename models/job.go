package models

import "gorm.io/gorm"

type Job struct {
	gorm.Model
	Name     string
	Id       string
	ServerID int
	Server   Server
}
