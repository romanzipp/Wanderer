package models

import "gorm.io/gorm"

type Server struct {
	gorm.Model
	Address              string
	Port                 string
	CfAccessClientId     string
	CfAccessClientSecret string
}
