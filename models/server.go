package models

import "gorm.io/gorm"

type Server struct {
	gorm.Model
	Name                 string
	Address              string
	Port                 string
	CfAccessClientId     string
	CfAccessClientSecret string
}
