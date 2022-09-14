package models

import "gorm.io/gorm"

type Template struct {
	gorm.Model
	JobID int
	Job   Job
}
